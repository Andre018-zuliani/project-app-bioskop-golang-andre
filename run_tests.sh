#!/bin/bash

# Cinema Booking System - Test Runner Script

echo "======================================"
echo "Cinema Booking System - Running Tests"
echo "======================================"
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Run unit tests
echo "Running unit tests for services..."
echo ""

go test ./internal/services/... -v -cover

TEST_RESULT=$?

echo ""
echo "======================================"

if [ $TEST_RESULT -eq 0 ]; then
    echo -e "${GREEN}✓ All tests passed successfully!${NC}"
else
    echo -e "${RED}✗ Some tests failed${NC}"
    exit 1
fi

echo "======================================"
echo ""

# Generate coverage report
echo "Generating coverage report..."
go test ./internal/services/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

echo ""
echo -e "${GREEN}✓ Coverage report generated: coverage.html${NC}"
echo ""
echo "======================================"
echo "Test execution completed!"
echo "======================================"
