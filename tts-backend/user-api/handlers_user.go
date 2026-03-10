package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/rest/httpx"
	"golang.org/x/crypto/bcrypt"

	"tts-backend/user-api/internal/config"
)

type User struct {
	Id             int64   `json:"id"`
	Username       string  `json:"username"`
	Email          string  `json:"email"`
	Balance        float64 `json:"balance"`
	CharacterCount int64   `json:"characterCount"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	PasswordEncrypted string `json:"passwordEncrypted"`
}

type LoginResp struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	PasswordEncrypted string `json:"passwordEncrypted"`
	Email    string `json:"email"`
}

type RegisterResp struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

func getUserInfoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := mustUserID(r.Context())
		var user User
		err := db.QueryRow(
			"SELECT id, username, COALESCE(email,''), COALESCE(balance,0), COALESCE(character_count,0) FROM user WHERE id = ?",
			userId,
		).Scan(&user.Id, &user.Username, &user.Email, &user.Balance, &user.CharacterCount)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, user)
	}
}

func loginHandler(c *config.Config, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "参数错误")
			return
		}
		req.Username = strings.TrimSpace(req.Username)
		if strings.TrimSpace(req.PasswordEncrypted) != "" {
			p, err := decryptPasswordEncrypted(strings.TrimSpace(req.PasswordEncrypted))
			if err != nil {
				writeJSONError(w, http.StatusBadRequest, "invalid password")
				return
			}
			req.Password = p
		}

		var user User
		var passwordHash string
		err := db.QueryRow(
			"SELECT id, username, COALESCE(email,''), COALESCE(balance,0), COALESCE(character_count,0), password FROM user WHERE username = ?",
			req.Username,
		).Scan(&user.Id, &user.Username, &user.Email, &user.Balance, &user.CharacterCount, &passwordHash)
		if err == sql.ErrNoRows {
			writeJSONError(w, http.StatusUnauthorized, "用户名或密码错误")
			return
		}
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		if looksLikeBcryptHash(passwordHash) {
			if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)) != nil {
				writeJSONError(w, http.StatusUnauthorized, "用户名或密码错误")
				return
			}
		} else {
			// Backward compatibility (plain text passwords)
			if req.Password != passwordHash {
				writeJSONError(w, http.StatusUnauthorized, "用户名或密码错误")
				return
			}
		}

		tokenString, err := signToken(c.JwtSecret, user.Id, user.Username)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, LoginResp{
			Token: tokenString,
			User:  user,
		})
	}
}

func registerHandler(c *config.Config, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "参数错误")
			return
		}
		req.Username = strings.TrimSpace(req.Username)
		if strings.TrimSpace(req.PasswordEncrypted) != "" {
			p, err := decryptPasswordEncrypted(strings.TrimSpace(req.PasswordEncrypted))
			if err != nil {
				writeJSONError(w, http.StatusBadRequest, "invalid password")
				return
			}
			req.Password = p
		}

		if req.Username == "" || req.Password == "" {
			writeJSONError(w, http.StatusBadRequest, "用户名和密码不能为空")
			return
		}

		var exists bool
		if err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM user WHERE username = ?)", req.Username).Scan(&exists); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if exists {
			writeJSONError(w, http.StatusBadRequest, "用户名已存在")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		result, err := db.Exec(
			"INSERT INTO user (username, password, email, balance, character_count) VALUES (?, ?, ?, 0, 0)",
			req.Username, string(hashedPassword), strings.TrimSpace(req.Email),
		)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		userId, _ := result.LastInsertId()
		tokenString, err := signToken(c.JwtSecret, userId, req.Username)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, RegisterResp{
			Token: tokenString,
			User: User{
				Id:             userId,
				Username:       req.Username,
				Email:          strings.TrimSpace(req.Email),
				Balance:        0,
				CharacterCount: 0,
			},
		})
	}
}

func signToken(secret string, userId int64, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   userId,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	return token.SignedString([]byte(secret))
}

func looksLikeBcryptHash(s string) bool {
	return strings.HasPrefix(s, "$2a$") || strings.HasPrefix(s, "$2b$") || strings.HasPrefix(s, "$2y$")
}
