package config

import (
	"os"
	"time"
)

type Config struct {
	RabbitMQURL      string
	QueueName        string
	Durable bool
	AutoDelete bool
	Exclusive bool
	NoWait bool
	AutoAck bool
	NoLocal bool

	MongoURI         string
	MongoDatabase    string
	MongoCollection  string

	OperationTimeout time.Duration
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func Load() Config {
	cfg := Config{
		RabbitMQURL:      env("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		QueueName:        env("SENSOR_QUEUE_NAME", "sensor.readings"),
		MongoURI:         env("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDatabase:    env("MONGODB_DATABASE", "sensor_monitoring"),
		MongoCollection:  env("MONGODB_COLLECTION", "sensor_readings"),
		OperationTimeout: 10 * time.Second,
	}
	return cfg
}



