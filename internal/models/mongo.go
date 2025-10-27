package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	Title     string             `bson:"title" json:"title"`
	Blocks    []Block            `bson:"blocks" json:"blocks"`
	Tags      []string           `bson:"tags" json:"tags"`
	IsPinned  bool               `bson:"is_pinned" json:"is_pinned"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type Block struct {
	ID        string     `bson:"id" json:"id"`
	Type      string     `bson:"type" json:"type"`
	Order     int        `bson:"order" json:"order"`
	ContentMD *string    `bson:"content_md,omitempty" json:"content_md,omitempty"`
	Items     []TodoItem `bson:"items,omitempty" json:"items,omitempty"`
}

type TodoItem struct {
	ID   string `bson:"id" json:"id"`
	Text string `bson:"text" json:"text"`
	Done bool   `bson:"done" json:"done"`
}

type Todo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	Title     string             `bson:"title" json:"title"`
	Done      bool               `bson:"done" json:"done"`
	Priority  string             `bson:"priority" json:"priority"`
	DueDate   *time.Time         `bson:"due_date,omitempty" json:"due_date,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type Task struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        string             `bson:"user_id" json:"user_id"`
	Title         string             `bson:"title" json:"title"`
	DescriptionMD *string            `bson:"description_md,omitempty" json:"description_md,omitempty"`
	Status        string             `bson:"status" json:"status"`
	Priority      string             `bson:"priority" json:"priority"`
	Deadline      *time.Time         `bson:"deadline,omitempty" json:"deadline,omitempty"`
	Tags          []string           `bson:"tags" json:"tags"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}
