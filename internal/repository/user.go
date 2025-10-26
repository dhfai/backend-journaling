package repository

import (
	"database/sql"
	"errors"

	"backend-journaling/internal/models"

	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(email string, username *string, passwordHash *string) (*models.User, error) {
	user := &models.User{}
	query := `
		INSERT INTO users (email, username, password_hash, is_active, is_verified, role, created_at, updated_at)
		VALUES ($1, $2, $3, false, false, 'user', now(), now())
		RETURNING id, email, username, password_hash, is_active, is_verified, role, created_at, updated_at
	`
	err := r.db.QueryRow(query, email, username, passwordHash).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.IsActive, &user.IsVerified, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, username, password_hash, is_active, is_verified, role, created_at, updated_at
		FROM users WHERE email = $1
	`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.IsActive, &user.IsVerified, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, username, password_hash, is_active, is_verified, role, created_at, updated_at
		FROM users WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.IsActive, &user.IsVerified, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Update(id uuid.UUID, username *string) error {
	query := `UPDATE users SET username = $1, updated_at = now() WHERE id = $2`
	_, err := r.db.Exec(query, username, id)
	return err
}

func (r *UserRepository) UpdatePassword(id uuid.UUID, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = now() WHERE id = $2`
	_, err := r.db.Exec(query, passwordHash, id)
	return err
}

func (r *UserRepository) Activate(id uuid.UUID) error {
	query := `UPDATE users SET is_active = true, is_verified = true, updated_at = now() WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *UserRepository) Upsert(email string, username *string, passwordHash *string) (*models.User, error) {
	user, err := r.FindByEmail(email)
	if err == ErrUserNotFound {
		return r.Create(email, username, passwordHash)
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
