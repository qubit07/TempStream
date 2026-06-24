package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type config struct {
	rabbitMQURL      string
	queueName        string
	mongoURI         string
	mongoDatabase    string
	mongoCollection  string
	operationTimeout time.Duration
}

type readingEvent struct {
	EventID    string    `json:"eventId"`
	SensorID   string    `json:"sensorId"`
	SensorType string    `json:"sensorType"`
	Value      float64   `json:"value"`
	Unit       string    `json:"unit"`
	Timestamp  time.Time `json:"timestamp"`
	ReceivedAt time.Time `json:"receivedAt"`
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

func main() {
	cfg := loadConfig()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	rabbitConn, err := amqp.Dial(cfg.rabbitMQURL)
	if err != nil {
		log.Fatalf("connect rabbitmq: %v", err)
	}
	defer rabbitConn.Close()

	channel, err := rabbitConn.Channel()
	if err != nil {
		log.Fatalf("open rabbitmq channel: %v", err)
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(cfg.queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("declare queue: %v", err)
	}

	if err := channel.Qos(1, 0, false); err != nil {
		log.Fatalf("set qos: %v", err)
	}

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.mongoURI))
	if err != nil {
		log.Fatalf("connect mongodb: %v", err)
	}
	defer func() {
		disconnectCtx, cancel := context.WithTimeout(context.Background(), cfg.operationTimeout)
		defer cancel()
		if err := mongoClient.Disconnect(disconnectCtx); err != nil {
			log.Printf("disconnect mongodb: %v", err)
		}
	}()

	collection := mongoClient.Database(cfg.mongoDatabase).Collection(cfg.mongoCollection)

	deliveries, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("consume queue: %v", err)
	}

	log.Printf("worker consuming queue %q", queue.Name)

	for {
		select {
		case <-ctx.Done():
			log.Println("worker shutting down")
			return
		case delivery, ok := <-deliveries:
			if !ok {
				log.Println("rabbitmq deliveries channel closed")
				return
			}
			handleDelivery(cfg, collection, delivery)
		}
	}
}

func loadConfig() config {
	return config{
		rabbitMQURL:      env("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		queueName:        env("SENSOR_QUEUE_NAME", "sensor.readings"),
		mongoURI:         env("MONGODB_URI", "mongodb://localhost:27017"),
		mongoDatabase:    env("MONGODB_DATABASE", "sensor_monitoring"),
		mongoCollection:  env("MONGODB_COLLECTION", "sensor_readings"),
		operationTimeout: 10 * time.Second,
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func handleDelivery(cfg config, collection *mongo.Collection, delivery amqp.Delivery) {
	event, err := decodeReadingEvent(delivery.Body)
	if err != nil {
		log.Printf("reject invalid message: %v", err)
		_ = delivery.Reject(false)
		return
	}

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

	ctx, cancel := context.WithTimeout(context.Background(), cfg.operationTimeout)
	defer cancel()

	if _, err := collection.InsertOne(ctx, doc); err != nil {
		log.Printf("mongodb insert failed, requeueing message: %v", err)
		_ = delivery.Nack(false, true)
		return
	}

	if err := delivery.Ack(false); err != nil {
		log.Printf("ack message: %v", err)
		return
	}

	log.Printf("stored reading eventId=%s sensorId=%s sensorType=%s", event.EventID, event.SensorID, event.SensorType)
}

func decodeReadingEvent(body []byte) (readingEvent, error) {
	var event readingEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return readingEvent{}, err
	}
	if event.EventID == "" {
		return readingEvent{}, errors.New("eventId is required")
	}
	if event.SensorID == "" {
		return readingEvent{}, errors.New("sensorId is required")
	}
	if event.SensorType == "" {
		return readingEvent{}, errors.New("sensorType is required")
	}
	if event.Unit == "" {
		return readingEvent{}, errors.New("unit is required")
	}
	if event.Timestamp.IsZero() {
		return readingEvent{}, errors.New("timestamp is required")
	}
	if event.ReceivedAt.IsZero() {
		return readingEvent{}, errors.New("receivedAt is required")
	}
	return event, nil
}
