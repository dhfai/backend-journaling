package models

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FlexibleTime handles multiple date/datetime formats
type FlexibleTime struct {
	time.Time
}

// UnmarshalJSON handles JSON decoding from HTTP requests
func (ft *FlexibleTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" || s == `""` {
		return nil
	}

	// Remove quotes
	s = s[1 : len(s)-1]

	// Try different formats
	formats := []string{
		time.RFC3339,               // "2006-01-02T15:04:05Z07:00"
		"2006-01-02T15:04:05.999Z", // ISO 8601 with milliseconds
		"2006-01-02T15:04:05",      // ISO 8601 without timezone
		"2006-01-02",               // Date only
		"2006-01-02 15:04:05",      // DateTime with space
	}

	var err error
	for _, format := range formats {
		ft.Time, err = time.Parse(format, s)
		if err == nil {
			return nil
		}
	}

	return err
}

// MarshalJSON handles JSON encoding for HTTP responses
func (ft FlexibleTime) MarshalJSON() ([]byte, error) {
	if ft.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(ft.Time.Format(time.RFC3339))
}

// UnmarshalBSONValue handles BSON decoding from MongoDB
func (ft *FlexibleTime) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	// Handle null/undefined
	if t == bsontype.Null || t == bsontype.Undefined {
		ft.Time = time.Time{}
		return nil
	}

	// Handle DateTime type (native BSON datetime)
	if t == bsontype.DateTime {
		dt := primitive.DateTime(0)
		if err := bson.Unmarshal(data, &dt); err != nil {
			return err
		}
		ft.Time = dt.Time()
		return nil
	}

	// Handle String type (stored as string in MongoDB)
	if t == bsontype.String {
		var s string
		if err := bson.Unmarshal(data, &s); err != nil {
			return err
		}

		// Try different formats
		formats := []string{
			time.RFC3339,               // "2006-01-02T15:04:05Z07:00"
			"2006-01-02T15:04:05.999Z", // ISO 8601 with milliseconds
			"2006-01-02T15:04:05",      // ISO 8601 without timezone
			"2006-01-02",               // Date only
			"2006-01-02 15:04:05",      // DateTime with space
		}

		var err error
		for _, format := range formats {
			ft.Time, err = time.Parse(format, s)
			if err == nil {
				return nil
			}
		}

		return err
	}

	return nil
}

// MarshalBSONValue handles BSON encoding to MongoDB
func (ft FlexibleTime) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if ft.IsZero() {
		return bsontype.Null, nil, nil
	}
	return bson.MarshalValue(primitive.NewDateTimeFromTime(ft.Time))
}

type Note struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID    string              `bson:"user_id" json:"user_id"`
	GroupID   *primitive.ObjectID `bson:"group_id,omitempty" json:"group_id,omitempty"`
	Title     string              `bson:"title" json:"title"`
	Blocks    []Block             `bson:"blocks" json:"blocks"`
	Tags      []string            `bson:"tags" json:"tags"`
	IsPinned  bool                `bson:"is_pinned" json:"is_pinned"`
	CreatedAt time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time           `bson:"updated_at" json:"updated_at"`
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
	DueDate   time.Time          `bson:"due_date,omitempty" json:"due_date,omitempty"`
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
	Deadline      time.Time          `bson:"deadline" json:"deadline,omitempty"`
	Tags          []string           `bson:"tags" json:"tags"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

type NoteGroup struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"user_id" json:"user_id"`
	Name        string             `bson:"name" json:"name"`
	Description *string            `bson:"description,omitempty" json:"description,omitempty"`
	Color       *string            `bson:"color,omitempty" json:"color,omitempty"`
	Icon        *string            `bson:"icon,omitempty" json:"icon,omitempty"`
	IsPinned    bool               `bson:"is_pinned" json:"is_pinned"`
	IsArchived  bool               `bson:"is_archived" json:"is_archived"`
	NotesCount  int                `bson:"notes_count" json:"notes_count"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
