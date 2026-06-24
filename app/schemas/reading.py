from datetime import datetime

from pydantic import BaseModel, ConfigDict, Field


class ReadingCreate(BaseModel):
    model_config = ConfigDict(populate_by_name=True)

    sensor_id: str = Field(..., alias="sensorId", min_length=1)
    sensor_type: str = Field(..., alias="sensorType", min_length=1)
    value: float
    unit: str = Field(..., min_length=1)
    timestamp: datetime | None = None


class ReadingResponse(BaseModel):
    model_config = ConfigDict(populate_by_name=True)

    sensor_id: str = Field(..., alias="sensorId")
    sensor_type: str = Field(..., alias="sensorType")
    value: float
    unit: str
    timestamp: datetime


class AcceptedReadingResponse(BaseModel):
    model_config = ConfigDict(populate_by_name=True)

    status: str
    event_id: str = Field(..., alias="eventId")
    reading: ReadingResponse
    received_at: datetime = Field(..., alias="receivedAt")
