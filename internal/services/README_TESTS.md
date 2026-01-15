# Unit Tests - Cinema Booking System

## Overview

This directory contains unit tests for the Cinema Booking System's service and repository layers.

## Test Coverage

### Implemented Tests

1. **User Service Tests** (`user_service_test.go`)

   - User registration (success, username exists, email exists)
   - User login (success, invalid credentials, wrong password)
   - Token verification (success, invalid token)
   - User logout
   - Get user by ID

2. **Booking Service Tests** (`booking_service_test.go`)
   - Create booking (success, seat not found, seat already booked)
   - Get user bookings (success, empty result)
   - Update booking status

## Running Tests

### Run All Tests

```bash
# Windows
go test ./internal/services/... -v

# Linux/Mac
go test ./internal/services/... -v
```

### Run Specific Test File

```bash
go test -run TestRegisterUser ./internal/services/ -v
```

### Generate Coverage Report

```bash
# Generate coverage profile
go test ./internal/services/... -coverprofile=coverage.out

# View coverage in HTML
go tool cover -html=coverage.out -o coverage.html
```

### Using Test Scripts

```bash
# Windows
run_tests.bat

# Linux/Mac
./run_tests.sh
```

## Test Structure

Tests follow the AAA (Arrange-Act-Assert) pattern:

```go
func TestFunctionName(t *testing.T) {
    // Arrange: Set up test data and mocks
    mockRepo := new(MockRepository)
    service := NewService(mockRepo)

    // Act: Execute the function being tested
    result, err := service.DoSomething(context.Background(), input)

    // Assert: Verify the results
    assert.NoError(t, err)
    assert.NotNil(t, result)
    mockRepo.AssertExpectations(t)
}
```

## Mocking

Tests use the `testify/mock` package for creating mock implementations of repositories. This allows testing service logic in isolation without requiring a real database connection.

### Example Mock

```go
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.User), args.Error(1)
}
```

## Test Dependencies

Tests require the following packages:

- `github.com/stretchr/testify` - Testing utilities and mocks
- Standard Go `testing` package

Install dependencies:

```bash
go get github.com/stretchr/testify
go mod tidy
```

## Notes

- Tests are designed to run independently without external dependencies
- Mock repositories simulate database behavior
- Tests focus on business logic in the service layer
- Integration tests with real database would be implemented separately

## Future Improvements

- Add repository layer tests with test database
- Add handler layer tests
- Increase test coverage to >80%
- Add integration tests
- Add benchmark tests for performance-critical functions
- Add table-driven tests for better coverage

## Test Conventions

- Test files are named `*_test.go`
- Test functions start with `Test`
- Use descriptive test names: `TestFunctionName_Scenario`
- Group related tests together
- Use subtests for variants: `t.Run("scenario", func(t *testing.T) {...})`

## Continuous Integration

Tests can be integrated into CI/CD pipelines:

```yaml
# Example GitHub Actions
- name: Run tests
  run: go test ./... -v -coverprofile=coverage.out

- name: Upload coverage
  uses: codecov/codecov-action@v2
  with:
    files: ./coverage.out
```
