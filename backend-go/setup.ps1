# Run this after installing Go from https://go.dev/dl/
Write-Host '=== Setting up Go Backend ===' -ForegroundColor Cyan

# Download all dependencies
go mod tidy
if ( -ne 0) { Write-Error 'go mod tidy failed'; exit 1 }

Write-Host '>> Dependencies downloaded' -ForegroundColor Green

# Copy env file
if (-not (Test-Path '.env')) {
    Copy-Item '.env.example' '.env'
    Write-Host '>> .env created from .env.example — add your DB connection string' -ForegroundColor Yellow
}

Write-Host ''
Write-Host '=== How to run ===' -ForegroundColor Cyan
Write-Host '  go run ./cmd/server   # starts server on :8080'
Write-Host '  go test ./...         # run all tests'
