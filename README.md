# URL Shortener

A high-performance URL shortening service built with Go, featuring JWT authentication, PostgreSQL database, and Redis caching.

## Features

- **User Authentication**: JWT-based authentication with access and refresh tokens
- **URL Shortening**: Custom algorithm using Base62 encoding with prefix/suffix randomization
- **Deep Link Support**: Mobile deep linking for iOS and Android apps
- **Redis Caching**: Fast URL lookups with 5-minute cache TTL
- **Rate Limiting**: Built-in rate limiting middleware
- **Soft Delete**: User and URL data can be soft deleted and restored

## Tech Stack

- **Language**: Go 1.26
- **Web Framework**: Gin
- **Database**: PostgreSQL (via pgx)
- **Cache**: Redis
- **Authentication**: JWT (golang-jwt)
- **Validation**: go-playground/validator
- **ORM**: SQLC

## Project Structure

```
url_shortener/
├── cmd/api/main.go           # Application entry point
├── internal/
│   ├── app/                  # Application module setup
│   ├── config/               # Configuration management
│   ├── database/             # Database connection and migrations
│   │   └── migrations/      # SQL migrations
│   ├── dto/                  # Data transfer objects
│   ├── handler/             # HTTP handlers
│   ├── middleware/          # HTTP middleware (auth, rate limit, CORS, etc.)
│   ├── repository/          # Data access layer
│   ├── routes/               # Route definitions
│   ├── services/            # Business logic
│   └── utils/               # Utility functions
├── pkg/
│   ├── auth/                 # JWT authentication
│   ├── cache/                # Redis cache service
│   └── logger/               # Logging utilities
└── Makefile                  # Build and development commands
```

## API Endpoints

### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/login` | Login with email/password |
| POST | `/api/v1/auth/logout` | Logout (invalidate tokens) |
| POST | `/api/v1/auth/refresh` | Refresh access token |

### Users
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/users` | Create new user |
| GET | `/api/v1/users` | List all users (paginated) |
| GET | `/api/v1/users/:uuid` | Get user by UUID |
| PUT | `/api/v1/users/:uuid` | Update user |
| DELETE | `/api/v1/users/:uuid` | Soft delete user |
| DELETE | `/api/v1/users/:uuid/hard` | Permanently delete user |
| PUT | `/api/v1/users/:uuid/restore` | Restore soft-deleted user |

### URLs
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/urls` | Create short URL |
| GET | `/api/v1/urls/:shortKey` | Redirect to original URL |

### Public
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/:shortKey` | Redirect shortened URL to target |

## Getting Started

### Prerequisites

- Go 1.26+
- PostgreSQL
- Redis
- Docker (optional)

### Environment Variables

Create a `.env` file:

```env
# Server
SERVER_HOST=localhost
SERVER_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=url_shortener
DB_SSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT
JWT_SECRET=your_secret_key
JWT_ACCESS_EXPIRE=900
JWT_REFRESH_EXPIRE=86400

# URL Settings
RAND_PREFIX_LEN=2
RANDOM_KEY_LEN=4
```

### Running with Docker

```bash
make dev
```

### Running Locally

1. Start PostgreSQL and Redis
2. Run migrations:
```bash
make migrate-up
```
3. Run the server:
```bash
make server
```

## Short URL Algorithm

The system uses a custom Base62 encoding algorithm:

1. Generate a random key (configurable length, default 4 characters)
2. Split into prefix and suffix
3. Encode the database ID using Base62 (0-9, a-z, A-Z)
4. Combine: `prefix + encoded_id + suffix`

This approach:
- Creates short, readable URLs
- Makes URLs hard to guess (random prefix/suffix)
- Allows efficient database lookups by ID

## Available Make Commands

| Command | Description |
|---------|-------------|
| `make server` | Run the API server |
| `make dev` | Start development environment with Docker |
| `make migrate-up` | Run all pending migrations |
| `make migrate-down` | Rollback last migration |
| `make sqlc` | Generate SQLC code |
| `make accessdb` | Access PostgreSQL container |
| `make importdb` | Import database from backup |
| `make exportdb` | Export database to backup |

## License

MIT