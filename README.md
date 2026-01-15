# Cinema Booking System API

A RESTful API built with Go for cinema seat booking and ticket management.

## Features

- User Registration and Authentication with JWT
- Cinema Selection with Pagination
- Seat Availability Checking
- Booking Management
- Payment Processing
- User Booking History

## Tech Stack

- **Language**: Go 1.21
- **Framework**: Chi Router
- **Database**: PostgreSQL
- **Authentication**: JWT (golang-jwt)
- **Validation**: go-playground/validator
- **Logging**: Zap Logger
- **Configuration**: Viper
- **Password Hashing**: bcrypt

## Project Structure

```
.
├── cmd/
│   ├── main/          # Main application
│   └── seeder/        # Database seeder
├── internal/
│   ├── config/        # Configuration management
│   ├── handlers/      # HTTP handlers
│   ├── middleware/    # HTTP middleware
│   ├── models/        # Data models
│   ├── repositories/  # Data access layer
│   └── services/      # Business logic layer
├── db/
│   └── schema.sql     # Database schema
├── .env               # Environment variables
├── go.mod             # Go module file
└── README.md          # This file
```

## Installation

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher

### Setup

1. Clone the repository and navigate to the project directory
2. Create `.env` file with necessary configuration:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=bioskop_db
SERVER_PORT=8080
SERVER_ENV=development
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
```

3. Create database:

```bash
createdb bioskop_db
```

4. Apply schema:

```bash
psql bioskop_db < db/schema.sql
```

5. Download dependencies:

```bash
go mod download
```

6. Seed initial data (optional):

```bash
go run cmd/seeder/main.go
```

## Running the Application

```bash
go run cmd/main/main.go
```

The server will start at `http://localhost:8080`

## API Endpoints

### Authentication

- `POST /api/register` - Register new user
- `POST /api/login` - Login user
- `POST /api/logout` - Logout user (requires auth)

### Cinema

- `GET /api/cinemas` - Get all cinemas (with pagination)
- `GET /api/cinemas/{cinemaId}` - Get cinema details

### Seats

- `GET /api/cinemas/{cinemaId}/seats?date=YYYY-MM-DD&time=HH:MM` - Get seat availability

### Booking

- `POST /api/booking` - Create booking (requires auth)
- `GET /api/user/bookings` - Get user booking history (requires auth)

### Payment

- `GET /api/payment-methods` - Get available payment methods
- `POST /api/pay` - Process payment (requires auth)

### User

- `GET /api/user/profile` - Get user profile (requires auth)

## Authentication

API endpoints that require authentication need the JWT token in the Authorization header:

```
Authorization: Bearer <token>
```

## Example Requests

### Register User

```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "securepassword123"
  }'
```

### Login User

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "securepassword123"
  }'
```

### Get Cinemas

```bash
curl -X GET "http://localhost:8080/api/cinemas?page=1&limit=10&city=Jakarta"
```

### Check Seat Availability

```bash
curl -X GET "http://localhost:8080/api/cinemas/1/seats?date=2026-01-20&time=19:00"
```

### Create Booking

```bash
curl -X POST http://localhost:8080/api/booking \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "cinema_id": 1,
    "seat_id": 1,
    "date": "2026-01-20",
    "time": "19:00",
    "payment_method": "Kartu Kredit"
  }'
```

### Process Payment

```bash
curl -X POST http://localhost:8080/api/pay \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "booking_id": 1,
    "payment_method": "Kartu Kredit",
    "amount": 50000.0
  }'
```

## Testing with Postman

Import the included Postman collection to test all endpoints. The collection includes:

- Pre-configured environment variables
- Example requests for all endpoints
- Authentication token management

## Architecture

The application follows the Clean Code architecture with clear separation of concerns:

1. **Handlers Layer**: HTTP request/response handling
2. **Services Layer**: Business logic implementation
3. **Repositories Layer**: Data access and database operations
4. **Models Layer**: Data structures

## Logging

The application uses Zap Logger for structured logging. Logs include:

- Request/response information
- Error details
- User actions (registration, login, booking, payment)

## Error Handling

The API returns appropriate HTTP status codes:

- `200 OK` - Successful GET request
- `201 Created` - Successful resource creation
- `400 Bad Request` - Invalid input or validation error
- `401 Unauthorized` - Missing or invalid authentication
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Future Enhancements

- Email OTP verification
- Unit tests and integration tests
- Async payment processing with goroutines
- Email notifications for bookings
- Movie information integration
- Advanced filtering and search

## Contributing

Feel free to submit issues and enhancement requests.

## License

This project is open source and available under the MIT License.
