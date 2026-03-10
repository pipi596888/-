package model

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrVoiceForbidden = errors.New("voice forbidden")
	ErrVoiceNotFound  = errors.New("voice not found")
	ErrInvalidVoiceID = errors.New("invalid voice id")
)

type VoiceAccessModel interface {
	ValidateVoiceIds(userId int64, isAdmin bool, voiceIds []int64) error
}

type DefaultVoiceAccessModel struct {
	db *sql.DB
}

func NewVoiceAccessModel(db *sql.DB) VoiceAccessModel {
	return &DefaultVoiceAccessModel{db: db}
}

func (m *DefaultVoiceAccessModel) ValidateVoiceIds(userId int64, isAdmin bool, voiceIds []int64) error {
	uniq := make([]int64, 0, len(voiceIds))
	seen := make(map[int64]struct{}, len(voiceIds))
	for _, id := range voiceIds {
		if id <= 0 {
			return fmt.Errorf("%w: %d", ErrInvalidVoiceID, id)
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		uniq = append(uniq, id)
	}
	if len(uniq) == 0 {
		return nil
	}

	placeholders := strings.Repeat("?,", len(uniq))
	placeholders = strings.TrimSuffix(placeholders, ",")

	// 1) Ensure voice exists
	exists := make(map[int64]struct{}, len(uniq))
	args := make([]interface{}, 0, len(uniq))
	for _, id := range uniq {
		args = append(args, id)
	}
	rows, err := m.db.Query(`SELECT id FROM voice WHERE id IN (`+placeholders+`)`, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return err
		}
		exists[id] = struct{}{}
	}
	for _, id := range uniq {
		if _, ok := exists[id]; !ok {
			return fmt.Errorf("%w: %d", ErrVoiceNotFound, id)
		}
	}

	if isAdmin {
		return nil
	}

	// 2) Enforce custom voice ownership:
	// If a voice_id appears as an approved result_voice_id, only the owner can use it.
	ownerByVoice := make(map[int64]int64)
	rows2, err := m.db.Query(
		`SELECT user_id, result_voice_id FROM custom_voice_request WHERE status = 'success' AND result_voice_id IN (`+placeholders+`)`,
		args...,
	)
	if err != nil {
		return err
	}
	defer rows2.Close()
	for rows2.Next() {
		var owner int64
		var voiceId int64
		if err := rows2.Scan(&owner, &voiceId); err != nil {
			return err
		}
		ownerByVoice[voiceId] = owner
	}

	for _, voiceId := range uniq {
		if owner, ok := ownerByVoice[voiceId]; ok && owner != userId {
			return fmt.Errorf("%w: voice %d user %d", ErrVoiceForbidden, voiceId, userId)
		}
	}

	return nil
}
