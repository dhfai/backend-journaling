# Note Groups API Documentation

## Overview
Note Groups API memungkinkan pengguna untuk mengorganisir notes mereka ke dalam grup/kategori. Fitur ini mendukung CRUD operations, pin/archive, dan management notes dalam grup.

## Endpoints

### 1. Create Note Group
**Endpoint:** `POST /api/v1/note-groups`

**Authentication:** Required

**Request Body:**
```json
{
  "name": "Work Projects",
  "description": "All notes related to work projects",
  "color": "#3B82F6",
  "icon": "briefcase"
}
```

**Response (201 Created):**
```json
{
  "id": "507f1f77bcf86cd799439011",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Work Projects",
  "description": "All notes related to work projects",
  "color": "#3B82F6",
  "icon": "briefcase",
  "is_pinned": false,
  "is_archived": false,
  "notes_count": 0,
  "created_at": "2025-11-02T10:00:00Z",
  "updated_at": "2025-11-02T10:00:00Z"
}
```

---

### 2. Get All Note Groups
**Endpoint:** `GET /api/v1/note-groups`

**Authentication:** Required

**Query Parameters:**
- `is_pinned` (optional): Filter by pinned status (true/false)
- `is_archived` (optional): Filter by archived status (true/false)

**Example:** `GET /api/v1/note-groups?is_pinned=true`

**Response (200 OK):**
```json
[
  {
    "id": "507f1f77bcf86cd799439011",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "Work Projects",
    "description": "All notes related to work projects",
    "color": "#3B82F6",
    "icon": "briefcase",
    "is_pinned": true,
    "is_archived": false,
    "notes_count": 5,
    "created_at": "2025-11-02T10:00:00Z",
    "updated_at": "2025-11-02T10:00:00Z"
  }
]
```

---

### 3. Get Single Note Group
**Endpoint:** `GET /api/v1/note-groups/{id}`

**Authentication:** Required

**Response (200 OK):**
```json
{
  "id": "507f1f77bcf86cd799439011",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Work Projects",
  "description": "All notes related to work projects",
  "color": "#3B82F6",
  "icon": "briefcase",
  "is_pinned": true,
  "is_archived": false,
  "notes_count": 5,
  "created_at": "2025-11-02T10:00:00Z",
  "updated_at": "2025-11-02T10:00:00Z"
}
```

---

### 4. Update Note Group
**Endpoint:** `PATCH /api/v1/note-groups/{id}`

**Authentication:** Required

**Request Body:**
```json
{
  "name": "Personal Notes",
  "description": "My personal thoughts and ideas",
  "color": "#10B981",
  "icon": "user"
}
```

**Note:** All fields are optional. Only send fields you want to update.

**Response (200 OK):**
```json
{
  "id": "507f1f77bcf86cd799439011",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Personal Notes",
  "description": "My personal thoughts and ideas",
  "color": "#10B981",
  "icon": "user",
  "is_pinned": false,
  "is_archived": false,
  "notes_count": 5,
  "created_at": "2025-11-02T10:00:00Z",
  "updated_at": "2025-11-02T11:00:00Z"
}
```

---

### 5. Pin/Unpin Note Group
**Endpoint:** `PATCH /api/v1/note-groups/{id}/pin`

**Authentication:** Required

**Request Body:**
```json
{
  "is_pinned": true
}
```

**Response (200 OK):**
```json
{
  "message": "Group pin status updated",
  "is_pinned": true
}
```

---

### 6. Archive/Unarchive Note Group
**Endpoint:** `PATCH /api/v1/note-groups/{id}/archive`

**Authentication:** Required

**Request Body:**
```json
{
  "is_archived": true
}
```

**Response (200 OK):**
```json
{
  "message": "Group archive status updated",
  "is_archived": true
}
```

---

### 7. Delete Note Group
**Endpoint:** `DELETE /api/v1/note-groups/{id}`

**Authentication:** Required

**Description:** Deletes a note group and automatically removes the group_id from all notes in the group. Notes themselves are not deleted, they just become ungrouped.

**Request:** No body needed

**Response (200 OK):**
```json
{
  "message": "Group deleted successfully"
}
```

**Note:** When a group is deleted:
- The group record is removed from database
- All notes in the group have their `group_id` field set to `null`
- Notes themselves are NOT deleted
- Notes can be re-grouped later

---

### 8. Add Note to Group
**Endpoint:** `POST /api/v1/note-groups/{id}/notes`

**Authentication:** Required

**Request Body:**
```json
{
  "note_id": "507f1f77bcf86cd799439011"
}
```

**Response (200 OK):**
```json
{
  "message": "Note added to group successfully"
}
```

---

### 9. Remove Note from Group
**Endpoint:** `DELETE /api/v1/note-groups/notes/{noteId}`

**Authentication:** Required

**Response (200 OK):**
```json
{
  "message": "Note removed from group successfully"
}
```

---

### 10. Move Multiple Notes to Group
**Endpoint:** `POST /api/v1/note-groups/{id}/move-notes`

**Authentication:** Required

**Description:** Memindahkan beberapa notes sekaligus ke grup tertentu. Notes akan dipindahkan dari grup lama (jika ada) ke grup baru.

**Request Body:**
```json
{
  "note_ids": [
    "507f1f77bcf86cd799439011",
    "507f1f77bcf86cd799439012",
    "507f1f77bcf86cd799439013"
  ]
}
```

**Response (200 OK):**
```json
{
  "message": "Notes moved to group successfully"
}
```

---

### 11. Get All Notes in Group
**Endpoint:** `GET /api/v1/note-groups/{id}/notes`

**Authentication:** Required

**Response (200 OK):**
```json
[
  {
    "id": "507f1f77bcf86cd799439011",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "group_id": "507f1f77bcf86cd799439020",
    "title": "Meeting Notes",
    "blocks": [],
    "tags": ["work", "meeting"],
    "is_pinned": false,
    "created_at": "2025-11-02T10:00:00Z",
    "updated_at": "2025-11-02T10:00:00Z"
  }
]
```

---

## Error Responses

### 400 Bad Request
```json
{
  "error": "Group name is required"
}
```

### 404 Not Found
```json
{
  "error": "Group not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Failed to create group"
}
```

---

## Features Summary

### ✅ Implemented Features:
1. **Buat Grup** - `POST /api/v1/note-groups`
2. **Tambah Notes ke Grup** - `POST /api/v1/note-groups/{id}/notes`
3. **Keluarkan Notes dari Grup** - `DELETE /api/v1/note-groups/notes/{noteId}`
4. **Ganti Kepemilikan Grup dari Notes** - `POST /api/v1/note-groups/{id}/move-notes`
5. **Arsip Grup** - `PATCH /api/v1/note-groups/{id}/archive`
6. **Delete Grup** - `DELETE /api/v1/note-groups/{id}` (langsung hapus, notes tidak ikut terhapus)
7. **Pin/Star Grup** - `PATCH /api/v1/note-groups/{id}/pin`

### Additional Features:
- Get all groups with filters
- Get single group
- Update group details
- Get all notes in a group
- Automatic notes count tracking
- Sorted by pinned status and updated date

---

## Database Schema

### NoteGroup Collection (MongoDB)
```javascript
{
  _id: ObjectId,
  user_id: String,
  name: String,
  description: String (optional),
  color: String (optional),
  icon: String (optional),
  is_pinned: Boolean,
  is_archived: Boolean,
  notes_count: Integer,
  created_at: DateTime,
  updated_at: DateTime
}
```

### Note Model Update
Added `group_id` field to existing Note model:
```javascript
{
  _id: ObjectId,
  user_id: String,
  group_id: ObjectId (optional),  // NEW FIELD
  title: String,
  blocks: Array,
  tags: Array,
  is_pinned: Boolean,
  created_at: DateTime,
  updated_at: DateTime
}
```

---

## Usage Examples

### Example Flow: Creating and Managing Groups

1. **Create a group:**
   ```bash
   curl -X POST http://localhost:8080/api/v1/note-groups \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"name": "Work Projects", "color": "#3B82F6"}'
   ```

2. **Add notes to group:**
   ```bash
   curl -X POST http://localhost:8080/api/v1/note-groups/GROUP_ID/notes \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"note_id": "NOTE_ID"}'
   ```

3. **Pin the group:**
   ```bash
   curl -X PATCH http://localhost:8080/api/v1/note-groups/GROUP_ID/pin \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"is_pinned": true}'
   ```

4. **Delete group:**
   ```bash
   # Deletes group and removes group_id from all notes
   curl -X DELETE http://localhost:8080/api/v1/note-groups/GROUP_ID \
     -H "Authorization: Bearer YOUR_TOKEN"
   ```

---

## Notes

- Groups akan di-sort berdasarkan `is_pinned` (descending) kemudian `updated_at` (descending)
- Notes count otomatis di-update saat add/remove notes
- Delete group TIDAK menghapus notes, hanya menghapus group dan reset group_id di notes
- Notes yang kehilangan group dapat di-group kembali ke group lain
