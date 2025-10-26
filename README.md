# Backend Journaling - Authentication System

A secure authentication system built with Go, PostgreSQL, and JWT tokens.

## Features

- User registration with email OTP verification
- Login with JWT access tokens and refresh tokens
- Password reset with OTP
- Profile management
- Role-based access control
- Rate limiting
- Secure password hashing (Argon2id)
- Email delivery via SMTP
- Audit logging

## Prerequisites

- Go 1.21+
- PostgreSQL 14+
- SMTP server (Gmail, SendGrid, etc.)

## Setup

1. **Clone and install dependencies**
   ```bash
   go mod download
   ```

2. **Generate JWT keys**
   ```bash
   chmod +x scripts/generate-keys.sh
   ./scripts/generate-keys.sh
   ```

3. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Create database**
   ```sql
   CREATE DATABASE journaling_auth;
   ```

5. **Run the server**
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080` (or your configured port).

## 🧪 Testing with Postman

Import the Postman collection for easy API testing:

1. **Import Collection**: `postman_collection.json`
2. **Import Environment**: `postman_environment_local.json` or `postman_environment_production.json`
3. Select the environment from dropdown (top right)
4. Start testing! Tokens are auto-saved after login.

📖 See [POSTMAN_GUIDE.md](POSTMAN_GUIDE.md) for detailed instructions.

## API Endpoints

### Authentication
- `POST /auth/register` - Register new user
- `POST /auth/verify-otp` - Verify OTP after registration
- `POST /auth/login` - Login user
- `POST /auth/refresh` - Refresh access token
- `POST /auth/logout` - Logout user
- `POST /auth/forgot-password` - Request password reset
- `POST /auth/reset-password` - Reset password with OTP
- `POST /auth/request-otp` - Request new OTP

### Profile (Authenticated)
- `GET /profile` - Get user profile
- `PUT /profile` - Update user profile
- `PUT /profile/change-password` - Change password

### Admin (Authenticated + Admin Role)
- `GET /users/:id` - Get user by ID

### Health
- `GET /health` - Health check

## Project Structure

```
.
├── config/              # Configuration management
├── internal/
│   ├── database/        # Database setup and migrations
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # Middleware components
│   ├── models/          # Data models
│   ├── repository/      # Database repositories
│   └── service/         # Business logic
├── pkg/
│   ├── email/           # Email sender
│   ├── jwt/             # JWT token management
│   ├── otp/             # OTP generation and verification
│   ├── password/        # Password hashing
│   └── token/           # Refresh token utilities
├── scripts/             # Utility scripts
├── main.go              # Application entry point
└── README.md
```

## Security Features

- Argon2id password hashing
- RS256 JWT signing
- OTP with SHA256 hashing + pepper
- Rate limiting
- CORS protection
- Security headers
- Audit logging
- Constant-time comparisons

## Environment Variables

See `.env.example` for all configuration options.

## License

MIT
