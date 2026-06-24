import json

import pika
from fastapi import HTTPException, status

from messaging.config import RABBITMQ_URL, READINGS_QUEUE_NAME


def publish_readings_event(event: dict[str, object]) -> None:
    try:
        parameters = pika.URLParameters(RABBITMQ_URL)
        connection = pika.BlockingConnection(parameters)
        channel = connection.channel()
        channel.queue_declare(queue=READINGS_QUEUE_NAME, durable=True)
        channel.basic_publish(
            exchange="",
            routing_key=READINGS_QUEUE_NAME,
            body=json.dumps(event).encode("utf-8"),
            properties=pika.BasicProperties(
                content_type="application/json",
                delivery_mode=pika.DeliveryMode.Persistent,
            ),
        )
        connection.close()
    except pika.exceptions.AMQPError as exc:
        raise HTTPException(
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
            detail="Message queue is unavailable",
        ) from exc
