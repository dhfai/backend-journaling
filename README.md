# Backend Journaling API

> Modern journaling backend API with Notes, Tasks, and Todos management. Built with **Go**, using **PostgreSQL** for authentication and **MongoDB** for journaling data.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![MongoDB](https://img.shields.io/badge/MongoDB-4.4+-47A248?style=flat&logo=mongodb&logoColor=white)](https://www.mongodb.com)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13+-316192?style=flat&logo=postgresql&logoColor=white)](https://www.postgresql.org)

## 📚 Documentation

- **[Quick Start Guide](./docs/QUICK_START.md)** - Get started in 5 minutes
- **[Complete API Documentation](./docs/API_DOCUMENTATION.md)** - Full endpoint reference
- **[Reorder Blocks Guide](./docs/REORDER_BLOCKS_GUIDE.md)** - Detailed guide for block reordering

## 🚀 Features

### ✅ Authentication & Authorization
- User registration dengan email verification (OTP)
- JWT-based authentication (Access & Refresh Tokens)
- Password reset via OTP
- Role-based access control (Admin/User)
- Secure password hashing dengan Argon2

### 📝 Notes Management
- Buat, edit, hapus notes
- Support untuk multiple block types:
  - **Paragraph** - Rich text dengan Markdown
  - **Heading** - Section headers
  - **Todo List** - Checklist dalam note
- Reorder blocks dengan drag & drop support
- Tagging dan pinning notes
- Per-user data isolation

### ✅ Todos
- Simple daily todo list management
- Priority levels (low, medium, high)
- Due dates support
- Mark as done/undone
- Quick task tracking

### 📊 Tasks
- Complex task management dengan status tracking
- Status workflow: todo → in_progress → done
- Priority levels
- Deadline support
- Tags untuk kategorisasi
- Markdown description support

## 🏗️ Tech Stack

- **Language**: Go 1.21+
- **Databases**:
  - PostgreSQL (Auth, Users, Profiles)
  - MongoDB (Notes, Tasks, Todos)
- **Router**: Chi v5
- **JWT**: golang-jwt/jwt v5
- **Password**: Argon2 hashing
- **Email**: SMTP support

## 📁 Project Structure

```
.
├── config/           # Configuration management
├── internal/
│   ├── database/     # Database connections (PostgreSQL & MongoDB)
│   ├── handlers/     # HTTP request handlers
│   ├── middleware/   # Auth, CORS, Rate limiting, Security
│   ├── models/       # Data models (PostgreSQL & MongoDB)
│   ├── repository/   # Data access layer
│   └── service/      # Business logic
├── pkg/              # Reusable packages
│   ├── email/        # Email sending
│   ├── jwt/          # JWT token management
│   ├── otp/          # OTP generation & verification
│   ├── password/     # Password hashing
│   └── token/        # Token generation
├── examples/         # API request examples
└── keys/             # RSA keys for JWT
```

## 🔧 Installation

### Prerequisites

- Go 1.21+
- PostgreSQL 12+
- MongoDB 4.4+
- SMTP server access (Gmail, etc.)

### Setup

1. Clone repository:
```bash
git clone <repository-url>
cd backend-journaling
```

2. Install dependencies:
```bash
go mod download
```

3. Generate JWT keys:
```bash
./scripts/generate-keys.sh
```

4. Setup environment variables:
```bash
cp .env.example .env
# Edit .env dengan konfigurasi Anda
```

5. Setup PostgreSQL database:
```bash
createdb journaling_auth
```

6. Setup MongoDB (optional, bisa menggunakan default):
```bash
# MongoDB akan membuat database otomatis
```

7. Run migrations (auto-run on startup):
```bash
go run main.go
```

## 🌍 Environment Variables

```env
# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=journaling_auth
DB_SSLMODE=disable

# MongoDB
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=journaling

# JWT
JWT_PRIVATE_KEY_PATH=./keys/jwt_private.pem
JWT_PUBLIC_KEY_PATH=./keys/jwt_public.pem
JWT_ACCESS_TOKEN_DURATION=15m
JWT_REFRESH_TOKEN_DURATION=168h

# OTP
OTP_PEPPER=your-secret-pepper-change-this
OTP_TTL_MINUTES=5
OTP_MAX_ATTEMPTS=5

# SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM_EMAIL=your-email@gmail.com
SMTP_FROM_NAME=Journaling App

# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
ENVIRONMENT=development
```

## 🚀 Running the Application

### Development
```bash
go run main.go
```

### Production (Build)
```bash
go build -o bin/backend-journaling
./bin/backend-journaling
```

### Using Make
```bash
make run      # Run in development
make build    # Build binary
make test     # Run tests
```

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/auth/register` | Register new user | No |
| POST | `/auth/verify-otp` | Verify OTP code | No |
| POST | `/auth/login` | Login user | No |
| POST | `/auth/refresh` | Refresh access token | No |
| POST | `/auth/logout` | Logout user | No |
| POST | `/auth/forgot-password` | Request password reset | No |
| POST | `/auth/reset-password` | Reset password with OTP | No |

### Profile Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/profile` | Get user profile | Yes |
| POST | `/profile` | Create profile | Yes |
| PUT | `/profile` | Update profile | Yes |
| PUT | `/profile/avatar` | Update avatar | Yes |
| PUT | `/profile/change-password` | Change password | Yes |

### Notes Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/notes` | List all notes | Yes |
| POST | `/notes` | Create new note | Yes |
| GET | `/notes/:id` | Get note details | Yes |
| PATCH | `/notes/:id` | Update note | Yes |
| DELETE | `/notes/:id` | Delete note | Yes |
| POST | `/notes/:id/blocks` | Add block to note | Yes |
| PATCH | `/notes/:id/blocks/:blockId` | Update block | Yes |
| DELETE | `/notes/:id/blocks/:blockId` | Delete block | Yes |
| PATCH | `/notes/:id/blocks/order` | Reorder blocks | Yes |

### Todos Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/todos` | List all todos | Yes |
| POST | `/todos` | Create new todo | Yes |
| PATCH | `/todos/:id` | Update todo | Yes |
| DELETE | `/todos/:id` | Delete todo | Yes |

### Tasks Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/tasks` | List all tasks | Yes |
| POST | `/tasks` | Create new task | Yes |
| GET | `/tasks/:id` | Get task details | Yes |
| PATCH | `/tasks/:id` | Update task | Yes |
| DELETE | `/tasks/:id` | Delete task | Yes |

## 📝 API Examples

### Create Note
```bash
curl -X POST http://localhost:8080/api/v1/notes \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d @examples/create-note.json
```

### Add Block to Note
```bash
curl -X POST http://localhost:8080/api/v1/notes/{note_id}/blocks \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d @examples/add-block-paragraph.json
```

### Create Todo
```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d @examples/create-todo.json
```

### Create Task
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d @examples/create-task.json
```

## 🏛️ Architecture

### Database Schema

**PostgreSQL (Auth)**:
- users
- otps
- refresh_tokens
- auth_events
- profiles

**MongoDB (Journaling)**:
- notes (dengan embedded blocks)
- todos
- tasks

### Key Integration Points

1. **User Authentication**: PostgreSQL JWT → MongoDB user_id (as string)
2. **Data Isolation**: Semua query MongoDB filtered by user_id dari JWT claims
3. **Clean Architecture**: Repository → Service → Handler pattern

## 🔒 Security Features

- JWT dengan RSA-256 signing
- Argon2 password hashing
- OTP-based email verification
- Rate limiting
- CORS protection
- Security headers
- SQL injection prevention (prepared statements)
- NoSQL injection prevention (bson queries)

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/service/...
```

## 📦 Building

```bash
# Build for current platform
go build -o bin/backend-journaling

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o bin/backend-journaling-linux

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o bin/backend-journaling.exe
```

## 🐳 Docker (Coming Soon)

```bash
docker build -t backend-journaling .
docker run -p 8080:8080 backend-journaling
```

## 🤝 Contributing

1. Fork repository
2. Create feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open Pull Request

## 📄 License

This project is licensed under the MIT License.

## 👤 Author

**dhfai**

## 🙏 Acknowledgments

- Chi Router
- MongoDB Go Driver
- PostgreSQL
- JWT-Go
- Argon2

---

**Happy Coding! 🚀**
