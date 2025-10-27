package handlers

import (
	"encoding/json"
	"net/http"

	"backend-journaling/internal/service"
	"backend-journaling/pkg/jwt"

	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

type CreateTaskRequest struct {
	Title         string      `json:"title"`
	DescriptionMD string      `json:"description_md,omitempty"`
	Priority      string      `json:"priority"`
	Deadline      interface{} `json:"deadline,omitempty"`
	Tags          []string    `json:"tags"`
}

type UpdateTaskRequest struct {
	Title         *string     `json:"title,omitempty"`
	DescriptionMD *string     `json:"description_md,omitempty"`
	Status        *string     `json:"status,omitempty"`
	Priority      *string     `json:"priority,omitempty"`
	Deadline      interface{} `json:"deadline,omitempty"`
	Tags          []string    `json:"tags,omitempty"`
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)

	var req CreateTaskRequest
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

	status := "todo"

	task, err := h.service.CreateTask(r.Context(), claims.UserID.String(), req.Title, req.DescriptionMD, status, req.Priority, req.Tags)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to create task")
		return
	}

	WriteJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)

	tasks, err := h.service.GetUserTasks(r.Context(), claims.UserID.String())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to fetch tasks")
		return
	}

	WriteJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)
	taskID := chi.URLParam(r, "id")

	task, err := h.service.GetTask(r.Context(), taskID, claims.UserID.String())
	if err != nil {
		if service.IsNotFound(err) {
			WriteError(w, http.StatusNotFound, "Task not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to fetch task")
		return
	}

	WriteJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)
	taskID := chi.URLParam(r, "id")

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.DescriptionMD != nil {
		updates["description_md"] = *req.DescriptionMD
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.Deadline != nil {
		updates["deadline"] = req.Deadline
	}
	if req.Tags != nil {
		updates["tags"] = req.Tags
	}

	if len(updates) == 0 {
		WriteError(w, http.StatusBadRequest, "No fields to update")
		return
	}

	if err := h.service.UpdateTask(r.Context(), taskID, claims.UserID.String(), updates); err != nil {
		if service.IsNotFound(err) {
			WriteError(w, http.StatusNotFound, "Task not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to update task")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "Task updated"})
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Claims)
	taskID := chi.URLParam(r, "id")

	if err := h.service.DeleteTask(r.Context(), taskID, claims.UserID.String()); err != nil {
		if service.IsNotFound(err) {
			WriteError(w, http.StatusNotFound, "Task not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to delete task")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "Task deleted"})
}
