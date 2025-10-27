# Notes API - Reorder Blocks Guide

## Overview
Endpoint untuk mengatur ulang urutan blocks dalam sebuah note.

## Endpoint
```
PATCH /api/v1/notes/:id/blocks/order
```

## Authentication
Memerlukan Bearer token di header `Authorization`

## Request Body
```json
{
  "order": ["block_id_1", "block_id_2", "block_id_3"]
}
```

### Field Description:
- `order` (array of strings, required): Array berisi ID dari semua blocks yang ingin diurutkan ulang

## Important Notes:
1. **Harus menyertakan SEMUA block IDs** yang ada di note tersebut
2. Jika ada block yang hilang atau ID tidak valid, akan error
3. Order dimulai dari index 0 (block pertama)
4. Block ID bisa didapatkan saat membuat block atau dari GET note endpoint

## Example Workflow:

### 1. Buat Note
```bash
POST /api/v1/notes
{
  "title": "My Note",
  "tags": ["example"]
}
```

Response:
```json
{
  "id": "68fff6e8bafd4f3b24cf67a0",
  "title": "My Note",
  "blocks": [],
  ...
}
```

### 2. Tambah Beberapa Blocks
```bash
POST /api/v1/notes/68fff6e8bafd4f3b24cf67a0/blocks
{
  "type": "heading",
  "content_md": "# Title"
}
```

Response block 1:
```json
{
  "id": "b1-uuid-here",
  "type": "heading",
  "order": 0,
  "content_md": "# Title"
}
```

Tambah block 2:
```bash
POST /api/v1/notes/68fff6e8bafd4f3b24cf67a0/blocks
{
  "type": "paragraph",
  "content_md": "First paragraph"
}
```

Response block 2:
```json
{
  "id": "b2-uuid-here",
  "type": "paragraph",
  "order": 1,
  "content_md": "First paragraph"
}
```

Tambah block 3:
```bash
POST /api/v1/notes/68fff6e8bafd4f3b24cf67a0/blocks
{
  "type": "paragraph",
  "content_md": "Second paragraph"
}
```

Response block 3:
```json
{
  "id": "b3-uuid-here",
  "type": "paragraph",
  "order": 2,
  "content_md": "Second paragraph"
}
```

### 3. Get Note untuk Melihat Current Blocks
```bash
GET /api/v1/notes/68fff6e8bafd4f3b24cf67a0
```

Response:
```json
{
  "id": "68fff6e8bafd4f3b24cf67a0",
  "title": "My Note",
  "blocks": [
    {
      "id": "b1-uuid-here",
      "type": "heading",
      "order": 0,
      "content_md": "# Title"
    },
    {
      "id": "b2-uuid-here",
      "type": "paragraph",
      "order": 1,
      "content_md": "First paragraph"
    },
    {
      "id": "b3-uuid-here",
      "type": "paragraph",
      "order": 2,
      "content_md": "Second paragraph"
    }
  ],
  ...
}
```

### 4. Reorder Blocks
Misalnya ingin memindahkan paragraph kedua ke atas (sebelum paragraph pertama):

```bash
PATCH /api/v1/notes/68fff6e8bafd4f3b24cf67a0/blocks/order
{
  "order": [
    "b1-uuid-here",
    "b3-uuid-here",
    "b2-uuid-here"
  ]
}
```

Response:
```json
{
  "message": "Blocks reordered"
}
```

### 5. Verify dengan GET Note
```bash
GET /api/v1/notes/68fff6e8bafd4f3b24cf67a0
```

Response (blocks sekarang terurut ulang):
```json
{
  "id": "68fff6e8bafd4f3b24cf67a0",
  "title": "My Note",
  "blocks": [
    {
      "id": "b1-uuid-here",
      "type": "heading",
      "order": 0,
      "content_md": "# Title"
    },
    {
      "id": "b3-uuid-here",
      "type": "paragraph",
      "order": 1,
      "content_md": "Second paragraph"
    },
    {
      "id": "b2-uuid-here",
      "type": "paragraph",
      "order": 2,
      "content_md": "First paragraph"
    }
  ],
  ...
}
```

## Error Responses

### 400 Bad Request
```json
{
  "error": "Order is required"
}
```
Terjadi ketika array `order` kosong atau tidak ada.

### 400 Bad Request
```json
{
  "error": "invalid block order"
}
```
Terjadi ketika:
- Jumlah block IDs dalam array tidak sama dengan jumlah blocks di note
- Ada block ID yang tidak valid/tidak ditemukan
- Ada block ID yang duplikat

### 404 Not Found
```json
{
  "error": "Note not found"
}
```
Note dengan ID tersebut tidak ditemukan atau bukan milik user yang login.

### 500 Internal Server Error
```json
{
  "error": "Failed to reorder blocks"
}
```
Terjadi kesalahan server saat menyimpan perubahan.

## Tips
1. Selalu ambil data note terbaru dengan GET sebelum melakukan reorder
2. Copy semua block IDs dari response GET note
3. Susun ulang array sesuai urutan yang diinginkan
4. Pastikan tidak ada block ID yang hilang atau salah
5. Gunakan UUID yang lengkap, bukan hanya sebagian

## Common Mistakes
❌ **SALAH** - Tidak menyertakan semua blocks:
```json
{
  "order": ["b1-uuid-here", "b2-uuid-here"]
}
```
Jika note punya 3 blocks, harus semua disertakan.

❌ **SALAH** - Block ID tidak valid:
```json
{
  "order": ["invalid-id", "b2-uuid-here", "b3-uuid-here"]
}
```

✅ **BENAR** - Semua block IDs valid dan lengkap:
```json
{
  "order": ["b1-uuid-here", "b3-uuid-here", "b2-uuid-here"]
}
```
