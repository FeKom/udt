package health

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Repository handles all database operations for healthcheck logs.
// It holds a reference to the MongoDB collection it writes to.
type Repository struct {
	collection *mongo.Collection
}

// NewRepository creates a Repository that writes to the given collection.
// Call this in main.go: health.NewRepository(db.Collection("healthchecks"))
func NewRepository(col *mongo.Collection) *Repository {
	return &Repository{collection: col}
}

// Save inserts a healthcheck log document into MongoDB.
// It builds the document from the caller IP and writes it with a 5-second deadline.
func (r *Repository) Save(ctx context.Context, ip string) error {
	doc := HealthLog{
		IP:        ip,
		Timestamp: time.Now(), // UTC timestamp of when the request arrived
		Message:   "health",
	}

	// InsertOne writes a single document to the collection.
	// The second return value is the inserted document ID — we don't need it here.
	result, err := r.collection.InsertOne(ctx, doc)
	log.Printf("InsertOne result=%v err=%v\n", result, err)
	return err
}
