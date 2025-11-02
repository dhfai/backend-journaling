package service

import (
	"context"
	"errors"

	"backend-journaling/internal/models"
	"backend-journaling/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrGroupNotFound         = errors.New("note group not found")
	ErrNoteNotFound          = errors.New("note not found")
	ErrGroupNameRequired     = errors.New("group name is required")
	ErrInvalidGroupID        = errors.New("invalid group id")
	ErrInvalidNoteID         = errors.New("invalid note id")
	ErrNoteNotInGroup        = errors.New("note is not in any group")
	ErrCannotMoveToSameGroup = errors.New("note is already in this group")
)

type NoteGroupService struct {
	repo *repository.NoteGroupRepository
}

func NewNoteGroupService(repo *repository.NoteGroupRepository) *NoteGroupService {
	return &NoteGroupService{repo: repo}
}

// CreateGroup creates a new note group
func (s *NoteGroupService) CreateGroup(ctx context.Context, userID, name string, description, color, icon *string) (*models.NoteGroup, error) {
	if name == "" {
		return nil, ErrGroupNameRequired
	}

	group := &models.NoteGroup{
		UserID:      userID,
		Name:        name,
		Description: description,
		Color:       color,
		Icon:        icon,
		IsPinned:    false,
		IsArchived:  false,
	}

	if err := s.repo.Create(ctx, group); err != nil {
		return nil, err
	}

	return group, nil
}

// GetGroup gets a note group by ID
func (s *NoteGroupService) GetGroup(ctx context.Context, groupID, userID string) (*models.NoteGroup, error) {
	id, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return nil, ErrInvalidGroupID
	}

	group, err := s.repo.FindByID(ctx, id, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrGroupNotFound
		}
		return nil, err
	}

	return group, nil
}

// GetGroups gets all note groups for a user
func (s *NoteGroupService) GetGroups(ctx context.Context, userID string, isPinned, isArchived *bool) ([]*models.NoteGroup, error) {
	groups, err := s.repo.FindAll(ctx, userID, isPinned, isArchived)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

// UpdateGroup updates a note group
func (s *NoteGroupService) UpdateGroup(ctx context.Context, groupID, userID string, name *string, description, color, icon *string) (*models.NoteGroup, error) {
	id, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return nil, ErrInvalidGroupID
	}

	updates := bson.M{}

	if name != nil {
		if *name == "" {
			return nil, ErrGroupNameRequired
		}
		updates["name"] = *name
	}

	if description != nil {
		updates["description"] = *description
	}

	if color != nil {
		updates["color"] = *color
	}

	if icon != nil {
		updates["icon"] = *icon
	}

	if len(updates) == 0 {
		return s.GetGroup(ctx, groupID, userID)
	}

	if err := s.repo.Update(ctx, id, userID, updates); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrGroupNotFound
		}
		return nil, err
	}

	return s.GetGroup(ctx, groupID, userID)
}

// PinGroup pins/unpins a note group
func (s *NoteGroupService) PinGroup(ctx context.Context, groupID, userID string, isPinned bool) error {
	id, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return ErrInvalidGroupID
	}

	updates := bson.M{"is_pinned": isPinned}

	if err := s.repo.Update(ctx, id, userID, updates); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrGroupNotFound
		}
		return err
	}

	return nil
}

// ArchiveGroup archives/unarchives a note group
func (s *NoteGroupService) ArchiveGroup(ctx context.Context, groupID, userID string, isArchived bool) error {
	id, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return ErrInvalidGroupID
	}

	updates := bson.M{"is_archived": isArchived}

	if err := s.repo.Update(ctx, id, userID, updates); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrGroupNotFound
		}
		return err
	}

	return nil
}

// DeleteGroup deletes a note group and removes group_id from all notes in the group
func (s *NoteGroupService) DeleteGroup(ctx context.Context, groupID, userID string) error {
	id, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return ErrInvalidGroupID
	}

	// Check if group exists
	_, err = s.repo.FindByID(ctx, id, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrGroupNotFound
		}
		return err
	}

	// Remove group_id from all notes in this group
	if err := s.repo.RemoveGroupFromAllNotes(ctx, id, userID); err != nil {
		return err
	}

	// Delete the group
	if err := s.repo.Delete(ctx, id, userID); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrGroupNotFound
		}
		return err
	}

	return nil
}

// AddNoteToGroup adds a note to a group
func (s *NoteGroupService) AddNoteToGroup(ctx context.Context, noteID, groupID, userID string) error {
	nid, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		return ErrInvalidNoteID
	}

	gid, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return ErrInvalidGroupID
	}

	// Verify group exists
	_, err = s.repo.FindByID(ctx, gid, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrGroupNotFound
		}
		return err
	}

	// Add note to group
	if err := s.repo.AddNoteToGroup(ctx, nid, gid, userID); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNoteNotFound
		}
		return err
	}

	return nil
}

// RemoveNoteFromGroup removes a note from its group
func (s *NoteGroupService) RemoveNoteFromGroup(ctx context.Context, noteID, userID string) error {
	nid, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		return ErrInvalidNoteID
	}

	// Remove note from group
	if err := s.repo.RemoveNoteFromGroup(ctx, nid, userID); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNoteNotFound
		}
		return err
	}

	return nil
}

// MoveNotesToGroup moves notes from one group to another
func (s *NoteGroupService) MoveNotesToGroup(ctx context.Context, noteIDs []string, newGroupID, userID string) error {
	// Convert note IDs
	nids := make([]primitive.ObjectID, len(noteIDs))
	for i, id := range noteIDs {
		nid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return ErrInvalidNoteID
		}
		nids[i] = nid
	}

	// Convert group ID
	gid, err := primitive.ObjectIDFromHex(newGroupID)
	if err != nil {
		return ErrInvalidGroupID
	}

	// Verify group exists
	_, err = s.repo.FindByID(ctx, gid, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrGroupNotFound
		}
		return err
	}

	// Move notes to new group
	if err := s.repo.MoveNotesToGroup(ctx, nids, gid, userID); err != nil {
		return err
	}

	return nil
}

// GetNotesInGroup gets all notes in a group
func (s *NoteGroupService) GetNotesInGroup(ctx context.Context, groupID, userID string) ([]*models.Note, error) {
	id, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return nil, ErrInvalidGroupID
	}

	// Verify group exists
	_, err = s.repo.FindByID(ctx, id, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrGroupNotFound
		}
		return nil, err
	}

	notes, err := s.repo.GetNotesInGroup(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	return notes, nil
}
