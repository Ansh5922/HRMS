import uuid
import os
import aiofiles
from fastapi import APIRouter, UploadFile, File, Depends, Form, HTTPException
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import text
from app.db.session import get_db
from app.core.security import verify_internal
from app.services.rag.ingestion import ingest_document

router = APIRouter(dependencies=[Depends(verify_internal)])

UPLOAD_DIR = "/tmp/hrms_uploads"
os.makedirs(UPLOAD_DIR, exist_ok=True)


@router.post("/upload")
async def upload_document(
    org_id: str = Form(...),
    title: str = Form(...),
    category: str = Form("general"),
    file: UploadFile = File(...),
    db: AsyncSession = Depends(get_db),
):
    doc_id = str(uuid.uuid4())
    ext = file.filename.split(".")[-1].lower()
    if ext not in ("pdf", "docx", "txt"):
        raise HTTPException(status_code=400, detail="Only PDF, DOCX, TXT files supported")

    # Save file temporarily
    tmp_path = f"{UPLOAD_DIR}/{doc_id}.{ext}"
    async with aiofiles.open(tmp_path, "wb") as f:
        await f.write(await file.read())

    # Persist document record
    await db.execute(text("""
        INSERT INTO kb_documents (id, org_id, title, file_url, file_type, category)
        VALUES (:id, :org_id, :title, :url, :ftype, :cat)
    """), {"id": doc_id, "org_id": org_id, "title": title,
           "url": tmp_path, "ftype": ext, "cat": category})
    await db.commit()

    # Run ingestion pipeline
    await ingest_document(doc_id, org_id, tmp_path, ext, db)
    return {"success": True, "doc_id": doc_id, "message": "Document ingested successfully"}


@router.get("/documents")
async def list_documents(org_id: str, db: AsyncSession = Depends(get_db)):
    result = await db.execute(text(
        "SELECT id, title, file_type, category, chunk_count, ingested_at FROM kb_documents "
        "WHERE org_id = :org_id AND is_active = TRUE ORDER BY created_at DESC"
    ), {"org_id": org_id})
    rows = result.mappings().all()
    return {"documents": [dict(r) for r in rows]}


@router.delete("/documents/{doc_id}")
async def delete_document(doc_id: str, db: AsyncSession = Depends(get_db)):
    await db.execute(text(
        "UPDATE kb_documents SET is_active = FALSE WHERE id = :id"
    ), {"id": doc_id})
    await db.execute(text("DELETE FROM kb_chunks WHERE doc_id = :id"), {"id": doc_id})
    await db.commit()
    return {"success": True}
