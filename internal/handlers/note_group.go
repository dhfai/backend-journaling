package handlers

import (
	"encoding/json"
	"net/http"

	"backend-journaling/internal/service"
	"backend-journaling/pkg/jwt"

	"github.com/go-chi/chi/v5"
)

type NoteGroupHandler struct {
	service *service.NoteGroupService
}

func NewNoteGroupHandler(service *service.NoteGroupService) *NoteGroupHandler {
	return &NoteGroupHandler{
		service: service,
	}
}

type CreateGroupRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
	Icon        *string `json:"icon,omitempty"`
}

type UpdateGroupRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
	Icon        *string `json:"icon,omitempty"`
}

type PinGroupRequest struct {
	IsPinned bool `json:"is_pinned"`
}

type ArchiveGroupRequest struct {
	IsArchived bool `json:"is_archived"`
}

type AddNoteToGroupRequest struct {
	NoteID string `json:"note_id"`
}

type RemoveNoteFromGroupRequest struct {
	NoteID string `json:"note_id"`
}

type MoveNotesRequest struct {
	NoteIDs []string `json:"note_ids"`
}

// CreateGroup creates a new note group
func (h *NoteGroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)

	var req CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		WriteError(w, http.StatusBadRequest, "Group name is required")
		return
	}

	group, err := h.service.CreateGroup(r.Context(), claims.UserID.String(), req.Name, req.Description, req.Color, req.Icon)
	if err != nil {
		if err == service.ErrGroupNameRequired {
			WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to create group")
		return
	}

	WriteJSON(w, http.StatusCreated, group)
}

// GetGroups gets all note groups for the user
func (h *NoteGroupHandler) GetGroups(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)

	// Optional query parameters
	var isPinned, isArchived *bool

	if r.URL.Query().Get("is_pinned") != "" {
		val := r.URL.Query().Get("is_pinned") == "true"
		isPinned = &val
	}

	if r.URL.Query().Get("is_archived") != "" {
		val := r.URL.Query().Get("is_archived") == "true"
		isArchived = &val
	}

	groups, err := h.service.GetGroups(r.Context(), claims.UserID.String(), isPinned, isArchived)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to fetch groups")
		return
	}

	WriteJSON(w, http.StatusOK, groups)
}

// GetGroup gets a single note group
func (h *NoteGroupHandler) GetGroup(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)
	groupID := chi.URLParam(r, "id")

	group, err := h.service.GetGroup(r.Context(), groupID, claims.UserID.String())
	if err != nil {
		if err == service.ErrGroupNotFound || err == service.ErrInvalidGroupID {
			WriteError(w, http.StatusNotFound, "Group not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to fetch group")
		return
	}

	WriteJSON(w, http.StatusOK, group)
}

// UpdateGroup updates a note group
func (h *NoteGroupHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)
	groupID := chi.URLParam(r, "id")

	var req UpdateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	group, err := h.service.UpdateGroup(r.Context(), groupID, claims.UserID.String(), req.Name, req.Description, req.Color, req.Icon)
	if err != nil {
		if err == service.ErrGroupNotFound || err == service.ErrInvalidGroupID {
			WriteError(w, http.StatusNotFound, "Group not found")
			return
		}
		if err == service.ErrGroupNameRequired {
			WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to update group")
		return
	}

	WriteJSON(w, http.StatusOK, group)
}

// PinGroup pins/unpins a note group
func (h *NoteGroupHandler) PinGroup(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)
	groupID := chi.URLParam(r, "id")

	var req PinGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.PinGroup(r.Context(), groupID, claims.UserID.String(), req.IsPinned); err != nil {
		if err == service.ErrGroupNotFound || err == service.ErrInvalidGroupID {
			WriteError(w, http.StatusNotFound, "Group not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to pin group")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message":   "Group pin status updated",
		"is_pinned": req.IsPinned,
	})
}

// ArchiveGroup archives/unarchives a note group
func (h *NoteGroupHandler) ArchiveGroup(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)
	groupID := chi.URLParam(r, "id")

	var req ArchiveGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.ArchiveGroup(r.Context(), groupID, claims.UserID.String(), req.IsArchived); err != nil {
		if err == service.ErrGroupNotFound || err == service.ErrInvalidGroupID {
			WriteError(w, http.StatusNotFound, "Group not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to archive group")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message":     "Group archive status updated",
		"is_archived": req.IsArchived,
	})
}

// DeleteGroup deletes a note group and removes group_id from all notes
func (h *NoteGroupHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)
	groupID := chi.URLParam(r, "id")

	// Delete group (will also remove group_id from all notes)
	err := h.service.DeleteGroup(r.Context(), groupID, claims.UserID.String())
	if err != nil {
		if err == service.ErrGroupNotFound || err == service.ErrInvalidGroupID {
			WriteError(w, http.StatusNotFound, "Group not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to delete group")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Group deleted successfully",
	})
}

// AddNoteToGroup adds a note to a group
func (h *NoteGroupHandler) AddNoteToGroup(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)
	groupID := chi.URLParam(r, "id")

	var req AddNoteToGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.NoteID == "" {
		WriteError(w, http.StatusBadRequest, "Note ID is required")
		return
	}

	if err := h.service.AddNoteToGroup(r.Context(), req.NoteID, groupID, claims.UserID.String()); err != nil {
		if err == service.ErrGroupNotFound || err == service.ErrInvalidGroupID {
			WriteError(w, http.StatusNotFound, "Group not found")
			return
		}
		if err == service.ErrNoteNotFound || err == service.ErrInvalidNoteID {
			WriteError(w, http.StatusNotFound, "Note not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to add note to group")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Note added to group successfully",
	})
}

// RemoveNoteFromGroup removes a note from a group
func (h *NoteGroupHandler) RemoveNoteFromGroup(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)
	noteID := chi.URLParam(r, "noteId")

	if err := h.service.RemoveNoteFromGroup(r.Context(), noteID, claims.UserID.String()); err != nil {
		if err == service.ErrNoteNotFound || err == service.ErrInvalidNoteID {
			WriteError(w, http.StatusNotFound, "Note not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to remove note from group")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Note removed from group successfully",
	})
}

// MoveNotesToGroup moves multiple notes to a group
func (h *NoteGroupHandler) MoveNotesToGroup(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)
	groupID := chi.URLParam(r, "id")

	var req MoveNotesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(req.NoteIDs) == 0 {
		WriteError(w, http.StatusBadRequest, "Note IDs are required")
		return
	}

	if err := h.service.MoveNotesToGroup(r.Context(), req.NoteIDs, groupID, claims.UserID.String()); err != nil {
		if err == service.ErrGroupNotFound || err == service.ErrInvalidGroupID {
			WriteError(w, http.StatusNotFound, "Group not found")
			return
		}
		if err == service.ErrInvalidNoteID {
			WriteError(w, http.StatusBadRequest, "Invalid note ID")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to move notes to group")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Notes moved to group successfully",
	})
}

// GetNotesInGroup gets all notes in a group
func (h *NoteGroupHandler) GetNotesInGroup(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)
	groupID := chi.URLParam(r, "id")

	notes, err := h.service.GetNotesInGroup(r.Context(), groupID, claims.UserID.String())
	if err != nil {
		if err == service.ErrGroupNotFound || err == service.ErrInvalidGroupID {
			WriteError(w, http.StatusNotFound, "Group not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to fetch notes in group")
		return
	}

	WriteJSON(w, http.StatusOK, notes)
}
