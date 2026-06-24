import logging

from fastapi import FastAPI

from api import health, readings

logging.basicConfig(
    level=logging.DEBUG,
    format="%(asctime)s %(levelname)s %(name)s: %(message)s"
)
logging.getLogger("pika").setLevel(logging.WARNING)

logger = logging.getLogger(__name__)

app = FastAPI(title="Sensor Monitoring API")

app.include_router(health.router)
app.include_router(readings.router)

logger.info("API started")
