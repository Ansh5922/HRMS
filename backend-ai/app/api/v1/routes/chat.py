from fastapi import APIRouter, Depends
from sqlalchemy.ext.asyncio import AsyncSession
from app.db.session import get_db
from app.schemas.chat import ChatRequest, ChatResponse, SessionCreate
from app.services.rag.pipeline import get_qa_chain

router = APIRouter()

@router.post("/session")
async def create_session(body: SessionCreate, db: AsyncSession = Depends(get_db)):
    # TODO: persist session to chat_sessions table
    return {"session_id": "new-session-uuid", "message": "Session created"}

@router.post("/message", response_model=ChatResponse)
async def send_message(body: ChatRequest, db: AsyncSession = Depends(get_db)):
    chain = get_qa_chain(str(body.org_id))
    result = chain.invoke({"query": body.message})
    sources = [{"content": d.page_content[:200]} for d in result.get("source_documents", [])]
    # TODO: persist message to chat_messages table
    return ChatResponse(answer=result["result"], sources=sources, session_id=body.session_id)

@router.get("/history/{session_id}")
async def get_history(session_id: str, db: AsyncSession = Depends(get_db)):
    # TODO: fetch from chat_messages table
    return {"session_id": session_id, "messages": []}
