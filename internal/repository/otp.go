package repository

import (
	"database/sql"
	"errors"
	"time"

	"backend-journaling/internal/models"

	"github.com/google/uuid"
)

var ErrOTPNotFound = errors.New("otp not found")

type OTPRepository struct {
	db *sql.DB
}

func NewOTPRepository(db *sql.DB) *OTPRepository {
	return &OTPRepository{db: db}
}

func (r *OTPRepository) Create(userID uuid.UUID, purpose, otpHash string, expiresAt time.Time, ip, userAgent *string) (*models.OTP, error) {
	otp := &models.OTP{}
	query := `
		INSERT INTO otps (user_id, purpose, otp_hash, created_at, expires_at, attempts, consumed, sent_ip, sent_user_agent)
		VALUES ($1, $2, $3, now(), $4, 0, false, $5, $6)
		RETURNING id, user_id, purpose, otp_hash, created_at, expires_at, attempts, consumed, sent_ip, sent_user_agent
	`
	err := r.db.QueryRow(query, userID, purpose, otpHash, expiresAt, ip, userAgent).Scan(
		&otp.ID, &otp.UserID, &otp.Purpose, &otp.OTPHash,
		&otp.CreatedAt, &otp.ExpiresAt, &otp.Attempts, &otp.Consumed,
		&otp.SentIP, &otp.SentUserAgent,
	)
	if err != nil {
		return nil, err
	}
	return otp, nil
}

func (r *OTPRepository) FindLatest(userID uuid.UUID, purpose string) (*models.OTP, error) {
	otp := &models.OTP{}
	query := `
		SELECT id, user_id, purpose, otp_hash, created_at, expires_at, attempts, consumed, sent_ip, sent_user_agent
		FROM otps
		WHERE user_id = $1 AND purpose = $2 AND consumed = false
		ORDER BY created_at DESC
		LIMIT 1
	`
	err := r.db.QueryRow(query, userID, purpose).Scan(
		&otp.ID, &otp.UserID, &otp.Purpose, &otp.OTPHash,
		&otp.CreatedAt, &otp.ExpiresAt, &otp.Attempts, &otp.Consumed,
		&otp.SentIP, &otp.SentUserAgent,
	)
	if err == sql.ErrNoRows {
		return nil, ErrOTPNotFound
	}
	if err != nil {
		return nil, err
	}
	return otp, nil
}

func (r *OTPRepository) IncrementAttempts(id uuid.UUID) error {
	query := `UPDATE otps SET attempts = attempts + 1 WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *OTPRepository) MarkConsumed(id uuid.UUID) error {
	query := `UPDATE otps SET consumed = true WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *OTPRepository) DeleteExpired() error {
	query := `DELETE FROM otps WHERE expires_at < now()`
	_, err := r.db.Exec(query)
	return err
}
