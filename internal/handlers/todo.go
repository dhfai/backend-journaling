package handlers

import (
	"encoding/json"
	"net/http"

	"backend-journaling/internal/models"
	"backend-journaling/internal/service"
	"backend-journaling/pkg/jwt"

	"github.com/go-chi/chi/v5"
)

type TodoHandler struct {
	service *service.TodoService
}

func NewTodoHandler(service *service.TodoService) *TodoHandler {
	return &TodoHandler{service: service}
}

type CreateTodoRequest struct {
	Title    string      `json:"title"`
	Priority string      `json:"priority"`
	DueDate  interface{} `json:"due_date,omitempty"`
}

type UpdateTodoRequest struct {
	Title    *string     `json:"title,omitempty"`
	Done     *bool       `json:"done,omitempty"`
	Priority *string     `json:"priority,omitempty"`
	DueDate  interface{} `json:"due_date,omitempty"`
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)

	var req CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" {
		WriteError(w, http.StatusBadRequest, "Title is required")
		return
	}

	if req.Priority == "" {
		req.Priority = "medium"
	}

	todo, err := h.service.CreateTodo(r.Context(), claims.UserID.String(), req.Title, req.Priority, req.DueDate)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to create todo")
		return
	}

	WriteJSON(w, http.StatusCreated, todo)
}

func (h *TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)

	todos, err := h.service.GetUserTodos(r.Context(), claims.UserID.String())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to fetch todos")
		return
	}

	// Ensure we return empty array instead of null
	if todos == nil {
		todos = []models.Todo{}
	}

	WriteJSON(w, http.StatusOK, todos)
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)
	todoID := chi.URLParam(r, "id")

	var req UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Done != nil {
		updates["done"] = *req.Done
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.DueDate != nil {
		updates["due_date"] = req.DueDate
	}

	if len(updates) == 0 {
		WriteError(w, http.StatusBadRequest, "No fields to update")
		return
	}

	if err := h.service.UpdateTodo(r.Context(), todoID, claims.UserID.String(), updates); err != nil {
		if service.IsNotFound(err) {
			WriteError(w, http.StatusNotFound, "Todo not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to update todo")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "Todo updated"})
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)
	todoID := chi.URLParam(r, "id")

	if err := h.service.DeleteTodo(r.Context(), todoID, claims.UserID.String()); err != nil {
		if service.IsNotFound(err) {
			WriteError(w, http.StatusNotFound, "Todo not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to delete todo")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "Todo deleted"})
}
