@echo off
REM Cinema Booking System - Test Runner Script (Windows)

echo ======================================
echo Cinema Booking System - Running Tests
echo ======================================
echo.

REM Run unit tests
echo Running unit tests for services...
echo.

go test ./internal/services/... -v -cover

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ======================================
    echo [32m✓ All tests passed successfully![0m
    echo ======================================
) else (
    echo.
    echo ======================================
    echo [31m✗ Some tests failed[0m
    echo ======================================
    exit /b 1
)

echo.
REM Generate coverage report
echo Generating coverage report...
go test ./internal/services/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

echo.
echo [32m✓ Coverage report generated: coverage.html[0m
echo.
echo ======================================
echo Test execution completed!
echo ======================================

pause
