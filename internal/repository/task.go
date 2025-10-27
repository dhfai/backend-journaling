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

type TaskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(db *mongo.Database) *TaskRepository {
	return &TaskRepository{
		collection: db.Collection("tasks"),
	}
}

func (r *TaskRepository) Create(ctx context.Context, task *models.Task) error {
	task.ID = primitive.NewObjectID()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	if task.Tags == nil {
		task.Tags = []string{}
	}

	_, err := r.collection.InsertOne(ctx, task)
	return err
}

func (r *TaskRepository) FindByID(ctx context.Context, id primitive.ObjectID, userID string) (*models.Task, error) {
	var task models.Task
	filter := bson.M{"_id": id, "user_id": userID}

	err := r.collection.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) FindByUserID(ctx context.Context, userID string) ([]models.Task, error) {
	filter := bson.M{"user_id": userID}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) Update(ctx context.Context, id primitive.ObjectID, userID string, update bson.M) error {
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

func (r *TaskRepository) Delete(ctx context.Context, id primitive.ObjectID, userID string) error {
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
