package model

import (
	"database/sql"
	"time"
)

type User struct {
	Id              int64     `json:"id"`
	Username        string    `json:"username"`
	Password        string    `json:"-"`
	Email           string    `json:"email"`
	Balance         float64   `json:"balance"`
	CharacterCount  int64     `json:"characterCount"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type UserModel interface {
	Insert(user *User) (int64, error)
	FindOne(id int64) (*User, error)
	FindByUsername(username string) (*User, error)
	UpdateBalance(id int64, balance float64) error
	UpdateCharacterCount(id int64, count int64) error
}

type DefaultUserModel struct {
	db *sql.DB
}

func NewUserModel(db *sql.DB) UserModel {
	return &DefaultUserModel{db: db}
}

func (m *DefaultUserModel) Insert(user *User) (int64, error) {
	query := `INSERT INTO user (username, password, email, balance, character_count) VALUES (?, ?, ?, ?, ?)`
	result, err := m.db.Exec(query, user.Username, user.Password, user.Email, user.Balance, user.CharacterCount)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (m *DefaultUserModel) FindOne(id int64) (*User, error) {
	query := `SELECT id, username, password, email, balance, character_count, created_at, updated_at FROM user WHERE id = ?`
	var user User
	err := m.db.QueryRow(query, id).Scan(
		&user.Id, &user.Username, &user.Password, &user.Emailuser.Balance,,
		& &user.CharacterCount, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *DefaultUserModel) FindByUsername(username string) (*User, error) {
	query := `SELECT id, username, password, email, balance, character_count, created_at, updated_at FROM user WHERE username = ?`
	var user User
	err := m.db.QueryRow(query, username).Scan(
		&user.Id, &user.Username, &user.Password, &user.Email,
		&user.Balance, &user.CharacterCount, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *DefaultUserModel) UpdateBalance(id int64, balance float64) error {
	query := `UPDATE user SET balance = ?, updated_at = ? WHERE id = ?`
	_, err := m.db.Exec(query, balance, time.Now(), id)
	return err
}

func (m *DefaultUserModel) UpdateCharacterCount(id int64, count int64) error {
	query := `UPDATE user SET character_count = character_count + ?, updated_at = ? WHERE id = ?`
	_, err := m.db.Exec(query, count, time.Now(), id)
	return err
}
