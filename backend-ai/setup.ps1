# Run this to set up the FastAPI AI backend
Write-Host '=== Setting up FastAPI AI Backend ===' -ForegroundColor Cyan

# Create virtual environment
python -m venv .venv
if ( -ne 0) { Write-Error 'python venv failed'; exit 1 }

# Activate and install
.\.venv\Scripts\Activate.ps1
pip install --upgrade pip
pip install -r requirements.txt
if ( -ne 0) { Write-Error 'pip install failed'; exit 1 }

Write-Host '>> Dependencies installed' -ForegroundColor Green

# Copy env file
if (-not (Test-Path '.env')) {
    Copy-Item '.env.example' '.env'
    Write-Host '>> .env created from .env.example — add your DB connection string' -ForegroundColor Yellow
}

Write-Host ''
Write-Host '=== How to run ===' -ForegroundColor Cyan
Write-Host '  .\.venv\Scripts\Activate.ps1'
Write-Host '  uvicorn app.main:app --reload --port 8001'
Write-Host ''
Write-Host '=== Apply DB Migrations (AI schema only) ==='
Write-Host '  alembic upgrade head'
