package model

import (
	"database/sql"
	"time"
)

type TtsTask struct {
	Id         int64     `json:"id"`
	TaskId     string    `json:"taskId"`
	UserId     int64     `json:"userId"`
	Status     string    `json:"status"`
	Progress   int       `json:"progress"`
	AudioUrl   string    `json:"audioUrl"`
	Format     string    `json:"format"`
	Channel    string    `json:"channel"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type TtsSegment struct {
	Id       int64  `json:"id"`
	TaskId   string `json:"taskId"`
	VoiceId  int64  `json:"voiceId"`
	Emotion  string `json:"emotion"`
	Text     string `json:"text"`
	Sort     int    `json:"sort"`
}

type TtsTaskModel interface {
	FindByTaskId(taskId string) (*TtsTask, error)
	UpdateStatus(taskId string, status string, progress int) error
	UpdateAudioUrl(taskId string, audioUrl string) error
	UpdateError(taskId string, errMsg string) error
}

type TtsSegmentModel interface {
	FindByTaskId(taskId string) ([]*TtsSegment, error)
}

type DefaultTtsTaskModel struct {
	db *sql.DB
}

func NewTtsTaskModel(db *sql.DB) TtsTaskModel {
	return &DefaultTtsTaskModel{db: db}
}

func (m *DefaultTtsTaskModel) FindByTaskId(taskId string) (*TtsTask, error) {
	query := `SELECT id, task_id, user_id, status, progress, audio_url, format, channel, created_at, updated_at FROM tts_task WHERE task_id = ?`
	var task TtsTask
	err := m.db.QueryRow(query, taskId).Scan(
		&task.Id, &task.TaskId, &task.UserId, &task.Status,
		&task.Progress, &task.AudioUrl, &task.Format, &task.Channel,
		&task.CreatedAt, &task.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (m *DefaultTtsTaskModel) UpdateStatus(taskId string, status string, progress int) error {
	query := `UPDATE tts_task SET status = ?, progress = ?, updated_at = ? WHERE task_id = ?`
	_, err := m.db.Exec(query, status, progress, time.Now(), taskId)
	return err
}

func (m *DefaultTtsTaskModel) UpdateAudioUrl(taskId string, audioUrl string) error {
	query := `UPDATE tts_task SET audio_url = ?, status = 'success', progress = 100, updated_at = ? WHERE task_id = ?`
	_, err := m.db.Exec(query, audioUrl, time.Now(), taskId)
	return err
}

func (m *DefaultTtsTaskModel) UpdateError(taskId string, errMsg string) error {
	query := `UPDATE tts_task SET status = 'failed', error_msg = ?, updated_at = ? WHERE task_id = ?`
	_, err := m.db.Exec(query, errMsg, time.Now(), taskId)
	return err
}

type DefaultTtsSegmentModel struct {
	db *sql.DB
}

func NewTtsSegmentModel(db *sql.DB) TtsSegmentModel {
	return &DefaultTtsSegmentModel{db: db}
}

func (m *DefaultTtsSegmentModel) FindByTaskId(taskId string) ([]*TtsSegment, error) {
	query := `SELECT id, task_id, voice_id, emotion, text, sort FROM tts_segment WHERE task_id = ? ORDER BY sort`
	rows, err := m.db.Query(query, taskId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var segments []*TtsSegment
	for rows.Next() {
		var seg TtsSegment
		err := rows.Scan(&seg.Id, &seg.TaskId, &seg.VoiceId, &seg.Emotion, &seg.Text, &seg.Sort)
		if err != nil {
			return nil, err
		}
		segments = append(segments, &seg)
	}
	return segments, nil
}
