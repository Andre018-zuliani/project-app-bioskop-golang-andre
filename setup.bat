@echo off
REM Setup script for Cinema Booking System (Windows)

echo Setting up Cinema Booking System...

REM Create database
echo Creating database...
psql -U postgres -c "CREATE DATABASE bioskop_db;"

REM Apply schema
echo Applying database schema...
psql -U postgres -d bioskop_db -f db/schema.sql

REM Download dependencies
echo Downloading Go dependencies...
go mod download

REM Build the application
echo Building the application...
go build -o bin/app.exe cmd/main/main.go
go build -o bin/seeder.exe cmd/seeder/main.go

REM Seed initial data
echo Seeding initial data...
bin/seeder.exe

echo Setup complete! Run with: go run cmd/main/main.go
