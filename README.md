# IntelliHR — AI-Powered HR Management System

## Tech Stack
| Layer | Technology |
|---|---|
| Web Frontend | React.js (Vite) + TypeScript |
| Mobile App | React Native (Expo) |
| Main Backend | Go (Gin) + PostgreSQL |
| AI Backend | Python FastAPI + LangChain |
| Database | PostgreSQL 16 + pgvector |
| Cache | Redis |
| Storage | MinIO (S3-compatible) |

## Project Structure
\\\
HRMS/
├── backend-go/          # Go main backend
├── backend-ai/          # FastAPI AI service
├── frontend-web/        # React.js web app (TBD)
├── frontend-mobile/     # React Native app (TBD)
├── database/            # Migrations & seeds
├── infrastructure/      # Docker, K8s, Terraform
└── docker-compose.yml   # Local dev environment
\\\

## Quick Start
\\\ash
# Start all infra services
docker-compose up postgres redis minio -d

# Run Go backend
cd backend-go && go run ./cmd/server

# Run FastAPI AI backend
cd backend-ai && uvicorn app.main:app --reload --port 8001
\\\

## Documentation
- Requirements: HRMS_Requirements.txt
- Database Schema: HRMS_Database_Schema.txt
