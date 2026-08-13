from pydantic import BaseModel
from uuid import UUID
from datetime import datetime

class ChatRequest(BaseModel):
    session_id: UUID
    message: str
    user_id: UUID
    org_id: UUID

class ChatResponse(BaseModel):
    answer: str
    sources: list[dict] = []
    session_id: UUID

class SessionCreate(BaseModel):
    user_id: UUID
    org_id: UUID
