# Start all local development services
Write-Host '=== Starting HRMS Dev Environment ===' -ForegroundColor Cyan

# Start infra (postgres + redis + minio)
Write-Host '>> Starting Docker services...' -ForegroundColor Yellow
docker-compose up postgres redis minio -d

Start-Sleep 5

# Run Go backend in new terminal
Write-Host '>> Starting Go backend on :8080...' -ForegroundColor Green
Start-Process powershell -ArgumentList '-NoExit', '-Command', "cd '\backend-go'; go run ./cmd/server"

# Run FastAPI backend in new terminal
Write-Host '>> Starting FastAPI AI backend on :8001...' -ForegroundColor Green
Start-Process powershell -ArgumentList '-NoExit', '-Command', "cd '\backend-ai'; .\.venv\Scripts\Activate.ps1; uvicorn app.main:app --reload --port 8001"

Write-Host ''
Write-Host 'Services running:' -ForegroundColor Cyan
Write-Host '  Go Backend  → http://localhost:8080'
Write-Host '  FastAPI AI  → http://localhost:8001'
Write-Host '  MinIO UI    → http://localhost:9001'
Write-Host '  DB          → localhost:5432 (hrms_dev)'
