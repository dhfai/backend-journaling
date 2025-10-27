package service

import (
	"context"
	"errors"

	"backend-journaling/internal/models"
	"backend-journaling/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(ctx context.Context, userID, title, descriptionMD, status, priority string, tags []string) (*models.Task, error) {
	task := &models.Task{
		UserID:   userID,
		Title:    title,
		Status:   status,
		Priority: priority,
		Tags:     tags,
	}

	if descriptionMD != "" {
		task.DescriptionMD = &descriptionMD
	}

	if err := s.repo.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetTask(ctx context.Context, taskID, userID string) (*models.Task, error) {
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return nil, errors.New("invalid task id")
	}

	return s.repo.FindByID(ctx, objID, userID)
}

func (s *TaskService) GetUserTasks(ctx context.Context, userID string) ([]models.Task, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *TaskService) UpdateTask(ctx context.Context, taskID, userID string, updates map[string]interface{}) error {
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return errors.New("invalid task id")
	}

	return s.repo.Update(ctx, objID, userID, bson.M(updates))
}

func (s *TaskService) DeleteTask(ctx context.Context, taskID, userID string) error {
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return errors.New("invalid task id")
	}

	return s.repo.Delete(ctx, objID, userID)
}
