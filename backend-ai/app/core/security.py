from fastapi import Header, HTTPException
from app.core.config import settings

async def verify_internal(x_internal_secret: str = Header(...)):
    """Service-to-service auth: Go backend must pass X-Internal-Secret header."""
    if x_internal_secret != settings.INTERNAL_SECRET:
        raise HTTPException(status_code=403, detail="Forbidden")
