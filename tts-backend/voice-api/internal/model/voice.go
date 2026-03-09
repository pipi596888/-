package model

import (
	"database/sql"
	"time"
)

type Voice struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	Tone       string    `json:"tone"`
	Gender     string    `json:"gender"`
	PreviewUrl string    `json:"previewUrl"`
	IsDefault  bool      `json:"isDefault"`
	CreatedAt  time.Time `json:"createdAt"`
}

type VoiceModel interface {
	Insert(voice *Voice) (int64, error)
	FindOne(id int64) (*Voice, error)
	FindAll() ([]*Voice, error)
	Delete(id int64) error
	SetDefault(id int64) error
	ClearDefault() error
}

type DefaultVoiceModel struct {
	db *sql.DB
}

func NewVoiceModel(db *sql.DB) VoiceModel {
	return &DefaultVoiceModel{db: db}
}

func (m *DefaultVoiceModel) Insert(voice *Voice) (int64, error) {
	query := `INSERT INTO voice (name, tone, gender, preview_url, is_default) VALUES (?, ?, ?, ?, ?)`
	result, err := m.db.Exec(query, voice.Name, voice.Tone, voice.Gender, voice.PreviewUrl, voice.IsDefault)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (m *DefaultVoiceModel) FindOne(id int64) (*Voice, error) {
	query := `SELECT id, name, tone, gender, COALESCE(preview_url,'') as preview_url, is_default, created_at FROM voice WHERE id = ?`
	var voice Voice
	err := m.db.QueryRow(query, id).Scan(
		&voice.Id, &voice.Name, &voice.Tone, &voice.Gender,
		&voice.PreviewUrl, &voice.IsDefault, &voice.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &voice, nil
}

func (m *DefaultVoiceModel) FindAll() ([]*Voice, error) {
	query := `SELECT id, name, tone, gender, COALESCE(preview_url,'') as preview_url, is_default, created_at FROM voice ORDER BY id DESC`
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var voices []*Voice
	for rows.Next() {
		var voice Voice
		err := rows.Scan(
			&voice.Id, &voice.Name, &voice.Tone, &voice.Gender,
			&voice.PreviewUrl, &voice.IsDefault, &voice.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		voices = append(voices, &voice)
	}
	return voices, nil
}

func (m *DefaultVoiceModel) Delete(id int64) error {
	query := `DELETE FROM voice WHERE id = ?`
	_, err := m.db.Exec(query, id)
	return err
}

func (m *DefaultVoiceModel) SetDefault(id int64) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`UPDATE voice SET is_default = 0`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE voice SET is_default = 1 WHERE id = ?`, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (m *DefaultVoiceModel) ClearDefault() error {
	_, err := m.db.Exec(`UPDATE voice SET is_default = 0`)
	return err
}
