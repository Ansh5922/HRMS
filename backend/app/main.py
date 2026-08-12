from fastapi import FastAPI

app = FastAPI(
    title="Face Recognition Attendance System API",
    version="1.0.0",
    description="Multi-tenant HRMS and Face Recognition Attendance API",
)

@app.get("/")
def read_root():
    return {"message": "Welcome to Face Recognition Attendance System API"}
