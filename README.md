# Meet Book API

A RESTful API for managing meeting room bookings, built with Go, Gin, and PostgreSQL.

## Features

- 🔐 JWT Authentication
- 📅 Meeting Room Booking System
- 🗄️ PostgreSQL Database
- 📚 Auto-generated API Documentation with Swagger
- 🐳 Docker Support
- 🔄 Database Migrations
- 🧪 Testing Setup

## Tech Stack

- **Framework**: Gin
- **Database**: PostgreSQL
- **ORM**: GORM
- **API Docs**: Swagger
- **Containerization**: Docker

## Prerequisites

- Go 1.23+
- PostgreSQL 13+
- Docker (optional)
- Make (recommended)

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/meet-book-api.git
cd meet-book-api
```

### 2. Setup Environment

```bash
cp .env.example .env
# Edit .env with your configuration
```

### 3. Install Dependencies

```bash
go mod download
```

### 4. Run Migrations

```bash
make migrate
```

### 5. Start the Server

```bash
make run
```

The API will be available at `http://localhost:8080`

## API Documentation

After starting the server, access the interactive API documentation at:
- Swagger UI: http://localhost:8080/swagger/index.html

## Available Commands

| Command         | Description                                      |
|-----------------|--------------------------------------------------|
| `make run`      | Start the development server                     |
| `make build`    | Build the application                            |
| `make test`     | Run tests                                        |
| `make migrate`  | Run database migrations                          |
| `make seed`     | Seed the database with sample data               |
| `make clean`    | Reset the database (drops all tables)            |
| `make docs`     | Generate API documentation                       |

## Environment Variables

| Variable               | Description                          | Default                          |
|------------------------|--------------------------------------|----------------------------------|
| `DATABASE_DIRECT_URL`  | PostgreSQL connection string         | `postgres://user:pass@host/db`   |
| `JWT_SECRET`           | Secret key for JWT tokens            | `your-secret-key`                |
| `MASTER_PASSWORD`      | Master password for admin creation    | `secret-master`                  |
| `PORT`                 | Server port                          | `8080`                           |

## Running with Docker

### Using Docker Compose (Recommended)

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f app
```

### Manual Docker Build

```bash
# Build the image
docker build -t meet-book-api .

# Run the container
docker run -p 8080:8080 --env-file .env meet-book-api
```

## Creating Admin Users

1. Register a new user through the `/auth/register` endpoint
2. Include the `master_password` field in your request with the value from your `.env` file
3. The user will be created with admin privileges

## Database Schema

The database schema includes the following tables:
- `users` - User accounts and authentication
- `rooms` - Meeting rooms
- `bookings` - Room reservations

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.