package handlers

import (
	"net/http"

	"backend-journaling/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

type RegisterRequest struct {
	Email    string  `json:"email" validate:"required,email"`
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=8"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ip := GetIPAddress(r)
	userAgent := GetUserAgent(r)

	if err := h.service.Register(req.Email, req.Username, req.Password, ip, userAgent); err != nil {
		println("Registration error:", err.Error())
		WriteError(w, http.StatusInternalServerError, "Failed to process registration")
		return
	}

	WriteSuccess(w, http.StatusOK, "Registration initiated. Please check your email for OTP.", nil)
}

type VerifyOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required,len=6"`
}

func (h *AuthHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req VerifyOTPRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ip := GetIPAddress(r)
	userAgent := GetUserAgent(r)

	resp, err := h.service.VerifyOTP(req.Email, req.OTP, ip, userAgent)
	if err != nil {
		status := http.StatusBadRequest
		message := err.Error()

		switch err {
		case service.ErrOTPExpired:
			message = "OTP has expired. Please request a new one."
		case service.ErrOTPMaxAttempts:
			message = "Maximum OTP attempts exceeded. Please request a new one."
		case service.ErrInvalidOTP:
			message = "Invalid OTP code."
		}

		WriteError(w, status, message)
		return
	}

	WriteSuccess(w, http.StatusOK, "Email verified successfully", resp)
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ip := GetIPAddress(r)
	userAgent := GetUserAgent(r)

	resp, err := h.service.Login(req.Email, req.Password, ip, userAgent)
	if err != nil {
		status := http.StatusUnauthorized
		message := "Invalid email or password"

		switch err {
		case service.ErrAccountNotVerified:
			message = "Account not verified. Please verify your email first."
		case service.ErrAccountNotActive:
			message = "Account is not active."
		}

		WriteError(w, status, message)
		return
	}

	WriteSuccess(w, http.StatusOK, "Login successful", resp)
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ip := GetIPAddress(r)
	userAgent := GetUserAgent(r)

	resp, err := h.service.Refresh(req.RefreshToken, ip, userAgent)
	if err != nil {
		status := http.StatusUnauthorized
		message := "Invalid or expired refresh token"

		WriteError(w, status, message)
		return
	}

	WriteSuccess(w, http.StatusOK, "Token refreshed successfully", resp)
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req LogoutRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ip := GetIPAddress(r)
	userAgent := GetUserAgent(r)

	h.service.Logout(req.RefreshToken, ip, userAgent)

	WriteSuccess(w, http.StatusOK, "Logout successful", nil)
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ip := GetIPAddress(r)
	userAgent := GetUserAgent(r)

	h.service.ForgotPassword(req.Email, ip, userAgent)

	WriteSuccess(w, http.StatusOK, "If the email exists, a reset code has been sent.", nil)
}

type ResetPasswordRequest struct {
	Email       string `json:"email" validate:"required,email"`
	OTP         string `json:"otp" validate:"required,len=6"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ip := GetIPAddress(r)
	userAgent := GetUserAgent(r)

	err := h.service.ResetPassword(req.Email, req.OTP, req.NewPassword, ip, userAgent)
	if err != nil {
		status := http.StatusBadRequest
		message := err.Error()

		switch err {
		case service.ErrOTPExpired:
			message = "OTP has expired. Please request a new one."
		case service.ErrOTPMaxAttempts:
			message = "Maximum OTP attempts exceeded. Please request a new one."
		case service.ErrInvalidOTP:
			message = "Invalid OTP code."
		}

		WriteError(w, status, message)
		return
	}

	WriteSuccess(w, http.StatusOK, "Password reset successfully", nil)
}

type RequestOTPRequest struct {
	Email   string `json:"email" validate:"required,email"`
	Purpose string `json:"purpose" validate:"required,oneof=register reset_password"`
}

func (h *AuthHandler) RequestOTP(w http.ResponseWriter, r *http.Request) {
	var req RequestOTPRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ip := GetIPAddress(r)
	userAgent := GetUserAgent(r)

	h.service.RequestOTP(req.Email, req.Purpose, ip, userAgent)

	WriteSuccess(w, http.StatusOK, "OTP has been sent if the account exists.", nil)
}
