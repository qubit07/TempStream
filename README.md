# Sensor Monitoring & Event Notification System

## Overview

This project simulates an IoT-style sensor monitoring system. Virtual sensors generate readings, which are sent to a backend API, processed asynchronously via a message queue, stored in a database, and optionally used to trigger notifications to a mobile application.

The system is designed to demonstrate event-driven architecture, cloud integration, and real-time data processing.

---

## Features

* Simulated sensors
* REST API for data ingestion
* Message queue-based processing
* Event-driven backend architecture
* Cloud-ready database storage (MongoDB)
* Rule-based alert generation
* Mobile push notifications
* Scalable worker-based processing

---

## System Architecture

### High-Level Flow

Sensor Simulator -> API -> RabbitMQ -> Go Worker -> MongoDB -> Notification Service -> Mobile App

---

### Components

#### 1. Sensor Simulator

* Generates sensor values at fixed intervals
* Simulates IoT devices
* Sends data via HTTP requests to backend API

---

#### 2. Backend API

* Receives sensor readings
* Validates incoming data
* Publishes reading events to RabbitMQ
* Returns `202 Accepted` after the event is queued
* Built using FastAPI

---

#### 3. Message Queue

* Decouples ingestion from processing
* Ensures reliability and scalability
* Handles asynchronous event processing
* Uses RabbitMQ queue: `sensor.readings`

---

#### 4. Worker Service

* Consumes messages from queue
* Applies business rules (threshold checks)
* Writes processed data to database
* Triggers alert events if needed
* Built using Go

---

#### 5. Database

* Stores sensor readings and events
* Schema-free document storage

Example record:

```json
{
  "sensorId": "sensor-1",
  "sensorType": "temperature",
  "value": 23.5,
  "unit": "C",
  "timestamp": "2026-06-11T12:00:00Z"
}
```

---

#### 6. Notification Service

* Sends alerts when thresholds are exceeded
* Integrates with push notification provider

---

#### 7. Mobile App

* Displays live sensor data
* Shows historical trends
* Receives push notifications

---

## Example Event Flow

1. Sensor generates a reading: `28.4°C`
2. API receives data
3. Event is published to queue
4. Worker processes event
5. Data stored in database
6. If the reading exceeds a threshold → alert triggered
7. Notification sent to mobile app

---

## Technologies

* Backend: FastAPI / NestJS
* Database: MongoDB
* Message Queue: RabbitMQ
* Mobile: React Native
* Cloud: MongoDB Atlas / Azure
* Notifications: Firebase Cloud Messaging

---

## Backend Quick Start

Install dependencies:

```bash
pip install -r requirements.txt
```

Run the FastAPI server:

```bash
uvicorn main:app --reload
```

Check the API:

```bash
curl http://127.0.0.1:8000/health
```

Send a sensor reading:

```bash
curl -X POST http://127.0.0.1:8000/readings \
  -H "Content-Type: application/json" \
  -d '{"sensorId":"sensor-1","sensorType":"temperature","value":23.5,"unit":"C"}'
```

---

## Worker Quick Start

The API requires RabbitMQ to accept readings. The Go worker requires RabbitMQ and MongoDB.

Environment variables:

```bash
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
SENSOR_QUEUE_NAME=sensor.readings
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=sensor_monitoring
MONGODB_COLLECTION=sensor_readings
```

Run the worker:

```bash
cd worker
go run ./cmd/worker
```

The API publishes reading events with this contract:

```json
{
  "eventId": "uuid",
  "sensorId": "sensor-1",
  "sensorType": "temperature",
  "value": 23.5,
  "unit": "C",
  "timestamp": "2026-06-11T12:00:00Z",
  "receivedAt": "2026-06-11T12:00:01Z"
}
```

---

## Project Goals

This project is designed to demonstrate:

* Event-driven system design
* Asynchronous processing
* Cloud-native architecture
* Scalable backend design
* Integration of multiple services
* Real-world IoT simulation

---

## Future Improvements

* Real IoT sensor integration
* Advanced analytics dashboard
* Stream processing (Kafka-based pipeline)
* Machine learning anomaly detection
* Multi-tenant support
