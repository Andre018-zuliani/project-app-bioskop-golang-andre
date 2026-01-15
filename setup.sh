#!/bin/bash
# Setup script for Cinema Booking System

echo "Setting up Cinema Booking System..."

# Create database
echo "Creating database..."
createdb bioskop_db

# Apply schema
echo "Applying database schema..."
psql bioskop_db < db/schema.sql

# Download dependencies
echo "Downloading Go dependencies..."
go mod download

# Build the application
echo "Building the application..."
go build -o bin/app cmd/main/main.go
go build -o bin/seeder cmd/seeder/main.go

# Seed initial data
echo "Seeding initial data..."
./bin/seeder

echo "Setup complete! Run with: go run cmd/main/main.go"
