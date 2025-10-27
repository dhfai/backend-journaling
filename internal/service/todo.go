package service

import (
	"context"
	"errors"

	"backend-journaling/internal/models"
	"backend-journaling/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoService struct {
	repo *repository.TodoRepository
}

func NewTodoService(repo *repository.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) CreateTodo(ctx context.Context, userID, title, priority string, dueDate interface{}) (*models.Todo, error) {
	todo := &models.Todo{
		UserID:   userID,
		Title:    title,
		Done:     false,
		Priority: priority,
	}

	if err := s.repo.Create(ctx, todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *TodoService) GetTodo(ctx context.Context, todoID, userID string) (*models.Todo, error) {
	objID, err := primitive.ObjectIDFromHex(todoID)
	if err != nil {
		return nil, errors.New("invalid todo id")
	}

	return s.repo.FindByID(ctx, objID, userID)
}

func (s *TodoService) GetUserTodos(ctx context.Context, userID string) ([]models.Todo, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *TodoService) UpdateTodo(ctx context.Context, todoID, userID string, updates map[string]interface{}) error {
	objID, err := primitive.ObjectIDFromHex(todoID)
	if err != nil {
		return errors.New("invalid todo id")
	}

	return s.repo.Update(ctx, objID, userID, bson.M(updates))
}

func (s *TodoService) DeleteTodo(ctx context.Context, todoID, userID string) error {
	objID, err := primitive.ObjectIDFromHex(todoID)
	if err != nil {
		return errors.New("invalid todo id")
	}

	return s.repo.Delete(ctx, objID, userID)
}
