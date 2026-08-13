from sqlalchemy import Column, String, DateTime, Text
from sqlalchemy.dialects.postgresql import UUID
from pgvector.sqlalchemy import Vector
from app.db.base import Base
import uuid

class FaceEmbedding(Base):
    __tablename__ = "face_embeddings"

    id            = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    emp_id        = Column(UUID(as_uuid=True), nullable=False, unique=True)
    embedding     = Column(Vector(512), nullable=False)
    model_version = Column(String(100))
    enrolled_at   = Column(DateTime(timezone=True))
