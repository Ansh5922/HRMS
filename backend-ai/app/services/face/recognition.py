"""
Face Recognition Service
Uses InsightFace / DeepFace for embedding generation.
pgvector cosine similarity for identity matching.
"""
import numpy as np
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import text
from app.core.config import settings

class FaceRecognitionService:
    def __init__(self):
        # Lazy-load model to avoid startup delay
        self._model = None

    def _get_model(self):
        if self._model is None:
            import insightface
            self._model = insightface.app.FaceAnalysis(name="buffalo_l")
            self._model.prepare(ctx_id=-1)   # -1 = CPU
        return self._model

    def get_embedding(self, image_bytes: bytes) -> np.ndarray | None:
        import cv2
        nparr = np.frombuffer(image_bytes, np.uint8)
        img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
        model = self._get_model()
        faces = model.get(img)
        if not faces:
            return None
        return faces[0].normed_embedding  # 512-dim

    async def enroll(self, emp_id: str, image_bytes: bytes, db: AsyncSession):
        embedding = self.get_embedding(image_bytes)
        if embedding is None:
            raise ValueError("No face detected in the image")
        vec = embedding.tolist()
        await db.execute(text("""
            INSERT INTO face_embeddings (emp_id, embedding, model_version)
            VALUES (:emp_id, :embedding, :model)
            ON CONFLICT (emp_id) DO UPDATE
            SET embedding = EXCLUDED.embedding, updated_at = NOW()
        """), {"emp_id": emp_id, "embedding": vec, "model": settings.FACE_MODEL})
        await db.commit()

    async def verify(self, image_bytes: bytes, db: AsyncSession) -> dict:
        embedding = self.get_embedding(image_bytes)
        if embedding is None:
            return {"matched": False, "emp_id": None, "confidence": 0.0}
        vec = embedding.tolist()
        threshold = settings.FACE_CONFIDENCE_THRESHOLD
        result = await db.execute(text("""
            SELECT emp_id, 1 - (embedding <=> :query_vec::vector) AS confidence
            FROM face_embeddings
            ORDER BY embedding <=> :query_vec::vector
            LIMIT 1
        """), {"query_vec": vec})
        row = result.fetchone()
        if row and row.confidence >= threshold:
            return {"matched": True, "emp_id": str(row.emp_id), "confidence": float(row.confidence)}
        return {"matched": False, "emp_id": None, "confidence": float(row.confidence) if row else 0.0}

face_service = FaceRecognitionService()
