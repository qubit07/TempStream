package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tempstream-worker/internal/config"
	"tempstream-worker/internal/db"
	"tempstream-worker/internal/messaging"
	"tempstream-worker/internal/worker"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	mongoCtx, cancel := context.WithTimeout(ctx, cfg.OperationTimeout)
	defer cancel()

	client, err := mongo.Connect(mongoCtx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("mongo connect: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.OperationTimeout)
		defer cancel()
		_ = client.Disconnect(ctx)
	}()

	collection := client.
		Database(cfg.MongoDatabase).
		Collection(cfg.MongoCollection)

	store := db.NewReadingStore(collection)
	worker := worker.New(cfg, store)

	consumer, err := messaging.NewConsumer(messaging.Config{
		URL:      cfg.RabbitMQURL,
		Queue:    cfg.QueueName,
		Prefetch: 1,
	})
	if err != nil {
		log.Fatalf("rabbitmq init: %v", err)
	}
	defer consumer.Close()

	log.Println("worker started")

	consumer.Run(ctx, worker.HandleDelivery)

	log.Println("worker stopped")
}