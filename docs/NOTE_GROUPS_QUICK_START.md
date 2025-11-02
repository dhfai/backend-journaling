# Note Groups - Quick Start Guide

## Fitur Baru: Note Groups Management

Fitur Note Groups memungkinkan Anda untuk mengorganisir notes ke dalam grup/kategori untuk manajemen yang lebih baik.

## 🎯 Fitur yang Tersedia

1. ✅ **Buat Grup** - Membuat grup baru untuk notes
2. ✅ **Tambah Notes ke Grup** - Menambahkan notes ke dalam grup
3. ✅ **Keluarkan Notes dari Grup** - Mengeluarkan notes dari grup
4. ✅ **Ganti Kepemilikan Grup** - Memindahkan notes antar grup
5. ✅ **Arsip Grup** - Mengarsipkan grup yang tidak aktif
6. ✅ **Delete Grup** - Menghapus grup (notes tidak ikut terhapus)
7. ✅ **Pin/Star Grup** - Menandai grup sebagai favorit

## 📁 Files yang Ditambahkan

### Core Files
- `internal/models/mongo.go` - Updated dengan model `NoteGroup` dan field `group_id` di `Note`
- `internal/repository/note_group.go` - Repository layer untuk note groups
- `internal/service/note_group.go` - Service layer dengan business logic
- `internal/handlers/note_group.go` - HTTP handlers untuk API endpoints

### Documentation & Examples
- `docs/NOTE_GROUPS_API.md` - Dokumentasi lengkap API
- `examples/create-note-group.json` - Contoh create group
- `examples/update-note-group.json` - Contoh update group
- `examples/pin-note-group.json` - Contoh pin group
- `examples/archive-note-group.json` - Contoh archive group
- `examples/add-note-to-group.json` - Contoh add note to group
- `examples/move-notes-to-group.json` - Contoh move notes

## 🚀 Cara Penggunaan

### 1. Membuat Grup Baru

```bash
curl -X POST http://localhost:8080/api/v1/note-groups \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Work Projects",
    "description": "All notes related to work",
    "color": "#3B82F6",
    "icon": "briefcase"
  }'
```

### 2. Melihat Semua Grup

```bash
# Semua grup
curl -X GET http://localhost:8080/api/v1/note-groups \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Filter grup yang di-pin
curl -X GET "http://localhost:8080/api/v1/note-groups?is_pinned=true" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Filter grup yang diarsipkan
curl -X GET "http://localhost:8080/api/v1/note-groups?is_archived=true" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 3. Menambahkan Note ke Grup

```bash
curl -X POST http://localhost:8080/api/v1/note-groups/{GROUP_ID}/notes \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "note_id": "507f1f77bcf86cd799439011"
  }'
```

### 4. Memindahkan Multiple Notes ke Grup

```bash
curl -X POST http://localhost:8080/api/v1/note-groups/{GROUP_ID}/move-notes \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "note_ids": [
      "507f1f77bcf86cd799439011",
      "507f1f77bcf86cd799439012",
      "507f1f77bcf86cd799439013"
    ]
  }'
```

### 5. Mengeluarkan Note dari Grup

```bash
curl -X DELETE http://localhost:8080/api/v1/note-groups/notes/{NOTE_ID} \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 6. Pin/Unpin Grup

```bash
# Pin grup
curl -X PATCH http://localhost:8080/api/v1/note-groups/{GROUP_ID}/pin \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"is_pinned": true}'

# Unpin grup
curl -X PATCH http://localhost:8080/api/v1/note-groups/{GROUP_ID}/pin \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"is_pinned": false}'
```

### 7. Archive/Unarchive Grup

```bash
# Archive grup
curl -X PATCH http://localhost:8080/api/v1/note-groups/{GROUP_ID}/archive \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"is_archived": true}'

# Unarchive grup
curl -X PATCH http://localhost:8080/api/v1/note-groups/{GROUP_ID}/archive \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"is_archived": false}'
```

### 8. Menghapus Grup

```bash
curl -X DELETE http://localhost:8080/api/v1/note-groups/{GROUP_ID} \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Catatan:**
- Menghapus grup TIDAK akan menghapus notes di dalamnya
- Notes hanya akan kehilangan referensi grup (group_id di-set null)
- Notes masih bisa diakses dan dapat di-group kembali ke grup lain

### 9. Melihat Notes dalam Grup

```bash
curl -X GET http://localhost:8080/api/v1/note-groups/{GROUP_ID}/notes \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 📊 API Endpoints Summary

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/note-groups` | Buat grup baru |
| GET | `/api/v1/note-groups` | List semua grup |
| GET | `/api/v1/note-groups/{id}` | Detail satu grup |
| PATCH | `/api/v1/note-groups/{id}` | Update grup |
| DELETE | `/api/v1/note-groups/{id}` | Hapus grup |
| PATCH | `/api/v1/note-groups/{id}/pin` | Pin/unpin grup |
| PATCH | `/api/v1/note-groups/{id}/archive` | Archive/unarchive grup |
| POST | `/api/v1/note-groups/{id}/notes` | Tambah note ke grup |
| DELETE | `/api/v1/note-groups/notes/{noteId}` | Keluarkan note dari grup |
| POST | `/api/v1/note-groups/{id}/move-notes` | Pindahkan multiple notes ke grup |
| GET | `/api/v1/note-groups/{id}/notes` | List notes dalam grup |

## 🔒 Security Features

- **Authentication Required**: Semua endpoints memerlukan valid JWT token
- **User Isolation**: User hanya bisa akses grup mereka sendiri
- **Safe Deletion**: Delete grup tidak menghapus notes, hanya menghapus grup dan reset group_id di notes

## 💾 Database Structure

### NoteGroup Collection
```javascript
{
  _id: ObjectId,
  user_id: String,
  name: String,
  description: String?,
  color: String?,
  icon: String?,
  is_pinned: Boolean,
  is_archived: Boolean,
  notes_count: Integer,
  created_at: DateTime,
  updated_at: DateTime
}
```

### Note Model (Updated)
```javascript
{
  _id: ObjectId,
  user_id: String,
  group_id: ObjectId?,  // NEW FIELD - reference ke NoteGroup
  title: String,
  blocks: Array,
  tags: Array,
  is_pinned: Boolean,
  created_at: DateTime,
  updated_at: DateTime
}
```

## 🎨 Best Practices

1. **Group Naming**: Gunakan nama yang deskriptif dan mudah diingat
2. **Color Coding**: Gunakan warna hex code untuk visual grouping
3. **Icons**: Gunakan icon names yang sesuai untuk identifikasi cepat
4. **Archiving**: Archive grup yang sudah tidak aktif daripada menghapusnya
5. **Pinning**: Pin grup yang sering diakses untuk akses cepat
6. **Organization**: Pisahkan notes berdasarkan context (work, personal, projects, dll)

## 📖 Documentation

Untuk dokumentasi API lengkap, lihat:
- [NOTE_GROUPS_API.md](./NOTE_GROUPS_API.md) - Dokumentasi lengkap dengan semua endpoint, request/response format, dan error handling

## 🧪 Testing

Gunakan contoh JSON di folder `examples/` untuk testing dengan tools seperti:
- Postman
- cURL
- Insomnia
- HTTPie

## 🐛 Troubleshooting

### Error: "Group not found"
- Pastikan GROUP_ID valid dan grup tersebut milik user Anda
- Check apakah grup sudah dihapus

### Error: "Note not found"
- Pastikan NOTE_ID valid dan note tersebut milik user Anda
- Check apakah note sudah dihapus

## 🆕 What's Next?

Fitur yang bisa ditambahkan di masa depan:
- [ ] Shared groups (kolaborasi)
- [ ] Group templates
- [ ] Group statistics
- [ ] Bulk operations
- [ ] Export group dengan semua notes
- [ ] Group permissions

## 💡 Tips

1. Gunakan filter `is_pinned` untuk mendapatkan grup favorit
2. Archive grup lama untuk keep workspace clean
3. Gunakan move-notes untuk reorganisasi cepat
4. Set color & icon untuk visual identification
5. Delete grup tidak menghapus notes, hanya menghapus grouping-nya

---

**Happy Organizing! 📚✨**
