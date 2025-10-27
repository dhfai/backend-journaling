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

type TodoRepository struct {
	collection *mongo.Collection
}

func NewTodoRepository(db *mongo.Database) *TodoRepository {
	return &TodoRepository{
		collection: db.Collection("todos"),
	}
}

func (r *TodoRepository) Create(ctx context.Context, todo *models.Todo) error {
	todo.ID = primitive.NewObjectID()
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, todo)
	return err
}

func (r *TodoRepository) FindByID(ctx context.Context, id primitive.ObjectID, userID string) (*models.Todo, error) {
	var todo models.Todo
	filter := bson.M{"_id": id, "user_id": userID}

	err := r.collection.FindOne(ctx, filter).Decode(&todo)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (r *TodoRepository) FindByUserID(ctx context.Context, userID string) ([]models.Todo, error) {
	filter := bson.M{"user_id": userID}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var todos []models.Todo
	if err := cursor.All(ctx, &todos); err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *TodoRepository) Update(ctx context.Context, id primitive.ObjectID, userID string, update bson.M) error {
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

func (r *TodoRepository) Delete(ctx context.Context, id primitive.ObjectID, userID string) error {
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
