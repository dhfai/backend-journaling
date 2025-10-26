package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewPostgres(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil
}

func RunMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`,
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email TEXT UNIQUE NOT NULL,
			username TEXT UNIQUE,
			password_hash TEXT,
			is_active BOOLEAN DEFAULT FALSE,
			is_verified BOOLEAN DEFAULT FALSE,
			role TEXT DEFAULT 'user',
			created_at TIMESTAMPTZ DEFAULT now(),
			updated_at TIMESTAMPTZ DEFAULT now()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`,
		`CREATE TABLE IF NOT EXISTS otps (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			purpose TEXT NOT NULL,
			otp_hash TEXT NOT NULL,
			created_at TIMESTAMPTZ DEFAULT now(),
			expires_at TIMESTAMPTZ NOT NULL,
			attempts INT DEFAULT 0,
			consumed BOOLEAN DEFAULT FALSE,
			sent_ip INET,
			sent_user_agent TEXT
		)`,
		`CREATE INDEX IF NOT EXISTS idx_otps_user_id ON otps(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_otps_expires_at ON otps(expires_at)`,
		`CREATE TABLE IF NOT EXISTS refresh_tokens (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			token_hash TEXT NOT NULL,
			created_at TIMESTAMPTZ DEFAULT now(),
			expires_at TIMESTAMPTZ NOT NULL,
			revoked BOOLEAN DEFAULT FALSE,
			replaced_by_uuid UUID
		)`,
		`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash ON refresh_tokens(token_hash)`,
		`CREATE TABLE IF NOT EXISTS auth_events (
			id BIGSERIAL PRIMARY KEY,
			user_id UUID,
			event_type TEXT,
			ip INET,
			user_agent TEXT,
			created_at TIMESTAMPTZ DEFAULT now(),
			meta JSONB
		)`,
		`CREATE INDEX IF NOT EXISTS idx_auth_events_user_id ON auth_events(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_auth_events_created_at ON auth_events(created_at)`,
		`CREATE TABLE IF NOT EXISTS profiles (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			full_name TEXT,
			bio TEXT,
			avatar TEXT,
			date_of_birth DATE,
			gender TEXT,
			phone_number TEXT,
			country TEXT,
			city TEXT,
			timezone TEXT,
			language TEXT DEFAULT 'en',
			theme TEXT DEFAULT 'light',
			created_at TIMESTAMPTZ DEFAULT now(),
			updated_at TIMESTAMPTZ DEFAULT now()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_profiles_user_id ON profiles(user_id)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}
