package repository

import (
	"database/sql"
	"errors"
	"time"

	"backend-journaling/internal/models"

	"github.com/google/uuid"
)

var ErrRefreshTokenNotFound = errors.New("refresh token not found")

type RefreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(userID uuid.UUID, tokenHash string, expiresAt time.Time) (*models.RefreshToken, error) {
	token := &models.RefreshToken{}
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, created_at, expires_at, revoked, replaced_by_uuid)
		VALUES ($1, $2, now(), $3, false, NULL)
		RETURNING id, user_id, token_hash, created_at, expires_at, revoked, replaced_by_uuid
	`
	err := r.db.QueryRow(query, userID, tokenHash, expiresAt).Scan(
		&token.ID, &token.UserID, &token.TokenHash, &token.CreatedAt,
		&token.ExpiresAt, &token.Revoked, &token.ReplacedByUUID,
	)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (r *RefreshTokenRepository) FindByTokenHash(tokenHash string) (*models.RefreshToken, error) {
	token := &models.RefreshToken{}
	query := `
		SELECT id, user_id, token_hash, created_at, expires_at, revoked, replaced_by_uuid
		FROM refresh_tokens
		WHERE token_hash = $1
	`
	err := r.db.QueryRow(query, tokenHash).Scan(
		&token.ID, &token.UserID, &token.TokenHash, &token.CreatedAt,
		&token.ExpiresAt, &token.Revoked, &token.ReplacedByUUID,
	)
	if err == sql.ErrNoRows {
		return nil, ErrRefreshTokenNotFound
	}
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (r *RefreshTokenRepository) Revoke(id uuid.UUID, replacedBy *uuid.UUID) error {
	query := `UPDATE refresh_tokens SET revoked = true, replaced_by_uuid = $1 WHERE id = $2`
	_, err := r.db.Exec(query, replacedBy, id)
	return err
}

func (r *RefreshTokenRepository) RevokeAllForUser(userID uuid.UUID) error {
	query := `UPDATE refresh_tokens SET revoked = true WHERE user_id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *RefreshTokenRepository) DeleteExpired() error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < now()`
	_, err := r.db.Exec(query)
	return err
}
