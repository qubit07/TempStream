from datetime import datetime, timezone
import logging
from uuid import uuid4

from messaging.publisher import publish_readings_event
from schemas.reading import AcceptedReadingResponse, ReadingCreate, ReadingResponse

logger = logging.getLogger(__name__)

def accept_reading(create: ReadingCreate) -> AcceptedReadingResponse:
    timestamp = create.timestamp or datetime.now(timezone.utc)
    received_at = datetime.now(timezone.utc)
    event_id = str(uuid4())

    payload = {
            "eventId": event_id,
            "sensorId": create.sensor_id,
            "sensorType": create.sensor_type,
            "value": create.value,
            "unit": create.unit,
            "timestamp": timestamp.isoformat(),
            "receivedAt": received_at.isoformat(),
        }

    logger.debug("Publishing event_id=%s payload=%s", event_id, payload)

    publish_readings_event(payload)

    logger.info("Published event_id=%s sensor_id=%s", event_id, create.sensor_id)


    return AcceptedReadingResponse(
        status="accepted",
        event_id=event_id,
        reading=ReadingResponse(
            sensor_id=create.sensor_id,
            sensor_type=create.sensor_type,
            value=create.value,
            unit=create.unit,
            timestamp=timestamp,
        ),
        received_at=received_at,
    )
