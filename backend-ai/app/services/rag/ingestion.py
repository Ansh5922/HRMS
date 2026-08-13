"""
Document ingestion pipeline:
  PDF/DOCX → text → chunks → OpenAI embeddings → INSERT INTO kb_chunks (pgvector)
"""
import uuid
import logging
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import text
from langchain_openai import OpenAIEmbeddings
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain_community.document_loaders import PyPDFLoader, Docx2txtLoader
from app.core.config import settings

logger = logging.getLogger(__name__)

splitter = RecursiveCharacterTextSplitter(chunk_size=800, chunk_overlap=100)
embedder = OpenAIEmbeddings(model=settings.EMBEDDING_MODEL, api_key=settings.OPENAI_API_KEY)


async def ingest_document(doc_id: str, org_id: str, file_path: str, file_type: str, db: AsyncSession):
    """Load file, chunk, embed and store in kb_chunks."""
    if file_type == 'pdf':
        loader = PyPDFLoader(file_path)
    else:
        loader = Docx2txtLoader(file_path)

    docs = loader.load()
    chunks = splitter.split_documents(docs)
    logger.info(f'Ingesting {len(chunks)} chunks for doc {doc_id}')

    texts = [c.page_content for c in chunks]
    embeddings = embedder.embed_documents(texts)

    for i, (chunk, emb) in enumerate(zip(chunks, embeddings)):
        chunk_id = str(uuid.uuid4())
        await db.execute(text("""
            INSERT INTO kb_chunks (id, doc_id, org_id, content, embedding, chunk_index, token_count, metadata)
            VALUES (:id, :doc_id, :org_id, :content, :emb::vector, :idx, :tokens, :meta)
        """), {
            "id": chunk_id,
            "doc_id": doc_id,
            "org_id": org_id,
            "content": chunk.page_content,
            "emb": emb,
            "idx": i,
            "tokens": len(chunk.page_content.split()),
            "meta": {"source": chunk.metadata.get("source", ""), "page": chunk.metadata.get("page", 0)},
        })

    await db.execute(text(
        "UPDATE kb_documents SET chunk_count = :n, ingested_at = NOW() WHERE id = :id"
    ), {"n": len(chunks), "id": doc_id})
    await db.commit()
    logger.info(f'Ingestion complete for doc {doc_id}')
