from fastapi import APIRouter, status

from schemas.reading import AcceptedReadingResponse, ReadingCreate
from services.readings import accept_reading


router = APIRouter()


@router.post(
    "/readings",
    response_model=AcceptedReadingResponse,
    status_code=status.HTTP_202_ACCEPTED,
)
def create_reading(create: ReadingCreate) -> AcceptedReadingResponse:
    return accept_reading(create)
