package handlers

import (
	"errors"
	"net/http"
	"time"

	"backend-journaling/internal/repository"
	"backend-journaling/internal/service"
	"backend-journaling/pkg/jwt"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	userRepo    *repository.UserRepository
	profileRepo *repository.ProfileRepository
	authService *service.AuthService
}

func NewUserHandler(userRepo *repository.UserRepository, profileRepo *repository.ProfileRepository, authService *service.AuthService) *UserHandler {
	return &UserHandler{
		userRepo:    userRepo,
		profileRepo: profileRepo,
		authService: authService,
	}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)

	user, err := h.userRepo.FindByID(claims.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			WriteError(w, http.StatusNotFound, "User not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to fetch user")
		return
	}

	profile, _ := h.profileRepo.GetOrCreate(claims.UserID)

	response := map[string]interface{}{
		"user":    user,
		"profile": profile,
	}

	WriteSuccess(w, http.StatusOK, "Profile retrieved successfully", response)
}

type CreateProfileRequest struct {
	FullName    *string `json:"full_name,omitempty"`
	Bio         *string `json:"bio,omitempty"`
	Avatar      *string `json:"avatar,omitempty"`
	DateOfBirth *string `json:"date_of_birth,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Country     *string `json:"country,omitempty"`
	City        *string `json:"city,omitempty"`
	Timezone    *string `json:"timezone,omitempty"`
	Language    *string `json:"language,omitempty"`
	Theme       *string `json:"theme,omitempty"`
}

func (h *UserHandler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)

	// Check if profile already exists
	existingProfile, _ := h.profileRepo.FindByUserID(claims.UserID)
	if existingProfile != nil {
		WriteError(w, http.StatusConflict, "Profile already exists")
		return
	}

	var req CreateProfileRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Create profile with default values
	profile, err := h.profileRepo.Create(claims.UserID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to create profile")
		return
	}

	// Update with provided values
	if req.FullName != nil {
		profile.FullName = req.FullName
	}
	if req.Bio != nil {
		profile.Bio = req.Bio
	}
	if req.Avatar != nil {
		profile.Avatar = req.Avatar
	}
	if req.DateOfBirth != nil {
		parsedDate, err := time.Parse("2006-01-02", *req.DateOfBirth)
		if err == nil {
			profile.DateOfBirth = &parsedDate
		}
	}
	if req.Gender != nil {
		profile.Gender = req.Gender
	}
	if req.PhoneNumber != nil {
		profile.PhoneNumber = req.PhoneNumber
	}
	if req.Country != nil {
		profile.Country = req.Country
	}
	if req.City != nil {
		profile.City = req.City
	}
	if req.Timezone != nil {
		profile.Timezone = req.Timezone
	}
	if req.Language != nil {
		profile.Language = *req.Language
	}
	if req.Theme != nil {
		profile.Theme = *req.Theme
	}

	// Update profile with new values
	if err := h.profileRepo.Update(claims.UserID, profile); err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	user, _ := h.userRepo.FindByID(claims.UserID)
	updatedProfile, _ := h.profileRepo.FindByUserID(claims.UserID)

	response := map[string]interface{}{
		"user":    user,
		"profile": updatedProfile,
	}

	WriteSuccess(w, http.StatusCreated, "Profile created successfully", response)
}

type UpdateProfileRequest struct {
	FullName    *string `json:"full_name,omitempty"`
	Bio         *string `json:"bio,omitempty"`
	Avatar      *string `json:"avatar,omitempty"`
	DateOfBirth *string `json:"date_of_birth,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Country     *string `json:"country,omitempty"`
	City        *string `json:"city,omitempty"`
	Timezone    *string `json:"timezone,omitempty"`
	Language    *string `json:"language,omitempty"`
	Theme       *string `json:"theme,omitempty"`
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)

	var req UpdateProfileRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	profile, err := h.profileRepo.GetOrCreate(claims.UserID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to get profile")
		return
	}

	if req.FullName != nil {
		profile.FullName = req.FullName
	}
	if req.Bio != nil {
		profile.Bio = req.Bio
	}
	if req.Avatar != nil {
		profile.Avatar = req.Avatar
	}
	if req.DateOfBirth != nil {
		parsedDate, err := time.Parse("2006-01-02", *req.DateOfBirth)
		if err == nil {
			profile.DateOfBirth = &parsedDate
		}
	}
	if req.Gender != nil {
		profile.Gender = req.Gender
	}
	if req.PhoneNumber != nil {
		profile.PhoneNumber = req.PhoneNumber
	}
	if req.Country != nil {
		profile.Country = req.Country
	}
	if req.City != nil {
		profile.City = req.City
	}
	if req.Timezone != nil {
		profile.Timezone = req.Timezone
	}
	if req.Language != nil {
		profile.Language = *req.Language
	}
	if req.Theme != nil {
		profile.Theme = *req.Theme
	}

	if err := h.profileRepo.Update(claims.UserID, profile); err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	user, _ := h.userRepo.FindByID(claims.UserID)
	updatedProfile, _ := h.profileRepo.FindByUserID(claims.UserID)

	response := map[string]interface{}{
		"user":    user,
		"profile": updatedProfile,
	}

	WriteSuccess(w, http.StatusOK, "Profile updated successfully", response)
}

type UpdateAvatarRequest struct {
	AvatarURL string `json:"avatar_url" validate:"required"`
}

func (h *UserHandler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)

	var req UpdateAvatarRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.profileRepo.UpdateAvatar(claims.UserID, req.AvatarURL); err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to update avatar")
		return
	}

	WriteSuccess(w, http.StatusOK, "Avatar updated successfully", nil)
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)

	var req ChangePasswordRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ip := GetIPAddress(r)
	userAgent := GetUserAgent(r)

	err := h.authService.ChangePassword(claims.UserID, req.OldPassword, req.NewPassword, ip, userAgent)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			WriteError(w, http.StatusUnauthorized, "Invalid current password")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to change password")
		return
	}

	WriteSuccess(w, http.StatusOK, "Password changed successfully", nil)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			WriteError(w, http.StatusNotFound, "User not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to fetch user")
		return
	}

	WriteSuccess(w, http.StatusOK, "User retrieved successfully", user)
}
