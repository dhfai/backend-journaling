# Quick Start Guide - Backend Journaling API

> Get started quickly with the Backend Journaling API

## Prerequisites

- MongoDB running on `localhost:27017`
- PostgreSQL running on `localhost:5432`
- Go 1.21 or higher

## Setup

1. **Clone and Install Dependencies**
```bash
cd backend-journaling
go mod download
```

2. **Configure Environment**
```bash
cp .env.example .env
# Edit .env with your settings
```

3. **Start MongoDB**
```bash
# Using Docker
docker run -d -p 27017:27017 --name mongodb mongo:latest

# Or using local MongoDB
mongod --dbpath /path/to/data
```

4. **Run the Server**
```bash
go run main.go
# or
make run
```

Server will start on `http://localhost:8080`

## Authentication Flow

### 1. Register User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!"
  }'
```

### 2. Verify OTP (Check Email)
```bash
curl -X POST http://localhost:8080/api/v1/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "otp": "123456"
  }'
```

### 3. Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!"
  }'
```

Response:
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "user": {...}
}
```

**Save the access_token!** Use it for all subsequent requests.

---

## Quick Examples

### 📝 Create a Note

```bash
TOKEN="your_access_token_here"

curl -X POST http://localhost:8080/api/v1/notes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "My First Note",
    "tags": ["personal"]
  }'
```

### ➕ Add Content Block

```bash
NOTE_ID="your_note_id_here"

curl -X POST http://localhost:8080/api/v1/notes/$NOTE_ID/blocks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "type": "paragraph",
    "content_md": "This is my **first** paragraph!"
  }'
```

### ✅ Create Todo

```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "Buy groceries",
    "priority": "high"
  }'
```

### 📊 Create Task

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "Build Feature",
    "description_md": "## TODO\n- Design\n- Code\n- Test",
    "priority": "high",
    "tags": ["project"]
  }'
```

---

## API Endpoints Summary

### 🔐 Authentication
- `POST /auth/register` - Register new user
- `POST /auth/verify-otp` - Verify email OTP
- `POST /auth/login` - Login user
- `POST /auth/refresh` - Refresh access token

### 📝 Notes
- `POST /notes` - Create note
- `GET /notes` - Get all notes
- `GET /notes/:id` - Get specific note
- `PATCH /notes/:id` - Update note
- `DELETE /notes/:id` - Delete note
- `POST /notes/:id/blocks` - Add block
- `PATCH /notes/:id/blocks/:blockId` - Update block
- `DELETE /notes/:id/blocks/:blockId` - Delete block
- `PATCH /notes/:id/blocks/order` - Reorder blocks

### ✅ Todos
- `POST /todos` - Create todo
- `GET /todos` - Get all todos
- `PATCH /todos/:id` - Update todo
- `DELETE /todos/:id` - Delete todo

### 📊 Tasks
- `POST /tasks` - Create task
- `GET /tasks` - Get all tasks
- `GET /tasks/:id` - Get specific task
- `PATCH /tasks/:id` - Update task
- `DELETE /tasks/:id` - Delete task

---

## Using Postman

1. **Import Collection**
   - Open Postman
   - Import `api_test/postman_collection.json`
   - Import `api_test/postman_environment_local.json`

2. **Set Environment**
   - Select "Local" environment
   - Variables will auto-populate from responses

3. **Test Flow**
   - Run "Register" request
   - Check email for OTP
   - Run "Verify OTP"
   - Run "Login" (token auto-saved)
   - Try other endpoints!

---

## Common Issues

### MongoDB Connection Error
```
Failed to connect to MongoDB
```
**Solution:** Make sure MongoDB is running on port 27017

### Token Expired
```
{"error": "Token has expired"}
```
**Solution:** Use the refresh token endpoint or login again

### Note/Todo/Task Not Found
```
{"error": "Note not found"}
```
**Solution:** Make sure you're using the correct ID and you own the resource

---

## Testing with Examples

All example JSON files are in the `examples/` directory:

**Notes:**
- `create-note.json`
- `update-note.json`
- `add-block-paragraph.json`
- `add-block-heading.json`
- `add-block-todo.json`
- `reorder-blocks.json`

**Todos:**
- `create-todo.json`
- `update-todo.json`

**Tasks:**
- `create-task.json`
- `update-task.json`

**Usage:**
```bash
curl -X POST http://localhost:8080/api/v1/notes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d @examples/create-note.json
```

---

## Project Structure

```
backend-journaling/
├── main.go                 # Application entry point
├── config/                 # Configuration
├── internal/
│   ├── database/          # Database connections
│   ├── handlers/          # HTTP handlers
│   ├── middleware/        # Middleware (auth, CORS, etc)
│   ├── models/           # Data models
│   ├── repository/       # Data access layer
│   └── service/          # Business logic
├── pkg/                  # Shared packages
├── docs/                 # Documentation
├── examples/             # Request examples
└── api_test/            # Postman collection
```

---

## Environment Variables

Key variables you need to set:

```env
# MongoDB
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=journaling

# PostgreSQL (for auth)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=journaling_auth

# JWT Keys
JWT_PRIVATE_KEY_PATH=./keys/jwt_private.pem
JWT_PUBLIC_KEY_PATH=./keys/jwt_public.pem

# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
```

See `.env.example` for complete list.

---

## Next Steps

1. ✅ Complete authentication flow
2. ✅ Create your first note
3. ✅ Add blocks to notes
4. ✅ Create todos and tasks
5. 📖 Read full [API Documentation](./API_DOCUMENTATION.md)
6. 🧪 Explore with Postman collection
7. 🔧 Build your frontend application!

---

## Getting Help

- **Full API Docs:** [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)
- **Reorder Blocks Guide:** [REORDER_BLOCKS_GUIDE.md](./REORDER_BLOCKS_GUIDE.md)
- **Issues:** Report on GitHub
- **Examples:** Check `examples/` folder

---

**Happy Coding! 🚀**
