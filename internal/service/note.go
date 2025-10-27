package service

import (
	"context"
	"errors"

	"backend-journaling/internal/models"
	"backend-journaling/internal/repository"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NoteService struct {
	repo *repository.NoteRepository
}

func NewNoteService(repo *repository.NoteRepository) *NoteService {
	return &NoteService{repo: repo}
}

func (s *NoteService) CreateNote(ctx context.Context, userID string, title string, tags []string) (*models.Note, error) {
	note := &models.Note{
		UserID:   userID,
		Title:    title,
		Tags:     tags,
		IsPinned: false,
		Blocks:   []models.Block{},
	}

	if err := s.repo.Create(ctx, note); err != nil {
		return nil, err
	}

	return note, nil
}

func (s *NoteService) GetNote(ctx context.Context, noteID, userID string) (*models.Note, error) {
	objID, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		return nil, errors.New("invalid note id")
	}

	return s.repo.FindByID(ctx, objID, userID)
}

func (s *NoteService) GetUserNotes(ctx context.Context, userID string) ([]models.Note, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *NoteService) UpdateNote(ctx context.Context, noteID, userID string, updates map[string]interface{}) error {
	objID, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		return errors.New("invalid note id")
	}

	return s.repo.Update(ctx, objID, userID, bson.M(updates))
}

func (s *NoteService) DeleteNote(ctx context.Context, noteID, userID string) error {
	objID, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		return errors.New("invalid note id")
	}

	return s.repo.Delete(ctx, objID, userID)
}

func (s *NoteService) AddBlock(ctx context.Context, noteID, userID string, blockType, contentMD string, items []models.TodoItem) (*models.Block, error) {
	objID, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		return nil, errors.New("invalid note id")
	}

	note, err := s.repo.FindByID(ctx, objID, userID)
	if err != nil {
		return nil, err
	}

	block := models.Block{
		ID:    uuid.New().String(),
		Type:  blockType,
		Order: len(note.Blocks),
	}

	if blockType == "todo" {
		block.Items = items
	} else {
		block.ContentMD = &contentMD
	}

	if err := s.repo.AddBlock(ctx, objID, userID, block); err != nil {
		return nil, err
	}

	return &block, nil
}

func (s *NoteService) UpdateBlock(ctx context.Context, noteID, userID, blockID string, updates map[string]interface{}) error {
	objID, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		return errors.New("invalid note id")
	}

	return s.repo.UpdateBlock(ctx, objID, userID, blockID, bson.M(updates))
}

func (s *NoteService) DeleteBlock(ctx context.Context, noteID, userID, blockID string) error {
	objID, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		return errors.New("invalid note id")
	}

	return s.repo.DeleteBlock(ctx, objID, userID, blockID)
}

func (s *NoteService) ReorderBlocks(ctx context.Context, noteID, userID string, blockOrder []string) error {
	objID, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		return errors.New("invalid note id")
	}

	note, err := s.repo.FindByID(ctx, objID, userID)
	if err != nil {
		return err
	}

	blockMap := make(map[string]*models.Block)
	for i := range note.Blocks {
		blockMap[note.Blocks[i].ID] = &note.Blocks[i]
	}

	var reorderedBlocks []models.Block
	for i, blockID := range blockOrder {
		if block, exists := blockMap[blockID]; exists {
			newBlock := *block
			newBlock.Order = i
			reorderedBlocks = append(reorderedBlocks, newBlock)
		}
	}

	if len(reorderedBlocks) != len(note.Blocks) {
		return errors.New("invalid block order")
	}

	return s.repo.ReorderBlocks(ctx, objID, userID, reorderedBlocks)
}

func IsNotFound(err error) bool {
	return errors.Is(err, mongo.ErrNoDocuments)
}
