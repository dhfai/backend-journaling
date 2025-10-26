package service

import (
	"errors"
	"time"

	"backend-journaling/config"
	"backend-journaling/internal/models"
	"backend-journaling/internal/repository"
	"backend-journaling/pkg/email"
	"backend-journaling/pkg/jwt"
	"backend-journaling/pkg/otp"
	"backend-journaling/pkg/password"
	"backend-journaling/pkg/token"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountNotActive   = errors.New("account not active")
	ErrAccountNotVerified = errors.New("account not verified")
	ErrOTPExpired         = errors.New("otp has expired")
	ErrOTPConsumed        = errors.New("otp already consumed")
	ErrOTPMaxAttempts     = errors.New("maximum otp attempts exceeded")
	ErrInvalidOTP         = errors.New("invalid otp")
	ErrInvalidToken       = errors.New("invalid refresh token")
	ErrTokenExpired       = errors.New("refresh token expired")
	ErrTokenRevoked       = errors.New("refresh token revoked")
)

type AuthService struct {
	userRepo         *repository.UserRepository
	otpRepo          *repository.OTPRepository
	refreshTokenRepo *repository.RefreshTokenRepository
	authEventRepo    *repository.AuthEventRepository
	jwtManager       *jwt.Manager
	emailSender      email.Sender
	config           *config.Config
}

func NewAuthService(
	userRepo *repository.UserRepository,
	otpRepo *repository.OTPRepository,
	refreshTokenRepo *repository.RefreshTokenRepository,
	authEventRepo *repository.AuthEventRepository,
	jwtManager *jwt.Manager,
	emailSender email.Sender,
	config *config.Config,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		otpRepo:          otpRepo,
		refreshTokenRepo: refreshTokenRepo,
		authEventRepo:    authEventRepo,
		jwtManager:       jwtManager,
		emailSender:      emailSender,
		config:           config,
	}
}

func (s *AuthService) Register(email string, username *string, pass *string, ip, userAgent *string) error {
	var passwordHash *string
	if pass != nil {
		hash, err := password.Hash(*pass)
		if err != nil {
			return err
		}
		passwordHash = &hash
	}

	user, err := s.userRepo.Upsert(email, username, passwordHash)
	if err != nil {
		return err
	}

	otpCode, err := otp.Generate()
	if err != nil {
		return err
	}

	otpHash := otp.Hash(otpCode, s.config.OTP.Pepper)
	expiresAt := time.Now().Add(s.config.OTP.TTL)

	_, err = s.otpRepo.Create(user.ID, "register", otpHash, expiresAt, ip, userAgent)
	if err != nil {
		return err
	}

	if err := s.emailSender.SendOTP(email, otpCode, "register"); err != nil {
		return err
	}

	s.authEventRepo.Create(&user.ID, "otp_sent", ip, userAgent, map[string]interface{}{
		"purpose": "register",
	})

	return nil
}

func (s *AuthService) VerifyOTP(email, otpCode string, ip, userAgent *string) (*LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, ErrInvalidOTP
	}

	otpRecord, err := s.otpRepo.FindLatest(user.ID, "register")
	if err != nil {
		return nil, ErrInvalidOTP
	}

	if otpRecord.Consumed {
		return nil, ErrOTPConsumed
	}

	if time.Now().After(otpRecord.ExpiresAt) {
		s.otpRepo.MarkConsumed(otpRecord.ID)
		return nil, ErrOTPExpired
	}

	if otpRecord.Attempts >= s.config.OTP.MaxAttempts {
		s.otpRepo.MarkConsumed(otpRecord.ID)
		return nil, ErrOTPMaxAttempts
	}

	if !otp.Verify(otpCode, s.config.OTP.Pepper, otpRecord.OTPHash) {
		s.otpRepo.IncrementAttempts(otpRecord.ID)
		s.authEventRepo.Create(&user.ID, "otp_failed", ip, userAgent, nil)
		return nil, ErrInvalidOTP
	}

	s.otpRepo.MarkConsumed(otpRecord.ID)
	s.userRepo.Activate(user.ID)

	s.authEventRepo.Create(&user.ID, "account_verified", ip, userAgent, nil)

	return s.generateTokens(user, ip, userAgent)
}

func (s *AuthService) Login(email, pass string, ip, userAgent *string) (*LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsVerified {
		return nil, ErrAccountNotVerified
	}

	if !user.IsActive {
		return nil, ErrAccountNotActive
	}

	if user.PasswordHash == nil {
		return nil, ErrInvalidCredentials
	}

	valid, err := password.Verify(pass, *user.PasswordHash)
	if err != nil || !valid {
		s.authEventRepo.Create(&user.ID, "login_failed", ip, userAgent, nil)
		return nil, ErrInvalidCredentials
	}

	s.authEventRepo.Create(&user.ID, "login", ip, userAgent, nil)

	return s.generateTokens(user, ip, userAgent)
}

func (s *AuthService) Refresh(refreshToken string, ip, userAgent *string) (*LoginResponse, error) {
	tokenHash := token.Hash(refreshToken)

	tokenRecord, err := s.refreshTokenRepo.FindByTokenHash(tokenHash)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if tokenRecord.Revoked {
		return nil, ErrTokenRevoked
	}

	if time.Now().After(tokenRecord.ExpiresAt) {
		return nil, ErrTokenExpired
	}

	user, err := s.userRepo.FindByID(tokenRecord.UserID)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := token.Generate()
	if err != nil {
		return nil, err
	}

	newTokenHash := token.Hash(newRefreshToken)
	expiresAt := time.Now().Add(s.config.JWT.RefreshTokenDuration)

	newTokenRecord, err := s.refreshTokenRepo.Create(user.ID, newTokenHash, expiresAt)
	if err != nil {
		return nil, err
	}

	s.refreshTokenRepo.Revoke(tokenRecord.ID, &newTokenRecord.ID)

	accessToken, err := s.jwtManager.Generate(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	s.authEventRepo.Create(&user.ID, "token_refreshed", ip, userAgent, nil)

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int(s.jwtManager.GetDuration().Seconds()),
		User:         user,
	}, nil
}

func (s *AuthService) Logout(refreshToken string, ip, userAgent *string) error {
	tokenHash := token.Hash(refreshToken)

	tokenRecord, err := s.refreshTokenRepo.FindByTokenHash(tokenHash)
	if err != nil {
		return nil
	}

	s.refreshTokenRepo.Revoke(tokenRecord.ID, nil)

	s.authEventRepo.Create(&tokenRecord.UserID, "logout", ip, userAgent, nil)

	return nil
}

func (s *AuthService) ForgotPassword(email string, ip, userAgent *string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil
	}

	otpCode, err := otp.Generate()
	if err != nil {
		return err
	}

	otpHash := otp.Hash(otpCode, s.config.OTP.Pepper)
	expiresAt := time.Now().Add(s.config.OTP.TTL)

	_, err = s.otpRepo.Create(user.ID, "reset_password", otpHash, expiresAt, ip, userAgent)
	if err != nil {
		return err
	}

	if err := s.emailSender.SendOTP(email, otpCode, "reset_password"); err != nil {
		return err
	}

	s.authEventRepo.Create(&user.ID, "otp_sent", ip, userAgent, map[string]interface{}{
		"purpose": "reset_password",
	})

	return nil
}

func (s *AuthService) ResetPassword(email, otpCode, newPassword string, ip, userAgent *string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return ErrInvalidOTP
	}

	otpRecord, err := s.otpRepo.FindLatest(user.ID, "reset_password")
	if err != nil {
		return ErrInvalidOTP
	}

	if otpRecord.Consumed {
		return ErrOTPConsumed
	}

	if time.Now().After(otpRecord.ExpiresAt) {
		s.otpRepo.MarkConsumed(otpRecord.ID)
		return ErrOTPExpired
	}

	if otpRecord.Attempts >= s.config.OTP.MaxAttempts {
		s.otpRepo.MarkConsumed(otpRecord.ID)
		return ErrOTPMaxAttempts
	}

	if !otp.Verify(otpCode, s.config.OTP.Pepper, otpRecord.OTPHash) {
		s.otpRepo.IncrementAttempts(otpRecord.ID)
		s.authEventRepo.Create(&user.ID, "otp_failed", ip, userAgent, nil)
		return ErrInvalidOTP
	}

	passwordHash, err := password.Hash(newPassword)
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdatePassword(user.ID, passwordHash); err != nil {
		return err
	}

	s.otpRepo.MarkConsumed(otpRecord.ID)
	s.refreshTokenRepo.RevokeAllForUser(user.ID)

	s.authEventRepo.Create(&user.ID, "password_reset", ip, userAgent, nil)

	return nil
}

func (s *AuthService) RequestOTP(email, purpose string, ip, userAgent *string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil
	}

	otpCode, err := otp.Generate()
	if err != nil {
		return err
	}

	otpHash := otp.Hash(otpCode, s.config.OTP.Pepper)
	expiresAt := time.Now().Add(s.config.OTP.TTL)

	_, err = s.otpRepo.Create(user.ID, purpose, otpHash, expiresAt, ip, userAgent)
	if err != nil {
		return err
	}

	if err := s.emailSender.SendOTP(email, otpCode, purpose); err != nil {
		return err
	}

	s.authEventRepo.Create(&user.ID, "otp_sent", ip, userAgent, map[string]interface{}{
		"purpose": purpose,
	})

	return nil
}

func (s *AuthService) ChangePassword(userID uuid.UUID, oldPassword, newPassword string, ip, userAgent *string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if user.PasswordHash == nil {
		return ErrInvalidCredentials
	}

	valid, err := password.Verify(oldPassword, *user.PasswordHash)
	if err != nil || !valid {
		return ErrInvalidCredentials
	}

	newHash, err := password.Hash(newPassword)
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdatePassword(userID, newHash); err != nil {
		return err
	}

	s.refreshTokenRepo.RevokeAllForUser(userID)

	s.authEventRepo.Create(&userID, "password_changed", ip, userAgent, nil)

	return nil
}

func (s *AuthService) generateTokens(user *models.User, ip, userAgent *string) (*LoginResponse, error) {
	accessToken, err := s.jwtManager.Generate(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := token.Generate()
	if err != nil {
		return nil, err
	}

	tokenHash := token.Hash(refreshToken)
	expiresAt := time.Now().Add(s.config.JWT.RefreshTokenDuration)

	_, err = s.refreshTokenRepo.Create(user.ID, tokenHash, expiresAt)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(s.jwtManager.GetDuration().Seconds()),
		User:         user,
	}, nil
}

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int          `json:"expires_in"`
	User         *models.User `json:"user"`
}
