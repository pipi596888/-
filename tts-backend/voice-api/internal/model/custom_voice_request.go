package model

import (
	"database/sql"
	"time"
)

type CustomVoiceRequest struct {
	Id            int64          `json:"id"`
	UserId        int64          `json:"userId"`
	Name          string         `json:"name"`
	Tone          string         `json:"tone"`
	Gender        string         `json:"gender"`
	SampleText    string         `json:"sampleText"`
	SampleUrls    string         `json:"sampleUrls"`
	Status        string         `json:"status"`
	ResultVoiceId sql.NullInt64  `json:"resultVoiceId"`
	ErrorMsg      sql.NullString `json:"errorMsg"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
}

type CustomVoiceRequestModel interface {
	Insert(req *CustomVoiceRequest) (int64, error)
	FindByUserId(userId int64, limit int) ([]*CustomVoiceRequest, error)
	FindAll(limit int) ([]*CustomVoiceRequest, error)
	FindOne(id int64) (*CustomVoiceRequest, error)
	FindApprovedVoiceOwners() (map[int64]int64, error)
	DeleteByIdAndUserId(id int64, userId int64) error
	Delete(id int64) error
	UpdateApproval(id int64, status string, resultVoiceId sql.NullInt64, errMsg sql.NullString) error
}

type DefaultCustomVoiceRequestModel struct {
	db *sql.DB
}

func NewCustomVoiceRequestModel(db *sql.DB) CustomVoiceRequestModel {
	return &DefaultCustomVoiceRequestModel{db: db}
}

func (m *DefaultCustomVoiceRequestModel) Insert(req *CustomVoiceRequest) (int64, error) {
	query := `INSERT INTO custom_voice_request (user_id, name, tone, gender, sample_text, sample_urls, status) VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := m.db.Exec(query, req.UserId, req.Name, req.Tone, req.Gender, req.SampleText, req.SampleUrls, req.Status)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (m *DefaultCustomVoiceRequestModel) FindByUserId(userId int64, limit int) ([]*CustomVoiceRequest, error) {
	if limit <= 0 {
		limit = 50
	}
	query := `SELECT id, user_id, name, COALESCE(tone,''), COALESCE(gender,''), COALESCE(sample_text,''), COALESCE(sample_urls,''), status, result_voice_id, error_msg, created_at, updated_at
	          FROM custom_voice_request WHERE user_id = ? ORDER BY created_at DESC LIMIT ?`
	rows, err := m.db.Query(query, userId, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*CustomVoiceRequest, 0)
	for rows.Next() {
		var item CustomVoiceRequest
		if err := rows.Scan(
			&item.Id, &item.UserId, &item.Name, &item.Tone, &item.Gender,
			&item.SampleText, &item.SampleUrls, &item.Status, &item.ResultVoiceId, &item.ErrorMsg,
			&item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, &item)
	}
	return list, nil
}

func (m *DefaultCustomVoiceRequestModel) FindAll(limit int) ([]*CustomVoiceRequest, error) {
	if limit <= 0 {
		limit = 200
	}
	query := `SELECT id, user_id, name, COALESCE(tone,''), COALESCE(gender,''), COALESCE(sample_text,''), COALESCE(sample_urls,''), status, result_voice_id, error_msg, created_at, updated_at
	          FROM custom_voice_request ORDER BY created_at DESC LIMIT ?`
	rows, err := m.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*CustomVoiceRequest, 0)
	for rows.Next() {
		var item CustomVoiceRequest
		if err := rows.Scan(
			&item.Id, &item.UserId, &item.Name, &item.Tone, &item.Gender,
			&item.SampleText, &item.SampleUrls, &item.Status, &item.ResultVoiceId, &item.ErrorMsg,
			&item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, &item)
	}
	return list, nil
}

func (m *DefaultCustomVoiceRequestModel) FindOne(id int64) (*CustomVoiceRequest, error) {
	query := `SELECT id, user_id, name, COALESCE(tone,''), COALESCE(gender,''), COALESCE(sample_text,''), COALESCE(sample_urls,''), status, result_voice_id, error_msg, created_at, updated_at
	          FROM custom_voice_request WHERE id = ?`
	var item CustomVoiceRequest
	err := m.db.QueryRow(query, id).Scan(
		&item.Id, &item.UserId, &item.Name, &item.Tone, &item.Gender,
		&item.SampleText, &item.SampleUrls, &item.Status, &item.ResultVoiceId, &item.ErrorMsg,
		&item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (m *DefaultCustomVoiceRequestModel) FindApprovedVoiceOwners() (map[int64]int64, error) {
	query := `SELECT user_id, result_voice_id FROM custom_voice_request WHERE status = 'success' AND result_voice_id IS NOT NULL`
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	owners := make(map[int64]int64)
	for rows.Next() {
		var userId int64
		var voiceId int64
		if err := rows.Scan(&userId, &voiceId); err != nil {
			return nil, err
		}
		owners[voiceId] = userId
	}
	return owners, nil
}

func (m *DefaultCustomVoiceRequestModel) DeleteByIdAndUserId(id int64, userId int64) error {
	_, err := m.db.Exec("DELETE FROM custom_voice_request WHERE id = ? AND user_id = ?", id, userId)
	return err
}

func (m *DefaultCustomVoiceRequestModel) Delete(id int64) error {
	_, err := m.db.Exec("DELETE FROM custom_voice_request WHERE id = ?", id)
	return err
}

func (m *DefaultCustomVoiceRequestModel) UpdateApproval(id int64, status string, resultVoiceId sql.NullInt64, errMsg sql.NullString) error {
	_, err := m.db.Exec(
		"UPDATE custom_voice_request SET status = ?, result_voice_id = ?, error_msg = ? WHERE id = ?",
		status, resultVoiceId, errMsg, id,
	)
	return err
}
