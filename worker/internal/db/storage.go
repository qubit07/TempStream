package db

import (
	"context"
	"tempstream-worker/internal/events"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type ReadingStore struct {
    collection *mongo.Collection
}

type readingDocument struct {
	EventID    string    `bson:"eventId"`
	SensorID   string    `bson:"sensorId"`
	SensorType string    `bson:"sensorType"`
	Value      float64   `bson:"value"`
	Unit       string    `bson:"unit"`
	Timestamp  time.Time `bson:"timestamp"`
	ReceivedAt time.Time `bson:"receivedAt"`
	StoredAt   time.Time `bson:"storedAt"`
}

func NewReadingStore(collection *mongo.Collection) *ReadingStore {
    return &ReadingStore{
        collection: collection,
    }
}

func (s *ReadingStore) Save(
    ctx context.Context,
    event events.ReadingEvent,
) error {
	doc := readingDocument{
		EventID:    event.EventID,
		SensorID:   event.SensorID,
		SensorType: event.SensorType,
		Value:      event.Value,
		Unit:       event.Unit,
		Timestamp:  event.Timestamp,
		ReceivedAt: event.ReceivedAt,
		StoredAt:   time.Now().UTC(),
	}

    _, err := s.collection.InsertOne(ctx, doc)
    return err
}