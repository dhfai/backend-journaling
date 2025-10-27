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

type NoteRepository struct {
	collection *mongo.Collection
}

func NewNoteRepository(db *mongo.Database) *NoteRepository {
	return &NoteRepository{
		collection: db.Collection("notes"),
	}
}

func (r *NoteRepository) Create(ctx context.Context, note *models.Note) error {
	note.ID = primitive.NewObjectID()
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()

	if note.Blocks == nil {
		note.Blocks = []models.Block{}
	}
	if note.Tags == nil {
		note.Tags = []string{}
	}

	_, err := r.collection.InsertOne(ctx, note)
	return err
}

func (r *NoteRepository) FindByID(ctx context.Context, id primitive.ObjectID, userID string) (*models.Note, error) {
	var note models.Note
	filter := bson.M{"_id": id, "user_id": userID}

	err := r.collection.FindOne(ctx, filter).Decode(&note)
	if err != nil {
		return nil, err
	}

	return &note, nil
}

func (r *NoteRepository) FindByUserID(ctx context.Context, userID string) ([]models.Note, error) {
	filter := bson.M{"user_id": userID}
	opts := options.Find().SetSort(bson.D{{Key: "updated_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notes []models.Note
	if err := cursor.All(ctx, &notes); err != nil {
		return nil, err
	}

	return notes, nil
}

func (r *NoteRepository) Update(ctx context.Context, id primitive.ObjectID, userID string, update bson.M) error {
	filter := bson.M{"_id": id, "user_id": userID}
	update["updated_at"] = time.Now()

	updateDoc := bson.M{"$set": update}
	result, err := r.collection.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *NoteRepository) Delete(ctx context.Context, id primitive.ObjectID, userID string) error {
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

func (r *NoteRepository) AddBlock(ctx context.Context, id primitive.ObjectID, userID string, block models.Block) error {
	filter := bson.M{"_id": id, "user_id": userID}
	update := bson.M{
		"$push": bson.M{"blocks": block},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *NoteRepository) UpdateBlock(ctx context.Context, noteID primitive.ObjectID, userID, blockID string, updateData bson.M) error {
	filter := bson.M{
		"_id":       noteID,
		"user_id":   userID,
		"blocks.id": blockID,
	}

	update := bson.M{"$set": bson.M{"updated_at": time.Now()}}
	for key, value := range updateData {
		update["$set"].(bson.M)["blocks.$."+key] = value
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *NoteRepository) DeleteBlock(ctx context.Context, noteID primitive.ObjectID, userID, blockID string) error {
	filter := bson.M{"_id": noteID, "user_id": userID}
	update := bson.M{
		"$pull": bson.M{"blocks": bson.M{"id": blockID}},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *NoteRepository) ReorderBlocks(ctx context.Context, noteID primitive.ObjectID, userID string, blocks []models.Block) error {
	filter := bson.M{"_id": noteID, "user_id": userID}
	update := bson.M{
		"$set": bson.M{
			"blocks":     blocks,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
