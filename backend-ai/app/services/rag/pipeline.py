"""
RAG Pipeline using LangChain + pgvector + OpenAI
"""
from langchain_openai import ChatOpenAI, OpenAIEmbeddings
from langchain_postgres import PGVector
from langchain.chains import RetrievalQA
from langchain.prompts import PromptTemplate
from app.core.config import settings

SYSTEM_PROMPT = """You are an HR assistant. Answer ONLY using the provided context.
If the answer is not in the context, say "I don't have information on that. Please contact HR."
Context: {context}
Question: {question}"""

def get_retriever(org_id: str):
    embeddings = OpenAIEmbeddings(model=settings.EMBEDDING_MODEL, api_key=settings.OPENAI_API_KEY)
    store = PGVector(
        embeddings=embeddings,
        collection_name=f"kb_chunks",
        connection=settings.DATABASE_URL.replace("+asyncpg", ""),
        use_jsonb=True,
    )
    return store.as_retriever(
        search_kwargs={"k": 5, "filter": {"org_id": org_id}}
    )

def get_qa_chain(org_id: str) -> RetrievalQA:
    llm = ChatOpenAI(model=settings.LLM_MODEL, api_key=settings.OPENAI_API_KEY, temperature=0.1)
    prompt = PromptTemplate(input_variables=["context","question"], template=SYSTEM_PROMPT)
    return RetrievalQA.from_chain_type(
        llm=llm,
        chain_type="stuff",
        retriever=get_retriever(org_id),
        chain_type_kwargs={"prompt": prompt},
        return_source_documents=True,
    )
