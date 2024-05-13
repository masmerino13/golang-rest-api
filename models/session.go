package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"lens.com/m/v2/helpers"
)

const (
	MinBytesPerToken = 32
)

type Session struct {
	ID        int
	UserID    int
	Token     string
	TokenHash string
}

type SessionService struct {
	DB            *sql.DB
	BytesPerToken int
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := ss.BytesPerToken

	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}

	token, err := helpers.String(bytesPerToken)

	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	tokenHash := ss.hash(token)

	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: tokenHash,
	}

	row := ss.DB.QueryRow(`
		UPDATE
		sessions
		SET token_hash = $2
		WHERE user_id = $1
		RETURNING id;
	`, userID, tokenHash)

	err = row.Scan(&session.ID)

	if err == sql.ErrNoRows {
		row = ss.DB.QueryRow("INSERT INTO sessions (user_id, token_hash) VALUES ($1, $2) RETURNING id", userID, tokenHash)
		err = row.Scan(&session.ID)
	}

	if err != nil {
		return nil, fmt.Errorf("error create session: %w", err)
	}

	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	tokenHash := ss.hash(token)
	var user User

	row := ss.DB.QueryRow("SELECT user_id FROM sessions WHERE token_hash = $1", tokenHash)
	err := row.Scan(&user.ID)

	if err != nil {
		return nil, fmt.Errorf("error finding session: %w", err)
	}

	row = ss.DB.QueryRow("SELECT email, password_hash FROM users WHERE id = $1", user.ID)
	err = row.Scan(&user.Email, &user.PasswordHash)

	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	return &user, nil
}

func (ss *SessionService) Delete(token string) error {
	hashToken := ss.hash(token)
	_, err := ss.DB.Exec("DELETE FROM sessions WHERE token_hash = $1", hashToken)

	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))

	// NOTE: [:] this turn the hash in
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
