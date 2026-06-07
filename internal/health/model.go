package health

import "time"

// HealthLog is the document stored in MongoDB for each healthcheck request.
// bson tags tell the MongoDB driver how to name each field in the database.
type HealthLog struct {
	IP        string    `bson:"ip"`        // IP address of the caller
	Timestamp time.Time `bson:"timestamp"` // when the request happened
	Message   string    `bson:"message"`   // always "health"
}
