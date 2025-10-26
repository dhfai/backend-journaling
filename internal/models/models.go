package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Username     *string   `json:"username,omitempty" db:"username"`
	PasswordHash *string   `json:"-" db:"password_hash"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	IsVerified   bool      `json:"is_verified" db:"is_verified"`
	Role         string    `json:"role" db:"role"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type OTP struct {
	ID            uuid.UUID `json:"id" db:"id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	Purpose       string    `json:"purpose" db:"purpose"`
	OTPHash       string    `json:"-" db:"otp_hash"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	ExpiresAt     time.Time `json:"expires_at" db:"expires_at"`
	Attempts      int       `json:"attempts" db:"attempts"`
	Consumed      bool      `json:"consumed" db:"consumed"`
	SentIP        *string   `json:"sent_ip,omitempty" db:"sent_ip"`
	SentUserAgent *string   `json:"sent_user_agent,omitempty" db:"sent_user_agent"`
}

type RefreshToken struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	UserID         uuid.UUID  `json:"user_id" db:"user_id"`
	TokenHash      string     `json:"-" db:"token_hash"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	ExpiresAt      time.Time  `json:"expires_at" db:"expires_at"`
	Revoked        bool       `json:"revoked" db:"revoked"`
	ReplacedByUUID *uuid.UUID `json:"replaced_by_uuid,omitempty" db:"replaced_by_uuid"`
}

type AuthEvent struct {
	ID        int64      `json:"id" db:"id"`
	UserID    *uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	EventType string     `json:"event_type" db:"event_type"`
	IP        *string    `json:"ip,omitempty" db:"ip"`
	UserAgent *string    `json:"user_agent,omitempty" db:"user_agent"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	Meta      *string    `json:"meta,omitempty" db:"meta"`
}

type Profile struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	FullName    *string    `json:"full_name,omitempty" db:"full_name"`
	Bio         *string    `json:"bio,omitempty" db:"bio"`
	Avatar      *string    `json:"avatar,omitempty" db:"avatar"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty" db:"date_of_birth"`
	Gender      *string    `json:"gender,omitempty" db:"gender"`
	PhoneNumber *string    `json:"phone_number,omitempty" db:"phone_number"`
	Country     *string    `json:"country,omitempty" db:"country"`
	City        *string    `json:"city,omitempty" db:"city"`
	Timezone    *string    `json:"timezone,omitempty" db:"timezone"`
	Language    string     `json:"language" db:"language"`
	Theme       string     `json:"theme" db:"theme"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type UserWithProfile struct {
	User
	Profile *Profile `json:"profile,omitempty"`
}
