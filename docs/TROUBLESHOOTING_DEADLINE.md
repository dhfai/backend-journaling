# 🔧 Troubleshooting: Deadline/Due Date Issues

## ❌ Masalah yang Terjadi

### Error Message
```json
{
    "success": false,
    "error": "error decoding key deadline: cannot decode string into a models.FlexibleTime"
}
```

### Root Cause
MongoDB sudah menyimpan data `deadline` sebagai **STRING** (`"2025-10-31"`), bukan sebagai **DateTime object**. Saat backend mencoba membaca data dari MongoDB, BSON decoder tidak bisa convert string ke FlexibleTime.

---

## 📊 Penjelasan Detail

### 1. **Bagaimana Data Tersimpan di MongoDB?**

Cek data Anda di MongoDB:
```javascript
// Data yang bermasalah (deadline sebagai STRING)
{
  "_id": {"$oid": "69002010e26d267da934b3eb"},
  "title": "dengan Nunu",
  "deadline": "2025-10-31",  // ❌ INI STRING, bukan DateTime!
  "user_id": "ccc0d5fb-6263-43a2-bc74-f3e6f16a062d"
}

// Data yang benar (deadline sebagai DateTime)
{
  "_id": {"$oid": "69002011e26d267da934b3ec"},
  "title": "New Task",
  "deadline": {"$date": "2025-10-31T00:00:00.000Z"},  // ✅ INI DateTime object
  "user_id": "ccc0d5fb-6263-43a2-bc74-f3e6f16a062d"
}
```

### 2. **Kenapa Terjadi?**

Kemungkinan besar task dengan deadline string dibuat sebelum kita implementasi FlexibleTime, atau dibuat langsung via MongoDB insert tanpa melalui API.

### 3. **Apa Solusinya?**

Ada **2 solusi**:

**SOLUSI 1 (Recommended):** Migrate data yang sudah ada di MongoDB
**SOLUSI 2:** Gunakan type `time.Time` biasa (tidak flexible)

---

## ✅ SOLUSI 1: Migrate Existing Data (RECOMMENDED)

### Cara Manual via MongoDB Shell

```bash
# Connect ke MongoDB
mongosh mongodb://dhfai:dhfai@103.151.145.172:27017/journaling

# Jalankan command ini
use journaling

// Lihat tasks dengan deadline string
db.tasks.find({ deadline: { $type: "string" } }).pretty()

// Convert semua deadline string ke DateTime
db.tasks.find({ deadline: { $type: "string" } }).forEach(function(task) {
    var dateStr = task.deadline;
    var dateObj = new Date(dateStr + "T00:00:00.000Z");

    db.tasks.updateOne(
        { _id: task._id },
        {
            $set: {
                deadline: dateObj,
                updated_at: new Date()
            }
        }
    );

    print("✅ Updated: " + task._id + " - " + task.title);
});

// Verify hasilnya
db.tasks.find({ deadline: { $exists: true } }).pretty()
```

### Cara Otomatis via Script

```bash
# Jalankan migration script
cd /home/dhfai/Documents/Program/Backend/backend-journaling
./scripts/migrate-dates.sh
```

**Output yang diharapkan:**
```
🔄 Starting MongoDB migration...
📦 Database: journaling

🔍 Finding tasks with string deadline...
Found 1 tasks to migrate
✅ Updated task: 69002010e26d267da934b3eb (dengan Nunu)

📊 Migration Summary:
   ✅ Successfully updated: 1
   ❌ Failed: 0
   📝 Total: 1

✨ Migration completed!
```

---

## ✅ SOLUSI 2: Simplified Approach (Alternatif)

Jika tidak ingin repot dengan FlexibleTime, gunakan `time.Time` biasa tapi dengan custom handler di service layer.

### Update Model

```go
// internal/models/mongo.go
type Task struct {
    ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    UserID        string             `bson:"user_id" json:"user_id"`
    Title         string             `bson:"title" json:"title"`
    DescriptionMD *string            `bson:"description_md,omitempty" json:"description_md,omitempty"`
    Status        string             `bson:"status" json:"status"`
    Priority      string             `bson:"priority" json:"priority"`
    Deadline      *time.Time         `bson:"deadline,omitempty" json:"deadline,omitempty"`  // ✅ time.Time biasa
    Tags          []string           `bson:"tags" json:"tags"`
    CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}
```

### Update Handler untuk Parse String Date

```go
// internal/handlers/task.go
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
    claims := r.Context().Value("user").(*jwt.Claims)

    var req struct {
        Title         string      `json:"title"`
        DescriptionMD string      `json:"description_md,omitempty"`
        Priority      string      `json:"priority"`
        DeadlineStr   string      `json:"deadline,omitempty"`  // ✅ Terima sebagai string
        Tags          []string    `json:"tags"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        WriteError(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    // Parse deadline string ke time.Time
    var deadline *time.Time
    if req.DeadlineStr != "" {
        // Try multiple formats
        formats := []string{
            "2006-01-02",
            "2006-01-02T15:04:05Z",
            "2006-01-02T15:04:05",
        }

        for _, format := range formats {
            if t, err := time.Parse(format, req.DeadlineStr); err == nil {
                deadline = &t
                break
            }
        }

        if deadline == nil {
            WriteError(w, http.StatusBadRequest, "Invalid deadline format. Use: YYYY-MM-DD")
            return
        }
    }

    // Create task with parsed deadline
    task, err := h.service.CreateTaskWithDeadline(r.Context(), claims.UserID.String(),
        req.Title, req.DescriptionMD, "todo", req.Priority, req.Tags, deadline)
    if err != nil {
        WriteError(w, http.StatusInternalServerError, "Failed to create task")
        return
    }

    WriteJSON(w, http.StatusCreated, task)
}
```

---

## 📝 Format Deadline yang Diterima

### Dari Frontend (JSON Request)

**✅ Format yang DITERIMA:**

```javascript
// 1. Date only (PALING UMUM)
{
  "title": "My Task",
  "deadline": "2025-10-31",  // ✅ YYYY-MM-DD
  "priority": "high"
}

// 2. DateTime dengan timezone
{
  "title": "My Task",
  "deadline": "2025-10-31T14:30:00Z",  // ✅ ISO 8601
  "priority": "high"
}

// 3. DateTime tanpa timezone
{
  "title": "My Task",
  "deadline": "2025-10-31T14:30:00",  // ✅ ISO 8601 local
  "priority": "high"
}

// 4. Null (no deadline)
{
  "title": "My Task",
  "deadline": null,  // ✅ Null untuk no deadline
  "priority": "high"
}
```

**❌ Format yang TIDAK DITERIMA:**

```javascript
// ❌ Format Indonesia
{
  "deadline": "31/10/2025"  // ERROR!
}

// ❌ Text format
{
  "deadline": "31 Oktober 2025"  // ERROR!
}

// ❌ Empty string
{
  "deadline": ""  // ERROR! Use null instead
}
```

### Disimpan di MongoDB

Setelah diproses backend, deadline akan disimpan sebagai **DateTime object**:

```javascript
{
  "_id": ObjectId("69002010e26d267da934b3eb"),
  "title": "My Task",
  "deadline": ISODate("2025-10-31T00:00:00.000Z"),  // ✅ DateTime object
  "created_at": ISODate("2025-10-28T01:44:48.844Z"),
  "updated_at": ISODate("2025-10-28T01:44:48.844Z")
}
```

### Response ke Frontend (JSON Response)

Backend mengembalikan dalam format ISO 8601:

```json
{
  "id": "69002010e26d267da934b3eb",
  "title": "My Task",
  "deadline": "2025-10-31T00:00:00Z",  // ✅ ISO 8601 string
  "created_at": "2025-10-28T01:44:48.844Z",
  "updated_at": "2025-10-28T01:44:48.844Z"
}
```

---

## 🧪 Testing

### 1. Test Create Task dengan Deadline

```bash
# Test dengan date only
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Task",
    "priority": "high",
    "deadline": "2025-10-31",
    "tags": []
  }'
```

**Expected Response:**
```json
{
  "id": "...",
  "title": "Test Task",
  "deadline": "2025-10-31T00:00:00Z",
  "priority": "high",
  "status": "todo",
  "tags": [],
  "created_at": "...",
  "updated_at": "..."
}
```

### 2. Test Get All Tasks

```bash
curl -X GET http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Expected Response:**
```json
[
  {
    "id": "69002010e26d267da934b3eb",
    "title": "dengan Nunu",
    "deadline": "2025-10-31T00:00:00Z",  // ✅ Sekarang format DateTime
    "status": "in_progress",
    "priority": "high"
  }
]
```

### 3. Verify di MongoDB

```bash
mongosh mongodb://dhfai:dhfai@103.151.145.172:27017/journaling --eval \
  'db.tasks.find({deadline: {$exists: true}}).pretty()'
```

**Expected Output:**
```javascript
{
  _id: ObjectId('69002010e26d267da934b3eb'),
  title: 'dengan Nunu',
  deadline: ISODate('2025-10-31T00:00:00.000Z'),  // ✅ DateTime object
  priority: 'high',
  status: 'in_progress'
}
```

---

## 🔍 Debug Steps

### Step 1: Cek Data di MongoDB

```bash
mongosh mongodb://dhfai:dhfai@103.151.145.172:27017/journaling

use journaling

// Cek type dari deadline
db.tasks.find({deadline: {$exists: true}}).forEach(function(task) {
    print("Task ID: " + task._id);
    print("Title: " + task.title);
    print("Deadline: " + task.deadline);
    print("Deadline Type: " + typeof task.deadline);
    print("---");
});
```

**Output yang SALAH:**
```
Task ID: 69002010e26d267da934b3eb
Title: dengan Nunu
Deadline: 2025-10-31
Deadline Type: string  ❌ INI MASALAHNYA!
```

**Output yang BENAR:**
```
Task ID: 69002010e26d267da934b3eb
Title: dengan Nunu
Deadline: 2025-10-31T00:00:00.000Z
Deadline Type: object  ✅ BENAR!
```

### Step 2: Test API

```bash
# Restart backend
go run main.go

# Test GET tasks
curl -X GET http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Step 3: Check Backend Logs

Jika masih error, cek log backend:
```
2025/10/28 09:34:17 [dhfai/Q39GBYnwBL-000001] "GET http://localhost:8080/api/v1/tasks HTTP/1.1" from [::1]:43200 - 500 150B in 38.155444ms
```

Error 500 berarti ada masalah di decode BSON dari MongoDB.

---

## 💡 Rekomendasi

**PILIH SOLUSI 1** jika:
- ✅ Mau support multiple date formats dari frontend
- ✅ Ada data existing yang perlu dimigrate
- ✅ Mau flexible untuk future changes

**PILIH SOLUSI 2** jika:
- ✅ Mau simple dan straightforward
- ✅ Hanya butuh format date standard
- ✅ Tidak masalah rebuild database dari scratch

---

## 📞 Next Steps

1. **Backup data dulu!**
   ```bash
   mongodump --uri="mongodb://dhfai:dhfai@103.151.145.172:27017/journaling" --out=backup
   ```

2. **Jalankan migration:**
   ```bash
   ./scripts/migrate-dates.sh
   ```

3. **Rebuild & restart backend:**
   ```bash
   go build -o bin/backend-journaling
   ./bin/backend-journaling
   ```

4. **Test API:**
   ```bash
   curl -X GET http://localhost:8080/api/v1/tasks -H "Authorization: Bearer TOKEN"
   ```

5. **Jika masih error, laporkan:**
   - Error message lengkap
   - Output dari MongoDB query
   - Sample data yang bermasalah

---

**Good Luck! 🚀**
