from pydantic import BaseModel

class ParsedResume(BaseModel):
    name: str | None
    email: str | None
    phone: str | None
    skills: list[str] = []
    experience_years: float | None
    education: list[dict] = []
    work_history: list[dict] = []
    summary: str | None
