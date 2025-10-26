package repository

import (
	"database/sql"
	"encoding/json"

	"backend-journaling/internal/models"

	"github.com/google/uuid"
)

type AuthEventRepository struct {
	db *sql.DB
}

func NewAuthEventRepository(db *sql.DB) *AuthEventRepository {
	return &AuthEventRepository{db: db}
}

func (r *AuthEventRepository) Create(userID *uuid.UUID, eventType string, ip, userAgent *string, meta map[string]interface{}) error {
	var metaJSON *string
	if meta != nil {
		b, err := json.Marshal(meta)
		if err != nil {
			return err
		}
		s := string(b)
		metaJSON = &s
	}

	query := `
		INSERT INTO auth_events (user_id, event_type, ip, user_agent, created_at, meta)
		VALUES ($1, $2, $3, $4, now(), $5)
	`
	_, err := r.db.Exec(query, userID, eventType, ip, userAgent, metaJSON)
	return err
}

func (r *AuthEventRepository) FindByUserID(userID uuid.UUID, limit int) ([]models.AuthEvent, error) {
	query := `
		SELECT id, user_id, event_type, ip, user_agent, created_at, meta
		FROM auth_events
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := r.db.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []models.AuthEvent{}
	for rows.Next() {
		var event models.AuthEvent
		err := rows.Scan(
			&event.ID, &event.UserID, &event.EventType,
			&event.IP, &event.UserAgent, &event.CreatedAt, &event.Meta,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}
