# Temperature Monitoring & Event Notification System

## Overview

This project simulates an IoT-style temperature monitoring system. Virtual sensors generate temperature data, which is sent to a backend API, processed asynchronously via a message queue, stored in a database, and optionally used to trigger notifications to a mobile application.

The system is designed to demonstrate event-driven architecture, cloud integration, and real-time data processing.

---

## Features

* Simulated temperature sensors
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

Sensor Simulator → API → Message Queue → Worker → Database → Notification Service → Mobile App

---

### Components

#### 1. Sensor Simulator

* Generates temperature values at fixed intervals
* Simulates IoT devices
* Sends data via HTTP requests to backend API

---

#### 2. Backend API

* Receives temperature readings
* Validates incoming data
* Publishes events to message queue
* Built using FastAPI

---

#### 3. Message Queue

* Decouples ingestion from processing
* Ensures reliability and scalability
* Handles asynchronous event processing

---

#### 4. Worker Service

* Consumes messages from queue
* Applies business rules (threshold checks)
* Writes processed data to database
* Triggers alert events if needed

---

#### 5. Database

* Stores temperature readings and events
* Schema-free document storage

Example record:

```json
{
  "deviceId": "sensor-1",
  "temperature": 23.5,
  "timestamp": "2026-06-11T12:00:00Z"
}
```

---

#### 6. Notification Service

* Sends alerts when thresholds are exceeded
* Integrates with push notification provider

---

#### 7. Mobile App

* Displays live temperature data
* Shows historical trends
* Receives push notifications

---

## Example Event Flow

1. Sensor generates temperature: `28.4°C`
2. API receives data
3. Event is published to queue
4. Worker processes event
5. Data stored in database
6. If temperature > threshold → alert triggered
7. Notification sent to mobile app

---

## Technologies

* Backend: FastAPI / NestJS
* Database: MongoDB
* Message Queue: RabbitMQ or Kafka
* Mobile: React Native
* Cloud: MongoDB Atlas / Azure
* Notifications: Firebase Cloud Messaging

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
