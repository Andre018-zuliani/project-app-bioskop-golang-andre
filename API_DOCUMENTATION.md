# API Documentation - Cinema Booking System

## Base URL

```
http://localhost:8080
```

## Authentication

All protected endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer <token>
```

---

## Endpoints

### 1. User Authentication

#### Register User

```http
POST /api/register
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Response (201 Created):**

```json
{
  "id": 1,
  "username": "john_doe",
  "email": "john@example.com",
  "is_verified": false,
  "created_at": "2026-01-13T10:00:00Z",
  "updated_at": "2026-01-13T10:00:00Z"
}
```

---

#### Login User

```http
POST /api/login
Content-Type: application/json

{
  "username": "john_doe",
  "password": "securepassword123"
}
```

**Response (200 OK):**

```json
{
  "id": 1,
  "username": "john_doe",
  "email": "john@example.com",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

#### Logout User

```http
POST /api/logout
Authorization: Bearer <token>
```

**Response (200 OK):**

```json
{
  "message": "Logout successful"
}
```

---

#### Get User Profile

```http
GET /api/user/profile
Authorization: Bearer <token>
```

**Response (200 OK):**

```json
{
  "id": 1,
  "username": "john_doe",
  "email": "john@example.com",
  "is_verified": false,
  "created_at": "2026-01-13T10:00:00Z",
  "updated_at": "2026-01-13T10:00:00Z"
}
```

---

### 2. Cinema Management

#### Get All Cinemas (with Pagination)

```http
GET /api/cinemas?page=1&limit=10&city=Jakarta&name=
```

**Parameters:**

- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10, max: 100)
- `city` (optional): Filter by city
- `name` (optional): Filter by cinema name

**Response (200 OK):**

```json
{
  "data": [
    {
      "id": 1,
      "name": "CGV Cinemas - Jakarta",
      "location": "Blok M Plaza",
      "city": "Jakarta",
      "address": "Jl. Melawai No. 1, Blok M, Jakarta Selatan",
      "total_seats": 150,
      "image_url": "https://via.placeholder.com/300x200?text=CGV+Jakarta",
      "created_at": "2026-01-13T10:00:00Z",
      "updated_at": "2026-01-13T10:00:00Z"
    },
    {
      "id": 2,
      "name": "Cinemaxx - Surabaya",
      "location": "Pakuwon Indah",
      "city": "Surabaya",
      "address": "Jl. Raya Pakuwon Indah, Surabaya",
      "total_seats": 200,
      "image_url": "https://via.placeholder.com/300x200?text=Cinemaxx+Surabaya",
      "created_at": "2026-01-13T10:00:00Z",
      "updated_at": "2026-01-13T10:00:00Z"
    }
  ],
  "page": 1,
  "limit": 10,
  "total": 5,
  "total_pages": 1
}
```

---

#### Get Cinema Details

```http
GET /api/cinemas/{cinemaId}
```

**Response (200 OK):**

```json
{
  "id": 1,
  "name": "CGV Cinemas - Jakarta",
  "location": "Blok M Plaza",
  "city": "Jakarta",
  "address": "Jl. Melawai No. 1, Blok M, Jakarta Selatan",
  "total_seats": 150,
  "image_url": "https://via.placeholder.com/300x200?text=CGV+Jakarta",
  "created_at": "2026-01-13T10:00:00Z",
  "updated_at": "2026-01-13T10:00:00Z"
}
```

---

### 3. Seat Availability

#### Check Seat Availability

```http
GET /api/cinemas/{cinemaId}/seats?date=2026-01-20&time=19:00
```

**Parameters:**

- `date` (required): Date in format YYYY-MM-DD
- `time` (required): Time in format HH:MM (e.g., 19:00)

**Response (200 OK):**

```json
{
  "cinema_id": 1,
  "date": "2026-01-20",
  "time": "19:00",
  "available_seats": [
    {
      "id": 1,
      "cinema_id": 1,
      "seat_id": 1,
      "show_date": "2026-01-20T00:00:00Z",
      "show_time": "19:00",
      "is_available": true,
      "created_at": "2026-01-13T10:00:00Z",
      "updated_at": "2026-01-13T10:00:00Z",
      "seat": {
        "id": 1,
        "cinema_id": 1,
        "seat_number": "1A",
        "row_number": 1,
        "seat_type": "standard",
        "price": 50000,
        "created_at": "2026-01-13T10:00:00Z",
        "updated_at": "2026-01-13T10:00:00Z"
      }
    }
  ],
  "unavailable_seats": [],
  "total_available": 120,
  "total_unavailable": 30
}
```

---

### 4. Booking Management

#### Create Booking

```http
POST /api/booking
Authorization: Bearer <token>
Content-Type: application/json

{
  "cinema_id": 1,
  "seat_id": 5,
  "date": "2026-01-20",
  "time": "19:00",
  "payment_method": "Kartu Kredit"
}
```

**Response (201 Created):**

```json
{
  "id": 1,
  "cinema_id": 1,
  "seat_id": 5,
  "show_date": "2026-01-20T00:00:00Z",
  "show_time": "19:00",
  "total_price": 50000,
  "payment_method": "Kartu Kredit",
  "status": "pending",
  "payment_status": "pending",
  "created_at": "2026-01-13T10:00:00Z"
}
```

---

#### Get User Bookings

```http
GET /api/user/bookings?page=1&limit=10
Authorization: Bearer <token>
```

**Parameters:**

- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10)

**Response (200 OK):**

```json
{
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "cinema_id": 1,
      "seat_id": 5,
      "show_date": "2026-01-20T00:00:00Z",
      "show_time": "19:00",
      "booking_date": "2026-01-13T10:00:00Z",
      "status": "confirmed",
      "total_price": 50000,
      "payment_method": "Kartu Kredit",
      "payment_status": "paid",
      "created_at": "2026-01-13T10:00:00Z",
      "updated_at": "2026-01-13T10:30:00Z",
      "cinema": {
        "id": 1,
        "name": "CGV Cinemas - Jakarta",
        "location": "Blok M Plaza",
        "city": "Jakarta",
        "address": "Jl. Melawai No. 1, Blok M, Jakarta Selatan",
        "total_seats": 150,
        "image_url": "https://via.placeholder.com/300x200?text=CGV+Jakarta",
        "created_at": "2026-01-13T10:00:00Z",
        "updated_at": "2026-01-13T10:00:00Z"
      },
      "seat": {
        "id": 5,
        "cinema_id": 1,
        "seat_number": "1E",
        "row_number": 1,
        "seat_type": "standard",
        "price": 50000,
        "created_at": "2026-01-13T10:00:00Z",
        "updated_at": "2026-01-13T10:00:00Z"
      }
    }
  ],
  "page": 1,
  "limit": 10,
  "total": 1,
  "total_pages": 1
}
```

---

### 5. Payment Methods

#### Get Available Payment Methods

```http
GET /api/payment-methods
```

**Response (200 OK):**

```json
[
  {
    "id": 1,
    "name": "Kartu Kredit",
    "type": "credit_card",
    "is_active": true,
    "created_at": "2026-01-13T10:00:00Z",
    "updated_at": "2026-01-13T10:00:00Z"
  },
  {
    "id": 2,
    "name": "Kartu Debit",
    "type": "debit_card",
    "is_active": true,
    "created_at": "2026-01-13T10:00:00Z",
    "updated_at": "2026-01-13T10:00:00Z"
  },
  {
    "id": 3,
    "name": "Transfer Bank",
    "type": "transfer",
    "is_active": true,
    "created_at": "2026-01-13T10:00:00Z",
    "updated_at": "2026-01-13T10:00:00Z"
  }
]
```

---

#### Process Payment

```http
POST /api/pay
Authorization: Bearer <token>
Content-Type: application/json

{
  "booking_id": 1,
  "payment_method": "Kartu Kredit",
  "amount": 50000.0
}
```

**Response (201 Created):**

```json
{
  "id": 1,
  "booking_id": 1,
  "amount": 50000.0,
  "payment_method": "Kartu Kredit",
  "status": "success",
  "transaction_id": "TXN-1-1",
  "created_at": "2026-01-13T10:30:00Z"
}
```

---

### 6. Health Check

#### Health Status

```http
GET /health
```

**Response (200 OK):**

```json
{
  "status": "ok"
}
```

---

## Error Responses

### 400 Bad Request

```json
{
  "error": "Validation error: field must be required"
}
```

### 401 Unauthorized

```json
{
  "error": "Invalid or expired token"
}
```

### 404 Not Found

```json
{
  "error": "Cinema not found"
}
```

### 500 Internal Server Error

```json
{
  "error": "Failed to process request"
}
```

---

## HTTP Status Codes

| Code | Meaning                                 |
| ---- | --------------------------------------- |
| 200  | OK - Request successful                 |
| 201  | Created - Resource created successfully |
| 400  | Bad Request - Invalid input             |
| 401  | Unauthorized - Missing/invalid token    |
| 404  | Not Found - Resource not found          |
| 500  | Internal Server Error                   |

---

## Rate Limiting & Pagination

- Default page size: 10 items
- Maximum page size: 100 items
- Cinemas are sorted by name (ASC)
- Bookings are sorted by date (DESC)

---

## Testing with cURL

### Register and Login

```bash
# Register
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"pass123"}'

# Login and save token
TOKEN=$(curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"pass123"}' | jq -r '.token')

# Use token in requests
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/user/profile
```

---

## Postman Collection

Import the `Postman_Collection.json` file to quickly test all endpoints with:

- Pre-configured requests
- Environment variables
- Example payloads
- Response samples
