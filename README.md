# Chirpy 🐦

A fully-featured web server built with Go that provides a Twitter-like social media API. Chirpy allows users to create accounts, post short messages (chirps), and interact with a modern REST API.

## What Chirpy Does

Chirpy is a social media backend API that enables:

-   **User Management**: Registration, authentication, and profile updates
-   **Content Creation**: Post and manage short messages called "chirps"
-   **Authentication**: JWT-based authentication with refresh tokens
-   **Premium Features**: Integration with Polka (imaginary) for user upgrades
-   **Admin Tools**: User management and metrics tracking

## Why You Should Care

-   **Modern Go Architecture**: Built with clean separation of concerns and best practices
-   **Production Ready**: Includes proper authentication, database integration, and error handling
-   **RESTful API**: Well-designed endpoints following REST conventions
-   **Scalable**: Uses PostgreSQL for data persistence and JWT for stateless authentication
-   **Educational**: Perfect example of a real-world Go web application

## Installation & Setup

### Prerequisites

-   **Go 1.25.0 or later** - [Download Go](https://golang.org/dl/)
-   **PostgreSQL 12 or later** - [Download PostgreSQL](https://www.postgresql.org/download/)

### 1. Clone and Install Dependencies

```bash
git clone https://github.com/HemahWeb/chirpy
go mod download
```

### 2. Database Setup

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE chirpy_db;

# Create user (optional)
CREATE USER chirpy_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE chirpy_db TO chirpy_user;
```

### 3. Environment Configuration

Create a `.env` file in the project root:

```bash
# Database Configuration
DB_URL="postgres://username:password@localhost:5432/chirpy_db?sslmode=disable"

# JWT Secret (generate a secure random string)
JWT_SECRET="your-super-secret-jwt-key"

# Platform (set to "dev" for development)
PLATFORM="dev"

# Polka API Key (for premium features)
POLKA_KEY="your-polka-api-key"
```

You can also have a look at the .env.example file to get started.

### 4. Run the Application

```bash
# Start the server
go run main.go
```

The server will start on `http://localhost:8080`

## API Documentation

### Base URL

```
http://localhost:8080
```

### Authentication

Most endpoints require authentication using JWT tokens. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

### Health Check

#### GET /api/healthz

Check if the API is running.

**Response:**

```
Status: 200 OK
Content-Type: text/plain

OK
```

### User Management

#### POST /api/users

Create a new user account.

**Request Body:**

```json
{
    "email": "user@example.com",
    "password": "securepassword"
}
```

**Response:**

```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "email": "user@example.com",
    "is_chirpy_red": false
}
```

#### POST /api/login

Authenticate a user and get access tokens.

**Request Body:**

```json
{
    "email": "user@example.com",
    "password": "securepassword"
}
```

**Response:**

```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "email": "user@example.com",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "550e8400-e29b-41d4-a716-446655440001",
    "is_chirpy_red": bool
}
```

#### PUT /api/users

Update user information (requires authentication).

**Request Body:**

```json
{
    "email": "newemail@example.com",
    "password": "newpassword"
}
```

**Response:**

```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z",
    "email": "newemail@example.com",
    "is_chirpy_red": bool
}
```

### Token Management

#### POST /api/refresh

Refresh an expired JWT token using a refresh token.

**Request Body:**

```json
{
    "refresh_token": "550e8400-e29b-41d4-a716-446655440001"
}
```

**Response:**

```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### POST /api/revoke

Revoke a refresh token (requires authentication).

**Request Body:**

```json
{
    "refresh_token": "550e8400-e29b-41d4-a716-446655440001"
}
```

**Response:**

```
Status: 204 No Content
```

### Chirps (Posts)

#### POST /api/chirps

Create a new chirp (requires authentication).

**Request Body:**

```json
{
    "body": "This is my first chirp!"
}
```

**Response:**

```json
{
    "id": "550e8400-e29b-41d4-a716-446655440002",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "body": "This is my first chirp!",
    "user_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

#### GET /api/chirps

Get all chirps with optional filtering and sorting.

**Query Parameters:**

-   `author_id` (optional): Filter chirps by user ID
-   `sort` (optional): Sort order - `"desc"` for newest first

**Example:**

```
GET /api/chirps?author_id=550e8400-e29b-41d4-a716-446655440000&sort=desc
```

**Response:**

```json
[
    {
        "id": "550e8400-e29b-41d4-a716-446655440002",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z",
        "body": "This is my first chirp!",
        "user_id": "550e8400-e29b-41d4-a716-446655440000"
    }
]
```

#### GET /api/chirps/{id}

Get a specific chirp by ID.

**Response:**

```json
{
    "id": "550e8400-e29b-41d4-a716-446655440002",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "body": "This is my first chirp!",
    "user_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

#### DELETE /api/chirps/{id}

Delete a chirp (requires authentication and ownership).

**Response:**

```
Status: 200 OK
```

### Premium Features

#### POST /api/polka/webhooks

Webhook endpoint for Polka integration (premium user upgrades).

**Headers:**

```
Authorization: ApiKey <polka-api-key>
```

**Request Body:**

```json
{
    "event": "user.upgraded",
    "data": {
        "user_id": "550e8400-e29b-41d4-a716-446655440000"
    }
}
```

**Response:**

```
Status: 204 No Content
```

### Admin Endpoints

#### POST /admin/reset

Reset all users and metrics (development only).

**Response:**

```
Status: 200 OK
```

#### GET /admin/metrics

Get server metrics and statistics.

**Response:**

```json
{
    "hits": 42
}
```

### Static Files

#### GET /app/\*

Serve static files from the `app/` directory.

**Example:**

```
GET /app/index.html
```

## Error Responses

All endpoints return consistent error responses:

```json
{
    "error": "Error description"
}
```

Common HTTP status codes:

-   `200 OK` - Success
-   `201 Created` - Resource created
-   `204 No Content` - Success with no response body
-   `400 Bad Request` - Invalid request data
-   `401 Unauthorized` - Authentication required
-   `403 Forbidden` - Access denied
-   `404 Not Found` - Resource not found
-   `500 Internal Server Error` - Server error

## Development

### Project Structure

```
chirpy/
├── app/              # Static files and frontend
├── internal/         # Internal packages
│   ├── auth/         # Authentication logic
│   ├── database/     # Database operations
│   ├── handlers/     # HTTP request handlers
│   ├── types/        # Type definitions
│   └── utils/        # Utility functions
├── sql/              # Database migrations
├── main.go           # Application entry point
└── go.mod            # Go module dependencies
```

### Building

```bash
# Build for your current platform
go build -o chirpy

# Build for specific platforms
GOOS=linux GOARCH=amd64 go build -o chirpy-linux
GOOS=darwin GOARCH=amd64 go build -o chirpy-macos
GOOS=windows GOARCH=amd64 go build -o chirpy-windows.exe
```

### Running with Air (live reload)

Use Air for hot-reloading during development. A preconfigured `.air.toml` is included.

```bash
# Install Air (first time only)
go install github.com/air-verse/air@latest

# From the chirpy/ directory, start the dev server with live reload
air

# Or explicitly specify the config file
air -c .air.toml
```

Notes:

-   Ensure your `.env` is present; variables are loaded via `godotenv` on start.
-   Run Air from the `chirpy/` directory so file paths match the config.

### Database Migrations

The application uses SQLC for type-safe database operations. Database schema is defined in the `sql/` directory.

## License

This project is part of the Boot.dev curriculum and is for educational purposes.
