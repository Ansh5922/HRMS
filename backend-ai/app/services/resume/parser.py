"""
AI Resume Parser using GPT-4o
"""
import json
from openai import AsyncOpenAI
from app.core.config import settings

client = AsyncOpenAI(api_key=settings.OPENAI_API_KEY)

PARSE_PROMPT = """Extract the following fields from the resume text and return as JSON:
name, email, phone, skills (list), experience_years (float), summary,
education (list of {degree, institution, year}),
work_history (list of {company, title, start_year, end_year, description}).
Resume:
{text}"""

async def parse_resume_text(text: str) -> dict:
    response = await client.chat.completions.create(
        model=settings.LLM_MODEL,
        messages=[{"role": "user", "content": PARSE_PROMPT.format(text=text)}],
        response_format={"type": "json_object"},
    )
    return json.loads(response.choices[0].message.content)
