from sqlalchemy import Column, String, Integer, Text, DateTime, JSON
from sqlalchemy.dialects.postgresql import UUID
from pgvector.sqlalchemy import Vector
from app.db.base import Base
import uuid

class KBChunk(Base):
    __tablename__ = "kb_chunks"

    id          = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    doc_id      = Column(UUID(as_uuid=True), nullable=False)
    org_id      = Column(UUID(as_uuid=True), nullable=False)
    content     = Column(Text, nullable=False)
    embedding   = Column(Vector(1536), nullable=False)
    chunk_index = Column(Integer)
    token_count = Column(Integer)
    metadata    = Column(JSON)
    created_at  = Column(DateTime(timezone=True))
