package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Hash      string    `json:"hash"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserData struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) CreateUser(username, email, password string) error {
	// 1 Hash the plaintext password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("ERROR: m.UserModel.CreateUser - bcrypt.GeneratePassword: %v", err)
		return fmt.Errorf("failed to hash password: %w", err)
	}

	stmt := `
		INSERT INTO users (username, email, hash, created_at, updated_at) 
		VALUES (?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())
	`
	_, err = m.DB.Exec(stmt, username, email, hashPassword)
	if err != nil {
		if IsDuplicateEntryError(err) {
			return ErrDuplicateRecord
		}
		log.Printf("ERROR: m.UserModel.CreateUser - m.DB.Exec: %v", err)
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

func (m *UserModel) AuthenticateByEmail(email, password string) (*UserData, error) {
	// 1. Retrieve the user data
	var user *User = &User{}

	query := `SELECT id, username, email, hash, created_at, updated_at FROM users WHERE email = ?`

	row := m.DB.QueryRow(query, email)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Hash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		log.Printf("ERROR: m.UserModel.Authenticate - m.QueryRow: %v", err)
		return nil, fmt.Errorf("invalid authentication credentials: %w", err)
	}

	// 2. Compare the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidCredentials
		}
		log.Printf("ERROR: m.UserModel.Authenticate - bcrypt.CompareHashAndPassword: %v", err)
		return nil, fmt.Errorf("failed to compare hash and password: %w", err)
	}

	userData := &UserData{
		ID: user.ID,
		Username: user.Username,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return userData, nil
}

func (m *UserModel) AuthenticateByUsername(username, password string) (*UserData, error) {
	// 1. Retrieve the user data
	var user *User = &User{}

	query := `SELECT id, username, email, hash, created_at, updated_at FROM users WHERE username = ?`

	row := m.DB.QueryRow(query, username)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Hash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		log.Printf("ERROR: m.UserModel.Authenticate - m.QueryRow: %v", err)
		return nil, fmt.Errorf("invalid authentication credentials: %w", err)
	}

	// 2. Compare the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidCredentials
		}
		log.Printf("ERROR: m.UserModel.Authenticate - bcrypt.CompareHashAndPassword: %v", err)
		return nil, fmt.Errorf("failed to compare hash and password: %w", err)
	}

	userData := &UserData{
		ID: user.ID,
		Username: user.Username,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return userData, nil
}