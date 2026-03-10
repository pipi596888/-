package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"
)

var (
	rsaOnce     sync.Once
	rsaPrivKey  *rsa.PrivateKey
	rsaSpkiB64  string
	rsaInitErr  error
)

func ensureRSAKeypair() error {
	rsaOnce.Do(func() {
		// Ephemeral keypair for encrypting credentials on the wire (UI-level obfuscation).
		// For real transport security, HTTPS is still required.
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			rsaInitErr = err
			return
		}
		spki, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
		if err != nil {
			rsaInitErr = err
			return
		}
		rsaPrivKey = key
		rsaSpkiB64 = base64.StdEncoding.EncodeToString(spki)
	})
	return rsaInitErr
}

type PublicKeyResp struct {
	Alg        string `json:"alg"`
	SpkiBase64 string `json:"spkiBase64"`
}

func getPublicKey() (PublicKeyResp, error) {
	if err := ensureRSAKeypair(); err != nil {
		return PublicKeyResp{}, err
	}
	return PublicKeyResp{Alg: "RSA-OAEP-256", SpkiBase64: rsaSpkiB64}, nil
}

func decryptPasswordEncrypted(passwordEncrypted string) (string, error) {
	if passwordEncrypted == "" {
		return "", errors.New("empty passwordEncrypted")
	}
	if err := ensureRSAKeypair(); err != nil {
		return "", err
	}
	cipherBytes, err := base64.StdEncoding.DecodeString(passwordEncrypted)
	if err != nil {
		return "", fmt.Errorf("invalid passwordEncrypted: %w", err)
	}

	plain, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaPrivKey, cipherBytes, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt failed: %w", err)
	}
	return string(plain), nil
}

