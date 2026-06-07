package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// DB wraps the MongoDB client and the target database name.
// Use New() to create it — never build this struct directly.
type DB struct {
	Client *mongo.Client
	Name   string
}

// New connects to MongoDB using the given URI and database name.
// It pings the server to confirm the connection is live before returning.
// Returns an error if the connection or ping fails.
func New(uri string, dbName string) (*DB, error) {
	// WithTimeout gives the connection attempt a 10-second deadline.
	// If MongoDB is unreachable, we fail fast instead of hanging forever.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // always release the context resources

	// ApplyURI parses the full connection string (host, auth, options).
	// mongo.Connect in v2 does NOT take a context — only options.
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Ping verifies the server is actually reachable.
	// Connect alone does not guarantee the server responded.
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &DB{Client: client, Name: dbName}, nil
}

// Collection returns a handle to a MongoDB collection by name.
// Use this in your repositories: db.Collection("healthchecks")
func (db *DB) Collection(name string) *mongo.Collection {
	return db.Client.Database(db.Name).Collection(name)
}
