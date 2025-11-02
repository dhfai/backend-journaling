package repository

import (
	"context"
	"time"

	"backend-journaling/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NoteGroupRepository struct {
	collection     *mongo.Collection
	noteCollection *mongo.Collection
}

func NewNoteGroupRepository(db *mongo.Database) *NoteGroupRepository {
	return &NoteGroupRepository{
		collection:     db.Collection("note_groups"),
		noteCollection: db.Collection("notes"),
	}
}

// Create creates a new note group
func (r *NoteGroupRepository) Create(ctx context.Context, group *models.NoteGroup) error {
	group.ID = primitive.NewObjectID()
	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()
	group.NotesCount = 0

	_, err := r.collection.InsertOne(ctx, group)
	return err
}

// FindByID finds a note group by ID for a specific user
func (r *NoteGroupRepository) FindByID(ctx context.Context, id primitive.ObjectID, userID string) (*models.NoteGroup, error) {
	var group models.NoteGroup
	filter := bson.M{"_id": id, "user_id": userID}

	err := r.collection.FindOne(ctx, filter).Decode(&group)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

// FindAll finds all note groups for a user with optional filters
func (r *NoteGroupRepository) FindAll(ctx context.Context, userID string, isPinned *bool, isArchived *bool) ([]*models.NoteGroup, error) {
	filter := bson.M{"user_id": userID}

	if isPinned != nil {
		filter["is_pinned"] = *isPinned
	}

	if isArchived != nil {
		filter["is_archived"] = *isArchived
	}

	opts := options.Find().SetSort(bson.D{
		{Key: "is_pinned", Value: -1},
		{Key: "updated_at", Value: -1},
	})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var groups []*models.NoteGroup
	if err := cursor.All(ctx, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

// Update updates a note group
func (r *NoteGroupRepository) Update(ctx context.Context, id primitive.ObjectID, userID string, updates bson.M) error {
	updates["updated_at"] = time.Now()

	filter := bson.M{"_id": id, "user_id": userID}
	update := bson.M{"$set": updates}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// Delete deletes a note group
func (r *NoteGroupRepository) Delete(ctx context.Context, id primitive.ObjectID, userID string) error {
	filter := bson.M{"_id": id, "user_id": userID}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// AddNoteToGroup adds a note to a group
func (r *NoteGroupRepository) AddNoteToGroup(ctx context.Context, noteID primitive.ObjectID, groupID primitive.ObjectID, userID string) error {
	// Update note with group_id
	noteFilter := bson.M{"_id": noteID, "user_id": userID}
	noteUpdate := bson.M{
		"$set": bson.M{
			"group_id":   groupID,
			"updated_at": time.Now(),
		},
	}

	result, err := r.noteCollection.UpdateOne(ctx, noteFilter, noteUpdate)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	// Increment notes_count in group
	groupFilter := bson.M{"_id": groupID, "user_id": userID}
	groupUpdate := bson.M{
		"$inc": bson.M{"notes_count": 1},
		"$set": bson.M{"updated_at": time.Now()},
	}

	_, err = r.collection.UpdateOne(ctx, groupFilter, groupUpdate)
	return err
}

// RemoveNoteFromGroup removes a note from a group
func (r *NoteGroupRepository) RemoveNoteFromGroup(ctx context.Context, noteID primitive.ObjectID, userID string) error {
	// Get current note to check if it has a group
	var note models.Note
	noteFilter := bson.M{"_id": noteID, "user_id": userID}
	err := r.noteCollection.FindOne(ctx, noteFilter).Decode(&note)
	if err != nil {
		return err
	}

	if note.GroupID == nil {
		return nil // Note is not in any group
	}

	// Remove group_id from note
	noteUpdate := bson.M{
		"$unset": bson.M{"group_id": ""},
		"$set":   bson.M{"updated_at": time.Now()},
	}

	result, err := r.noteCollection.UpdateOne(ctx, noteFilter, noteUpdate)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	// Decrement notes_count in group
	groupFilter := bson.M{"_id": note.GroupID, "user_id": userID}
	groupUpdate := bson.M{
		"$inc": bson.M{"notes_count": -1},
		"$set": bson.M{"updated_at": time.Now()},
	}

	_, err = r.collection.UpdateOne(ctx, groupFilter, groupUpdate)
	return err
}

// MoveNotesToGroup moves multiple notes to a different group
func (r *NoteGroupRepository) MoveNotesToGroup(ctx context.Context, noteIDs []primitive.ObjectID, newGroupID primitive.ObjectID, userID string) error {
	// Get all notes to calculate old groups
	notesFilter := bson.M{
		"_id":     bson.M{"$in": noteIDs},
		"user_id": userID,
	}

	cursor, err := r.noteCollection.Find(ctx, notesFilter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var notes []models.Note
	if err := cursor.All(ctx, &notes); err != nil {
		return err
	}

	// Count notes per old group
	oldGroupCounts := make(map[primitive.ObjectID]int)
	for _, note := range notes {
		if note.GroupID != nil {
			oldGroupCounts[*note.GroupID]++
		}
	}

	// Update all notes with new group_id
	notesUpdate := bson.M{
		"$set": bson.M{
			"group_id":   newGroupID,
			"updated_at": time.Now(),
		},
	}

	_, err = r.noteCollection.UpdateMany(ctx, notesFilter, notesUpdate)
	if err != nil {
		return err
	}

	// Decrement count in old groups
	for oldGroupID, count := range oldGroupCounts {
		groupFilter := bson.M{"_id": oldGroupID, "user_id": userID}
		groupUpdate := bson.M{
			"$inc": bson.M{"notes_count": -count},
			"$set": bson.M{"updated_at": time.Now()},
		}
		r.collection.UpdateOne(ctx, groupFilter, groupUpdate)
	}

	// Increment count in new group
	newGroupFilter := bson.M{"_id": newGroupID, "user_id": userID}
	newGroupUpdate := bson.M{
		"$inc": bson.M{"notes_count": len(notes)},
		"$set": bson.M{"updated_at": time.Now()},
	}

	_, err = r.collection.UpdateOne(ctx, newGroupFilter, newGroupUpdate)
	return err
}

// CountNotesInGroup counts notes in a group
func (r *NoteGroupRepository) CountNotesInGroup(ctx context.Context, groupID primitive.ObjectID, userID string) (int64, error) {
	filter := bson.M{
		"group_id": groupID,
		"user_id":  userID,
	}

	count, err := r.noteCollection.CountDocuments(ctx, filter)
	return count, err
}

// GetNotesInGroup gets all notes in a group
func (r *NoteGroupRepository) GetNotesInGroup(ctx context.Context, groupID primitive.ObjectID, userID string) ([]*models.Note, error) {
	filter := bson.M{
		"group_id": groupID,
		"user_id":  userID,
	}

	opts := options.Find().SetSort(bson.D{{Key: "updated_at", Value: -1}})

	cursor, err := r.noteCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notes []*models.Note
	if err := cursor.All(ctx, &notes); err != nil {
		return nil, err
	}

	return notes, nil
}

// RecalculateNotesCount recalculates the notes_count for a group
func (r *NoteGroupRepository) RecalculateNotesCount(ctx context.Context, groupID primitive.ObjectID, userID string) error {
	count, err := r.CountNotesInGroup(ctx, groupID, userID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": groupID, "user_id": userID}
	update := bson.M{
		"$set": bson.M{
			"notes_count": count,
			"updated_at":  time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

// RemoveGroupFromAllNotes removes group_id from all notes in a specific group
func (r *NoteGroupRepository) RemoveGroupFromAllNotes(ctx context.Context, groupID primitive.ObjectID, userID string) error {
	filter := bson.M{
		"group_id": groupID,
		"user_id":  userID,
	}

	update := bson.M{
		"$unset": bson.M{"group_id": ""},
		"$set":   bson.M{"updated_at": time.Now()},
	}

	_, err := r.noteCollection.UpdateMany(ctx, filter, update)
	return err
}
