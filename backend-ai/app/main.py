import logging
from contextlib import asynccontextmanager
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from app.api.v1.routes import face, chat, resume, kb
from app.core.config import settings
from app.db.session import check_db_connection

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup
    logger.info(f'Starting HRMS AI Service (env: {settings.APP_ENV})')
    db_ok = await check_db_connection()
    if not db_ok:
        logger.warning('DB unavailable at startup — some features may fail')
    yield
    # Shutdown
    logger.info('Shutting down HRMS AI Service')


app = FastAPI(
    title='HRMS AI Service',
    description='Face Recognition + RAG Chatbot + Resume Parser',
    version='1.0.0',
    lifespan=lifespan,
    docs_url='/docs' if settings.APP_ENV != 'production' else None,
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=['*'],
    allow_methods=['*'],
    allow_headers=['*'],
)

app.include_router(face.router,   prefix='/ai/v1/face',   tags=['Face Recognition'])
app.include_router(chat.router,   prefix='/ai/v1/chat',   tags=['RAG Chatbot'])
app.include_router(resume.router, prefix='/ai/v1/resume', tags=['Resume Parser'])
app.include_router(kb.router,     prefix='/ai/v1/kb',     tags=['Knowledge Base'])


@app.get('/health')
async def health():
    return {'status': 'ok', 'service': 'hrms-ai', 'env': settings.APP_ENV}
