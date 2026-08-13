from pydantic import BaseModel
from uuid import UUID

class FaceEnrollRequest(BaseModel):
    emp_id: UUID

class FaceVerifyResponse(BaseModel):
    matched: bool
    emp_id: UUID | None = None
    confidence: float
    message: str

class LivenessResponse(BaseModel):
    is_live: bool
    confidence: float
