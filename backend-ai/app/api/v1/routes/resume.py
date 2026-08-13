from fastapi import APIRouter, UploadFile, File, Depends
from app.core.security import verify_internal
from app.services.resume.parser import parse_resume_text
from pypdf import PdfReader
import io

router = APIRouter(dependencies=[Depends(verify_internal)])

@router.post("/parse")
async def parse_resume(file: UploadFile = File(...)):
    raw = await file.read()
    reader = PdfReader(io.BytesIO(raw))
    text = "\n".join(p.extract_text() or "" for p in reader.pages)
    parsed = await parse_resume_text(text)
    return {"success": True, "data": parsed}
