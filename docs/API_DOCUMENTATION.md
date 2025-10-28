# Backend Journaling API Documentation

> Complete API documentation for Notes, Todos, and Tasks management system

**Version:** 2.0.0
**Base URL:** `http://localhost:8080/api/v1`
**Authentication:** Bearer Token (JWT)

---

## Table of Contents

1. [Authentication](#authentication)
2. [Notes API](#notes-api)
3. [Todos API](#todos-api)
4. [Tasks API](#tasks-api)
5. [Error Responses](#error-responses)
6. [Status Codes](#status-codes)

---

## Authentication

All endpoints (except health check) require JWT authentication via Bearer token in the `Authorization` header.

```http
Authorization: Bearer <your_access_token>
```

To get an access token, use the authentication endpoints:
- `POST /auth/register` - Register new user
- `POST /auth/login` - Login and get tokens
- `POST /auth/refresh` - Refresh access token

---

## Notes API

Notes adalah catatan yang bisa berisi beberapa "blocks" (paragraph, heading, todo list, dll).

### Endpoints Overview

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/notes` | Create new note |
| GET | `/notes` | Get all user's notes |
| GET | `/notes/:id` | Get specific note with blocks |
| PATCH | `/notes/:id` | Update note properties |
| DELETE | `/notes/:id` | Delete note |
| POST | `/notes/:id/blocks` | Add block to note |
| PATCH | `/notes/:id/blocks/:blockId` | Update specific block |
| DELETE | `/notes/:id/blocks/:blockId` | Delete specific block |
| PATCH | `/notes/:id/blocks/order` | Reorder all blocks |

---

### 1. Create Note

Create a new empty note with title and optional tags.

**Endpoint:** `POST /notes`

**Headers:**
```http
Content-Type: application/json
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "title": "My Daily Journal",
  "tags": ["daily", "personal"]
}
```

**Request Fields:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| title | string | Yes | Note title |
| tags | array[string] | No | Tags for categorization |

**Response:** `201 Created`
```json
{
  "id": "671e5c9e8bafd4f3b24cf67a0",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "My Daily Journal",
  "blocks": [],
  "tags": ["daily", "personal"],
  "is_pinned": false,
  "created_at": "2025-10-28T10:30:00Z",
  "updated_at": "2025-10-28T10:30:00Z"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid request body or missing title
- `401 Unauthorized` - Missing or invalid token
- `500 Internal Server Error` - Server error

**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/notes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "title": "My Daily Journal",
    "tags": ["daily", "personal"]
  }'
```

---

### 2. Get All Notes

Retrieve all notes for the authenticated user, sorted by updated date (newest first).

**Endpoint:** `GET /notes`

**Headers:**
```http
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
[
  {
    "id": "671e5c9e8bafd4f3b24cf67a0",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "title": "My Daily Journal",
    "blocks": [
      {
        "id": "b1",
        "type": "paragraph",
        "order": 0,
        "content_md": "Today was productive"
      }
    ],
    "tags": ["daily", "personal"],
    "is_pinned": true,
    "created_at": "2025-10-28T10:30:00Z",
    "updated_at": "2025-10-28T15:45:00Z"
  },
  {
    "id": "671e5c9e8bafd4f3b24cf67b1",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "title": "Project Notes",
    "blocks": [],
    "tags": ["work"],
    "is_pinned": false,
    "created_at": "2025-10-27T09:00:00Z",
    "updated_at": "2025-10-27T09:00:00Z"
  }
]
```

**Error Responses:**
- `401 Unauthorized` - Missing or invalid token
- `500 Internal Server Error` - Server error

---

### 3. Get Note by ID

Retrieve a specific note with all its blocks.

**Endpoint:** `GET /notes/:id`

**Headers:**
```http
Authorization: Bearer <token>
```

**URL Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | string | MongoDB ObjectID of the note |

**Response:** `200 OK`
```json
{
  "id": "671e5c9e8bafd4f3b24cf67a0",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "My Daily Journal",
  "blocks": [
    {
      "id": "b1-uuid",
      "type": "heading",
      "order": 0,
      "content_md": "# Today's Highlights"
    },
    {
      "id": "b2-uuid",
      "type": "paragraph",
      "order": 1,
      "content_md": "Completed the **project** on time."
    },
    {
      "id": "b3-uuid",
      "type": "todo",
      "order": 2,
      "items": [
        {
          "id": "t1",
          "text": "Review code",
          "done": true
        },
        {
          "id": "t2",
          "text": "Deploy to production",
          "done": false
        }
      ]
    }
  ],
  "tags": ["daily", "personal"],
  "is_pinned": false,
  "created_at": "2025-10-28T10:30:00Z",
  "updated_at": "2025-10-28T15:45:00Z"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid note ID format
- `401 Unauthorized` - Missing or invalid token
- `404 Not Found` - Note not found or not owned by user
- `500 Internal Server Error` - Server error

---

### 4. Update Note

Update note properties (title, tags, pinned status).

**Endpoint:** `PATCH /notes/:id`

**Headers:**
```http
Content-Type: application/json
Authorization: Bearer <token>
```

**Request Body:** (all fields are optional)
```json
{
  "title": "Updated Title",
  "tags": ["work", "important"],
  "is_pinned": true
}
```

**Request Fields:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| title | string | No | New note title |
| tags | array[string] | No | New tags array |
| is_pinned | boolean | No | Pin/unpin note |

**Response:** `200 OK`
```json
{
  "message": "Note updated"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid request body or no fields to update
- `401 Unauthorized` - Missing or invalid token
- `404 Not Found` - Note not found
- `500 Internal Server Error` - Server error

---

### 5. Delete Note

Delete a note and all its blocks permanently.

**Endpoint:** `DELETE /notes/:id`

**Headers:**
```http
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
{
  "message": "Note deleted"
}
```

**Error Responses:**
- `401 Unauthorized` - Missing or invalid token
- `404 Not Found` - Note not found
- `500 Internal Server Error` - Server error

---

### 6. Add Block to Note

Add a new content block to a note. Blocks can be paragraphs, headings, or todo lists.

**Endpoint:** `POST /notes/:id/blocks`

**Headers:**
```http
Content-Type: application/json
Authorization: Bearer <token>
```

**Request Body for Paragraph/Heading:**
```json
{
  "type": "paragraph",
  "content_md": "This is a **paragraph** with markdown support."
}
```

**Request Body for Todo Block:**
```json
{
  "type": "todo",
  "items": [
    {
      "id": "t1",
      "text": "Complete documentation",
      "done": false
    },
    {
      "id": "t2",
      "text": "Review code",
      "done": true
    }
  ]
}
```

**Request Fields:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| type | string | Yes | Block type: "paragraph", "heading", "todo" |
| content_md | string | Conditional | Required for paragraph/heading. Markdown content |
| items | array[object] | Conditional | Required for todo type. Array of todo items |

**Todo Item Structure:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| id | string | Yes | Unique ID for todo item |
| text | string | Yes | Todo item text |
| done | boolean | Yes | Completion status |

**Response:** `201 Created`
```json
{
  "id": "b1-uuid-generated",
  "type": "paragraph",
  "order": 0,
  "content_md": "This is a **paragraph** with markdown support."
}
```

**Error Responses:**
- `400 Bad Request` - Invalid request or missing type
- `401 Unauthorized` - Missing or invalid token
- `404 Not Found` - Note not found
- `500 Internal Server Error` - Server error

**Block Types:**
- `paragraph` - Regular text block with markdown
- `heading` - Heading block (use # syntax in content_md)
- `todo` - Checklist block with multiple items

---

### 7. Update Block

Update the content of a specific block in a note.

**Endpoint:** `PATCH /notes/:id/blocks/:blockId`

**Headers:**
```http
Content-Type: application/json
Authorization: Bearer <token>
```

**URL Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | string | Note ID |
| blockId | string | Block ID to update |

**Request Body for Paragraph/Heading:**
```json
{
  "content_md": "Updated content with **new** information."
}
```

**Request Body for Todo Block:**
```json
{
  "items": [
    {
      "id": "t1",
      "text": "Updated task",
      "done": true
    },
    {
      "id": "t2",
      "text": "New task",
      "done": false
    }
  ]
}
```

**Response:** `200 OK`
```json
{
  "message": "Block updated"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid request or no fields to update
- `401 Unauthorized` - Missing or invalid token
- `404 Not Found` - Note or block not found
- `500 Internal Server Error` - Server error

---

### 8. Delete Block

Delete a specific block from a note.

**Endpoint:** `DELETE /notes/:id/blocks/:blockId`

**Headers:**
```http
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
{
  "message": "Block deleted"
}
```

**Error Responses:**
- `401 Unauthorized` - Missing or invalid token
- `404 Not Found` - Note or block not found
- `500 Internal Server Error` - Server error

---

### 9. Reorder Blocks

Reorder all blocks in a note by providing the complete list of block IDs in desired order.

**Endpoint:** `PATCH /notes/:id/blocks/order`

**Headers:**
```http
Content-Type: application/json
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "order": ["b3-uuid", "b1-uuid", "b2-uuid"]
}
```

**Request Fields:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| order | array[string] | Yes | Complete array of all block IDs in new order |

**Important Notes:**
- Must include ALL block IDs from the note
- IDs must be valid (existing block IDs)
- Order determines the new sequence (index 0 = first block)

**Response:** `200 OK`
```json
{
  "message": "Blocks reordered"
}
```

**Error Responses:**
- `400 Bad Request` - Empty order array or invalid block order (missing IDs, wrong count, invalid IDs)
- `401 Unauthorized` - Missing or invalid token
- `404 Not Found` - Note not found
- `500 Internal Server Error` - Server error

**Example Workflow:**
1. Get note to see current blocks and their IDs
2. Copy all block IDs
3. Arrange IDs in desired order
4. Send PATCH request with complete order array

See [Reorder Blocks Guide](./REORDER_BLOCKS_GUIDE.md) for detailed examples.

---

## Todos API

Simple todo list management for daily tasks.

### Endpoints Overview

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/todos` | Create new todo |
| GET | `/todos` | Get all user's todos |
| PATCH | `/todos/:id` | Update todo |
| DELETE | `/todos/:id` | Delete todo |

---

### 1. Create Todo

Create a new todo item.

**Endpoint:** `POST /todos`

**Headers:**
```http
Content-Type: application/json
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "title": "Buy groceries",
  "priority": "high",
  "due_date": "2025-10-30T00:00:00Z"
}
```

**Request Fields:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| title | string | Yes | Todo title |
| priority | string | No | Priority level: "low", "medium", "high" (default: "medium") |
| due_date | string | No | ISO 8601 datetime string |

**Response:** `201 Created`
```json
{
  "id": "6720a7c2bafd4f3b24cf67c1",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "Buy groceries",
  "done": false,
  "priority": "high",
  "due_date": "2025-10-30T00:00:00Z",
  "created_at": "2025-10-28T10:30:00Z",
  "updated_at": "2025-10-28T10:30:00Z"
}
```

**Error Responses:**
- `400 Bad Request` - Missing title or invalid request
- `401 Unauthorized` - Missing or invalid token
- `500 Internal Server Error` - Server error

---

### 2. Get All Todos

Retrieve all todos for the authenticated user, sorted by creation date (newest first).

**Endpoint:** `GET /todos`

**Headers:**
```http
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
[
  {
    "id": "6720a7c2bafd4f3b24cf67c1",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "title": "Buy groceries",
    "done": false,
    "priority": "high",
    "due_date": "2025-10-30T00:00:00Z",
    "created_at": "2025-10-28T10:30:00Z",
    "updated_at": "2025-10-28T10:30:00Z"
  },
  {
    "id": "6720a7c2bafd4f3b24cf67c2",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "title": "Finish report",
    "done": true,
    "priority": "medium",
    "due_date": null,
    "created_at": "2025-10-27T09:00:00Z",
    "updated_at": "2025-10-28T14:00:00Z"
  }
]
```

**Error Responses:**
- `401 Unauthorized` - Missing or invalid token
- `500 Internal Server Error` - Server error

---

### 3. Update Todo

Update todo properties (title, done status, priority, due date).

**Endpoint:** `PATCH /todos/:id`

**Headers:**
```http
Content-Type: application/json
Authorization: Bearer <token>
```

**Request Body:** (all fields are optional)
```json
{
  "title": "Updated title",
  "done": true,
  "priority": "low",
  "due_date": "2025-11-01T00:00:00Z"
}
```

**Request Fields:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| title | string | No | New todo title |
| done | boolean | No | Completion status |
| priority | string | No | Priority: "low", "medium", "high" |
| due_date | string | No | ISO 8601 datetime or null |

**Response:** `200 OK`
```json
{
  "message": "Todo updated"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid request or no fields to update
- `401 Unauthorized` - Missing or invalid token
- `404 Not Found` - Todo not found
- `500 Internal Server Error` - Server error

---

### 4. Delete Todo

Delete a todo permanently.

**Endpoint:** `DELETE /todos/:id`

**Headers:**
```http
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
{
  "message": "Todo deleted"
}
```

**Error Responses:**
- `401 Unauthorized` - Missing or invalid token
- `404 Not Found` - Todo not found
- `500 Internal Server Error` - Server error

---

## Tasks API

Complex task management with descriptions, status tracking, and tags.

### Endpoints Overview

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/tasks` | Create new task |
| GET | `/tasks` | Get all user's tasks |
| GET | `/tasks/:id` | Get specific task |
| PATCH | `/tasks/:id` | Update task |
| DELETE | `/tasks/:id` | Delete task |

---

### 1. Create Task

Create a new task with detailed description and metadata.

**Endpoint:** `POST /tasks`

**Headers:**
```http
Content-Type: application/json
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "title": "Build Backend API",
  "description_md": "## Objectives\n- Setup MongoDB\n- Implement endpoints\n- Write tests",
  "priority": "high",
  "deadline": "2025-11-05T00:00:00Z",
  "tags": ["project", "backend"]
}
```

**Request Fields:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| title | string | Yes | Task title |
| description_md | string | No | Markdown description |
| priority | string | No | Priority: "low", "medium", "high" (default: "medium") |
| deadline | string | No | ISO 8601 datetime |
| tags | array[string] | No | Tags for categorization |

**Response:** `201 Created`
```json
{
  "id": "6720b33cbafd4f3b24cf67d1",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "Build Backend API",
  "description_md": "## Objectives\n- Setup MongoDB\n- Implement endpoints\n- Write tests",
  "status": "todo",
  "priority": "high",
  "deadline": "2025-11-05T00:00:00Z",
  "tags": ["project", "backend"],
  "created_at": "2025-10-28T10:30:00Z",
  "updated_at": "2025-10-28T10:30:00Z"
}
```

**Status Values:**
- `todo` - Not started (default)
- `in_progress` - Currently working on it
- `done` - Completed

**Error Responses:**
- `400 Bad Request` - Missing title or invalid request
- `401 Unauthorized` - Missing or invalid token
- `500 Internal Server Error` - Server error

---

### 2. Get All Tasks

Retrieve all tasks for the authenticated user, sorted by creation date (newest first).

**Endpoint:** `GET /tasks`

**Headers:**
```http
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
[
  {
    "id": "6720b33cbafd4f3b24cf67d1",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "title": "Build Backend API",
    "description_md": "## Objectives\n- Setup MongoDB\n- Implement endpoints",
    "status": "in_progress",
    "priority": "high",
    "deadline": "2025-11-05T00:00:00Z",
    "tags": ["project", "backend"],
    "created_at": "2025-10-28T10:30:00Z",
    "updated_at": "2025-10-28T15:00:00Z"
  }
]
```

**Error Responses:**
- `401 Unauthorized` - Missing or invalid token
- `500 Internal Server Error` - Server error

---

### 3. Get Task by ID

Retrieve a specific task with all details.

**Endpoint:** `GET /tasks/:id`

**Headers:**
```http
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
{
  "id": "6720b33cbafd4f3b24cf67d1",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "Build Backend API",
  "description_md": "## Objectives\n- Setup MongoDB\n- Implement endpoints\n- Write tests",
  "status": "in_progress",
  "priority": "high",
  "deadline": "2025-11-05T00:00:00Z",
  "tags": ["project", "backend"],
  "created_at": "2025-10-28T10:30:00Z",
  "updated_at": "2025-10-28T15:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid task ID format
- `401 Unauthorized` - Missing or invalid token
- `404 Not Found` - Task not found
- `500 Internal Server Error` - Server error

---

### 4. Update Task

Update task properties.

**Endpoint:** `PATCH /tasks/:id`

**Headers:**
```http
Content-Type: application/json
Authorization: Bearer <token>
```

**Request Body:** (all fields are optional)
```json
{
  "title": "Updated Task Title",
  "description_md": "## Updated\n- New item",
  "status": "in_progress",
  "priority": "medium",
  "deadline": "2025-11-10T00:00:00Z",
  "tags": ["updated", "project"]
}
```

**Request Fields:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| title | string | No | New task title |
| description_md | string | No | New description |
| status | string | No | Status: "todo", "in_progress", "done" |
| priority | string | No | Priority: "low", "medium", "high" |
| deadline | string | No | ISO 8601 datetime or null |
| tags | array[string] | No | New tags array |

**Response:** `200 OK`
```json
{
  "message": "Task updated"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid request or no fields to update
- `401 Unauthorized` - Missing or invalid token
- `404 Not Found` - Task not found
- `500 Internal Server Error` - Server error

---

### 5. Delete Task

Delete a task permanently.

**Endpoint:** `DELETE /tasks/:id`

**Headers:**
```http
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
{
  "message": "Task deleted"
}
```

**Error Responses:**
- `401 Unauthorized` - Missing or invalid token
- `404 Not Found` - Task not found
- `500 Internal Server Error` - Server error

---

## Error Responses

All error responses follow this format:

```json
{
  "error": "Error message description"
}
```

### Common Error Messages

**Authentication Errors:**
```json
{
  "error": "Missing authorization header"
}
```
```json
{
  "error": "Invalid authorization header format"
}
```
```json
{
  "error": "Token has expired"
}
```
```json
{
  "error": "Invalid token"
}
```

**Validation Errors:**
```json
{
  "error": "Invalid request body"
}
```
```json
{
  "error": "Title is required"
}
```
```json
{
  "error": "No fields to update"
}
```

**Resource Errors:**
```json
{
  "error": "Note not found"
}
```
```json
{
  "error": "Todo not found"
}
```
```json
{
  "error": "Task not found"
}
```
```json
{
  "error": "Note or block not found"
}
```

**Server Errors:**
```json
{
  "error": "Failed to create note"
}
```
```json
{
  "error": "Failed to fetch notes"
}
```

---

## Status Codes

| Code | Description | When It Happens |
|------|-------------|-----------------|
| 200 | OK | Successful GET, PATCH, DELETE request |
| 201 | Created | Successful POST request (resource created) |
| 400 | Bad Request | Invalid request body, missing required fields, validation error |
| 401 | Unauthorized | Missing token, invalid token, expired token |
| 403 | Forbidden | User doesn't have permission (admin endpoints) |
| 404 | Not Found | Resource not found or not owned by user |
| 500 | Internal Server Error | Server-side error, database error |

---

## Data Types & Formats

### Date/Time Format
All timestamps use ISO 8601 format:
```
2025-10-28T10:30:00Z
```

### MongoDB ObjectID Format
24-character hexadecimal string:
```
671e5c9e8bafd4f3b24cf67a0
```

### UUID Format
PostgreSQL user IDs use UUID v4:
```
123e4567-e89b-12d3-a456-426614174000
```

### Priority Values
```
"low" | "medium" | "high"
```

### Task Status Values
```
"todo" | "in_progress" | "done"
```

### Block Type Values
```
"paragraph" | "heading" | "todo"
```

---

## Best Practices

### 1. Always Include Authorization Header
```http
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 2. Use Proper Content-Type for POST/PATCH
```http
Content-Type: application/json
```

### 3. Handle Token Expiration
- Monitor for 401 errors with "Token has expired"
- Use refresh token endpoint to get new access token
- Implement automatic token refresh logic

### 4. Validate Before Sending
- Check required fields are present
- Ensure proper data types
- Validate date formats

### 5. Handle Errors Gracefully
- Parse error messages from response
- Show user-friendly error messages
- Log errors for debugging

### 6. Use Pagination (Future Enhancement)
Currently, all GET endpoints return all records. Consider implementing pagination for large datasets.

---

## Examples with cURL

### Create Note with Blocks
```bash
# 1. Create note
NOTE_ID=$(curl -s -X POST http://localhost:8080/api/v1/notes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title":"My Note","tags":["example"]}' \
  | jq -r '.id')

# 2. Add heading block
curl -X POST http://localhost:8080/api/v1/notes/$NOTE_ID/blocks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"type":"heading","content_md":"# Main Title"}'

# 3. Add paragraph block
curl -X POST http://localhost:8080/api/v1/notes/$NOTE_ID/blocks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"type":"paragraph","content_md":"This is the content."}'

# 4. Get note with all blocks
curl -X GET http://localhost:8080/api/v1/notes/$NOTE_ID \
  -H "Authorization: Bearer $TOKEN"
```

### Create and Complete Todo
```bash
# 1. Create todo
TODO_ID=$(curl -s -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title":"Buy groceries","priority":"high"}' \
  | jq -r '.id')

# 2. Mark as done
curl -X PATCH http://localhost:8080/api/v1/todos/$TODO_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"done":true}'
```

### Create and Track Task
```bash
# 1. Create task
TASK_ID=$(curl -s -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title":"Build API",
    "description_md":"## Steps\n- Design\n- Implement\n- Test",
    "priority":"high",
    "tags":["project"]
  }' | jq -r '.id')

# 2. Update to in_progress
curl -X PATCH http://localhost:8080/api/v1/tasks/$TASK_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"status":"in_progress"}'

# 3. Complete task
curl -X PATCH http://localhost:8080/api/v1/tasks/$TASK_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"status":"done"}'
```

---

## Postman Collection

Import the Postman collection for easy testing:
- **File:** `api_test/postman_collection.json`
- **Environment:** `api_test/postman_environment_local.json`

The collection includes:
- Pre-configured requests for all endpoints
- Auto-saving of tokens and IDs to variables
- Example requests with sample data
- Test scripts for automation

---

## Rate Limiting

Currently, no rate limiting is implemented on these endpoints. Future versions may include rate limiting for production use.

---

## Version History

### v2.0.0 (2025-10-28)
- Added Notes API with blocks support
- Added Todos API
- Added Tasks API
- Integrated MongoDB for content storage
- User isolation via PostgreSQL JWT authentication

### v1.0.0
- Initial release with authentication system

---

## Support & Contact

For issues, questions, or contributions:
- **Repository:** github.com/dhfai/backend-journaling
- **Issues:** Create an issue on GitHub
- **Documentation:** See `/docs` folder for additional guides

---

**Last Updated:** October 28, 2025
