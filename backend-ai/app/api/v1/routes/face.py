from fastapi import APIRouter, UploadFile, File, Depends, HTTPException
from sqlalchemy.ext.asyncio import AsyncSession
from app.db.session import get_db
from app.services.face.recognition import face_service
from app.core.security import verify_internal

router = APIRouter(dependencies=[Depends(verify_internal)])

@router.post("/enroll/{emp_id}")
async def enroll_face(emp_id: str, file: UploadFile = File(...), db: AsyncSession = Depends(get_db)):
    image_bytes = await file.read()
    try:
        await face_service.enroll(emp_id, image_bytes, db)
        return {"success": True, "message": "Face enrolled successfully"}
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))

@router.post("/verify")
async def verify_face(file: UploadFile = File(...), db: AsyncSession = Depends(get_db)):
    image_bytes = await file.read()
    result = await face_service.verify(image_bytes, db)
    return result

@router.post("/liveness-check")
async def liveness_check(file: UploadFile = File(...)):
    # TODO: implement liveness detection model
    return {"is_live": True, "confidence": 0.95}
