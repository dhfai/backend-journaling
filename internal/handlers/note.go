package handlers

import (
	"encoding/json"
	"net/http"

	"backend-journaling/internal/models"
	"backend-journaling/internal/service"
	"backend-journaling/pkg/jwt"

	"github.com/go-chi/chi/v5"
)

type NoteHandler struct {
	service *service.NoteService
}

func NewNoteHandler(service *service.NoteService) *NoteHandler {
	return &NoteHandler{service: service}
}

type CreateNoteRequest struct {
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}

type UpdateNoteRequest struct {
	Title    *string  `json:"title,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	IsPinned *bool    `json:"is_pinned,omitempty"`
}

type AddBlockRequest struct {
	Type      string            `json:"type"`
	ContentMD string            `json:"content_md,omitempty"`
	Items     []models.TodoItem `json:"items,omitempty"`
}

type UpdateBlockRequest struct {
	ContentMD *string           `json:"content_md,omitempty"`
	Items     []models.TodoItem `json:"items,omitempty"`
}

type ReorderBlocksRequest struct {
	Order []string `json:"order"`
}

func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)

	var req CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" {
		WriteError(w, http.StatusBadRequest, "Title is required")
		return
	}

	note, err := h.service.CreateNote(r.Context(), claims.UserID.String(), req.Title, req.Tags)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to create note")
		return
	}

	WriteJSON(w, http.StatusCreated, note)
}

func (h *NoteHandler) GetNotes(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)

	notes, err := h.service.GetUserNotes(r.Context(), claims.UserID.String())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to fetch notes")
		return
	}

	// Ensure we return empty array instead of null
	if notes == nil {
		notes = []models.Note{}
	}

	WriteJSON(w, http.StatusOK, notes)
}

func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)
	noteID := chi.URLParam(r, "id")

	note, err := h.service.GetNote(r.Context(), noteID, claims.UserID.String())
	if err != nil {
		if service.IsNotFound(err) {
			WriteError(w, http.StatusNotFound, "Note not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to fetch note")
		return
	}

	WriteJSON(w, http.StatusOK, note)
}

func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)
	noteID := chi.URLParam(r, "id")

	var req UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Tags != nil {
		updates["tags"] = req.Tags
	}
	if req.IsPinned != nil {
		updates["is_pinned"] = *req.IsPinned
	}

	if len(updates) == 0 {
		WriteError(w, http.StatusBadRequest, "No fields to update")
		return
	}

	if err := h.service.UpdateNote(r.Context(), noteID, claims.UserID.String(), updates); err != nil {
		if service.IsNotFound(err) {
			WriteError(w, http.StatusNotFound, "Note not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to update note")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "Note updated"})
}

func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)
	noteID := chi.URLParam(r, "id")

	if err := h.service.DeleteNote(r.Context(), noteID, claims.UserID.String()); err != nil {
		if service.IsNotFound(err) {
			WriteError(w, http.StatusNotFound, "Note not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to delete note")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "Note deleted"})
}

func (h *NoteHandler) AddBlock(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)
	noteID := chi.URLParam(r, "id")

	var req AddBlockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Type == "" {
		WriteError(w, http.StatusBadRequest, "Block type is required")
		return
	}

	block, err := h.service.AddBlock(r.Context(), noteID, claims.UserID.String(), req.Type, req.ContentMD, req.Items)
	if err != nil {
		if service.IsNotFound(err) {
			WriteError(w, http.StatusNotFound, "Note not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to add block")
		return
	}

	WriteJSON(w, http.StatusCreated, block)
}

func (h *NoteHandler) UpdateBlock(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)
	noteID := chi.URLParam(r, "id")
	blockID := chi.URLParam(r, "blockId")

	var req UpdateBlockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updates := make(map[string]interface{})
	if req.ContentMD != nil {
		updates["content_md"] = *req.ContentMD
	}
	if req.Items != nil {
		updates["items"] = req.Items
	}

	if len(updates) == 0 {
		WriteError(w, http.StatusBadRequest, "No fields to update")
		return
	}

	if err := h.service.UpdateBlock(r.Context(), noteID, claims.UserID.String(), blockID, updates); err != nil {
		if service.IsNotFound(err) {
			WriteError(w, http.StatusNotFound, "Note or block not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to update block")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "Block updated"})
}

func (h *NoteHandler) DeleteBlock(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)
	noteID := chi.URLParam(r, "id")
	blockID := chi.URLParam(r, "blockId")

	if err := h.service.DeleteBlock(r.Context(), noteID, claims.UserID.String(), blockID); err != nil {
		if service.IsNotFound(err) {
			WriteError(w, http.StatusNotFound, "Note or block not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to delete block")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "Block deleted"})
}

func (h *NoteHandler) ReorderBlocks(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)
	noteID := chi.URLParam(r, "id")

	var req ReorderBlocksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(req.Order) == 0 {
		WriteError(w, http.StatusBadRequest, "Order is required")
		return
	}

	if err := h.service.ReorderBlocks(r.Context(), noteID, claims.UserID.String(), req.Order); err != nil {
		if service.IsNotFound(err) {
			WriteError(w, http.StatusNotFound, "Note not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to reorder blocks")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "Blocks reordered"})
}
