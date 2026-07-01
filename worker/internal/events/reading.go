package events

import (
	"encoding/json"
	"errors"
	"time"
)

type ReadingEvent struct {
	EventID    string
	SensorID   string
	SensorType string
	Value      float64
	Unit       string
	Timestamp  time.Time
	ReceivedAt time.Time
}

func DecodeReadingEvent(body []byte) (ReadingEvent, error) {
	var event ReadingEvent

	if err := json.Unmarshal(body, &event); err != nil {
		return ReadingEvent{}, err
	}

	if err := event.Validate(); err != nil {
		return ReadingEvent{}, err
	}

	return event, nil
}

func (event ReadingEvent) Validate() (error) {
	if event.EventID == "" {
		return errors.New("eventId is required")
	}
	if event.SensorID == "" {
		return errors.New("sensorId is required")
	}
	if event.SensorType == "" {
		return errors.New("sensorType is required")
	}
	if event.Unit == "" {
		return errors.New("unit is required")
	}
	if event.Timestamp.IsZero() {
		return errors.New("timestamp is required")
	}
	if event.ReceivedAt.IsZero() {
		return errors.New("receivedAt is required")
	}
	return nil
}