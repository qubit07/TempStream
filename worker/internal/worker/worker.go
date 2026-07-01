package worker

import (
	"context"
	"tempstream-worker/internal/config"
	"tempstream-worker/internal/db"
	"tempstream-worker/internal/events"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Worker struct {
    store *db.ReadingStore
    cfg   config.Config
}

func New(
    cfg config.Config,
    store *db.ReadingStore,
) *Worker {
    return &Worker{
        cfg: cfg,
        store: store,
    }
}

func (w *Worker) HandleDelivery(
ctx context.Context,
    delivery amqp.Delivery,
) error {
    event, err := events.DecodeReadingEvent(delivery.Body)

    if err != nil {
        _ = delivery.Reject(false)
        return nil
    }

	opCtx, cancel := context.WithTimeout(
		ctx,
		w.cfg.OperationTimeout,
	)
    defer cancel()

    err = w.store.Save(opCtx, event)

    if err != nil {
        _ = delivery.Nack(false, false)
        return err
    }

	if err := delivery.Ack(false); err != nil {
		return err
	}

	return nil
}