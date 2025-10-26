package repository

import (
	"database/sql"
	"errors"

	"backend-journaling/internal/models"

	"github.com/google/uuid"
)

var ErrProfileNotFound = errors.New("profile not found")

type ProfileRepository struct {
	db *sql.DB
}

func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) Create(userID uuid.UUID) (*models.Profile, error) {
	profile := &models.Profile{}
	query := `
		INSERT INTO profiles (user_id, language, theme, created_at, updated_at)
		VALUES ($1, 'en', 'light', now(), now())
		RETURNING id, user_id, full_name, bio, avatar, date_of_birth, gender,
		          phone_number, country, city, timezone, language, theme, created_at, updated_at
	`
	err := r.db.QueryRow(query, userID).Scan(
		&profile.ID, &profile.UserID, &profile.FullName, &profile.Bio,
		&profile.Avatar, &profile.DateOfBirth, &profile.Gender, &profile.PhoneNumber,
		&profile.Country, &profile.City, &profile.Timezone, &profile.Language,
		&profile.Theme, &profile.CreatedAt, &profile.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (r *ProfileRepository) FindByUserID(userID uuid.UUID) (*models.Profile, error) {
	profile := &models.Profile{}
	query := `
		SELECT id, user_id, full_name, bio, avatar, date_of_birth, gender,
		       phone_number, country, city, timezone, language, theme, created_at, updated_at
		FROM profiles WHERE user_id = $1
	`
	err := r.db.QueryRow(query, userID).Scan(
		&profile.ID, &profile.UserID, &profile.FullName, &profile.Bio,
		&profile.Avatar, &profile.DateOfBirth, &profile.Gender, &profile.PhoneNumber,
		&profile.Country, &profile.City, &profile.Timezone, &profile.Language,
		&profile.Theme, &profile.CreatedAt, &profile.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrProfileNotFound
	}
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (r *ProfileRepository) Update(userID uuid.UUID, profile *models.Profile) error {
	query := `
		UPDATE profiles SET
			full_name = $1,
			bio = $2,
			avatar = $3,
			date_of_birth = $4,
			gender = $5,
			phone_number = $6,
			country = $7,
			city = $8,
			timezone = $9,
			language = $10,
			theme = $11,
			updated_at = now()
		WHERE user_id = $12
	`
	_, err := r.db.Exec(query,
		profile.FullName, profile.Bio, profile.Avatar, profile.DateOfBirth,
		profile.Gender, profile.PhoneNumber, profile.Country, profile.City,
		profile.Timezone, profile.Language, profile.Theme, userID,
	)
	return err
}

func (r *ProfileRepository) GetOrCreate(userID uuid.UUID) (*models.Profile, error) {
	profile, err := r.FindByUserID(userID)
	if err == ErrProfileNotFound {
		return r.Create(userID)
	}
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (r *ProfileRepository) UpdateAvatar(userID uuid.UUID, avatarURL string) error {
	query := `UPDATE profiles SET avatar = $1, updated_at = now() WHERE user_id = $2`
	_, err := r.db.Exec(query, avatarURL, userID)
	return err
}

func (r *ProfileRepository) UpdateTheme(userID uuid.UUID, theme string) error {
	query := `UPDATE profiles SET theme = $1, updated_at = now() WHERE user_id = $2`
	_, err := r.db.Exec(query, theme, userID)
	return err
}

func (r *ProfileRepository) UpdateLanguage(userID uuid.UUID, language string) error {
	query := `UPDATE profiles SET language = $1, updated_at = now() WHERE user_id = $2`
	_, err := r.db.Exec(query, language, userID)
	return err
}
