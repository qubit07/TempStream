from os import getenv


RABBITMQ_URL = getenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
READINGS_QUEUE_NAME = getenv("READINGS_QUEUE_NAME", "sensor.readings")
