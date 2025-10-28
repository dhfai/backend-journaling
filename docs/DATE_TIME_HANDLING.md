# 📅 Date and Time Handling Guide

## Overview

Backend Journaling API menggunakan flexible date/time parsing yang mendukung berbagai format input untuk memudahkan integrasi frontend.

---

## 🎯 Supported Formats

### Input Formats (JSON Request)

API menerima berbagai format untuk field `deadline` (Task) dan `due_date` (Todo):

| Format | Example | Use Case |
|--------|---------|----------|
| **Date Only** | `"2025-10-15"` | Paling umum untuk deadline/due date |
| **ISO 8601 Full** | `"2025-10-15T14:30:00Z"` | Dengan timezone UTC |
| **ISO 8601 Local** | `"2025-10-15T14:30:00"` | Tanpa timezone |
| **ISO 8601 with TZ** | `"2025-10-15T14:30:00+07:00"` | Dengan timezone offset |
| **DateTime Space** | `"2025-10-15 14:30:00"` | Format alternatif |

### Output Format (JSON Response)

API selalu mengembalikan dalam format ISO 8601 lengkap:

```json
{
  "deadline": "2025-10-15T00:00:00Z"
}
```

---

## 📝 Field Reference

### Task Model

```go
type Task struct {
    ID            string    `json:"id"`
    UserID        string    `json:"user_id"`
    Title         string    `json:"title"`
    DescriptionMD string    `json:"description_md,omitempty"`
    Status        string    `json:"status"`
    Priority      string    `json:"priority"`
    Deadline      *time     `json:"deadline,omitempty"`      // ⏰ Flexible format
    Tags          []string  `json:"tags"`
    CreatedAt     time      `json:"created_at"`
    UpdatedAt     time      `json:"updated_at"`
}
```

### Todo Model

```go
type Todo struct {
    ID        string   `json:"id"`
    UserID    string   `json:"user_id"`
    Title     string   `json:"title"`
    Done      bool     `json:"done"`
    Priority  string   `json:"priority"`
    DueDate   *time    `json:"due_date,omitempty"`          // ⏰ Flexible format
    CreatedAt time     `json:"created_at"`
    UpdatedAt time     `json:"updated_at"`
}
```

---

## 💻 Frontend Examples

### JavaScript / Fetch API

#### Creating Task with Date Only

```javascript
const createTask = async (title, deadline) => {
  const response = await fetch('http://localhost:8080/api/v1/tasks', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({
      title: title,
      priority: 'high',
      deadline: deadline,  // "2025-10-15" works!
      tags: []
    })
  });

  return await response.json();
};

// Usage
await createTask('Complete project', '2025-10-15');
```

#### Creating Task with Full DateTime

```javascript
const createTask = async (title, deadline) => {
  const response = await fetch('http://localhost:8080/api/v1/tasks', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({
      title: title,
      priority: 'high',
      deadline: deadline,  // "2025-10-15T14:30:00Z" also works!
      tags: []
    })
  });

  return await response.json();
};

// Usage
await createTask('Meeting', '2025-10-15T14:30:00Z');
```

### React Date Input

#### Using HTML Date Input

```jsx
import { useState } from 'react';

function TaskForm() {
  const [title, setTitle] = useState('');
  const [deadline, setDeadline] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();

    // deadline akan dalam format "2025-10-15" dari input[type="date"]
    const response = await fetch('http://localhost:8080/api/v1/tasks', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({
        title,
        deadline,  // "2025-10-15" format
        priority: 'medium',
        tags: []
      })
    });

    const data = await response.json();
    console.log('Task created:', data);
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="text"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        placeholder="Task title"
      />
      <input
        type="date"
        value={deadline}
        onChange={(e) => setDeadline(e.target.value)}
      />
      <button type="submit">Create Task</button>
    </form>
  );
}
```

#### Using DateTime Input

```jsx
function TaskForm() {
  const [title, setTitle] = useState('');
  const [deadline, setDeadline] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();

    // deadline akan dalam format "2025-10-15T14:30" dari input[type="datetime-local"]
    const response = await fetch('http://localhost:8080/api/v1/tasks', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({
        title,
        deadline,  // "2025-10-15T14:30" format
        priority: 'medium',
        tags: []
      })
    });

    const data = await response.json();
    console.log('Task created:', data);
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="text"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        placeholder="Task title"
      />
      <input
        type="datetime-local"
        value={deadline}
        onChange={(e) => setDeadline(e.target.value)}
      />
      <button type="submit">Create Task</button>
    </form>
  );
}
```

### Using Date Libraries

#### Moment.js

```javascript
import moment from 'moment';

const createTask = async (title, deadlineDate) => {
  // Format berbagai cara, semua diterima!
  const formats = [
    deadlineDate.format('YYYY-MM-DD'),           // "2025-10-15"
    deadlineDate.format('YYYY-MM-DDTHH:mm:ss'),  // "2025-10-15T14:30:00"
    deadlineDate.toISOString(),                  // "2025-10-15T14:30:00.000Z"
  ];

  const response = await fetch('http://localhost:8080/api/v1/tasks', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({
      title: title,
      deadline: formats[0],  // Pilih format yang diinginkan
      priority: 'high',
      tags: []
    })
  });

  return await response.json();
};

// Usage
const tomorrow = moment().add(1, 'days');
await createTask('Complete report', tomorrow);
```

#### Day.js

```javascript
import dayjs from 'dayjs';

const createTask = async (title, deadlineDate) => {
  const response = await fetch('http://localhost:8080/api/v1/tasks', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({
      title: title,
      deadline: deadlineDate.format('YYYY-MM-DD'),  // "2025-10-15"
      priority: 'high',
      tags: []
    })
  });

  return await response.json();
};

// Usage
const nextWeek = dayjs().add(7, 'day');
await createTask('Review code', nextWeek);
```

#### Date-fns

```javascript
import { format } from 'date-fns';

const createTask = async (title, deadlineDate) => {
  const response = await fetch('http://localhost:8080/api/v1/tasks', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({
      title: title,
      deadline: format(deadlineDate, 'yyyy-MM-dd'),  // "2025-10-15"
      priority: 'high',
      tags: []
    })
  });

  return await response.json();
};

// Usage
const deadline = new Date(2025, 9, 15);  // October 15, 2025
await createTask('Submit proposal', deadline);
```

---

## 🔄 Updating Dates

### PATCH Task Deadline

```javascript
const updateDeadline = async (taskId, newDeadline) => {
  const response = await fetch(`http://localhost:8080/api/v1/tasks/${taskId}`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({
      deadline: newDeadline  // Any supported format
    })
  });

  return await response.json();
};

// Usage
await updateDeadline('670009c234e0eff411bf6f83', '2025-10-20');
```

### PATCH Todo Due Date

```javascript
const updateDueDate = async (todoId, newDueDate) => {
  const response = await fetch(`http://localhost:8080/api/v1/todos/${todoId}`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({
      due_date: newDueDate  // Any supported format
    })
  });

  return await response.json();
};

// Usage
await updateDueDate('670009c234e0eff411bf6f84', '2025-10-18');
```

---

## 📊 Parsing Response Dates

### JavaScript Date Object

```javascript
const response = await fetch('http://localhost:8080/api/v1/tasks/123', {
  headers: { 'Authorization': `Bearer ${token}` }
});

const task = await response.json();

// deadline dalam format: "2025-10-15T00:00:00Z"
const deadlineDate = new Date(task.deadline);

console.log('Deadline:', deadlineDate.toLocaleDateString());
// Output: "10/15/2025"
```

### Display in UI

```jsx
function TaskItem({ task }) {
  const formatDate = (dateString) => {
    if (!dateString) return 'No deadline';

    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  return (
    <div>
      <h3>{task.title}</h3>
      <p>Deadline: {formatDate(task.deadline)}</p>
      <p>Status: {task.status}</p>
    </div>
  );
}
```

### Using Date Libraries

```javascript
import moment from 'moment';

const displayTask = (task) => {
  const deadline = moment(task.deadline);

  console.log('Formatted:', deadline.format('MMMM DD, YYYY'));
  // Output: "October 15, 2025"

  console.log('From now:', deadline.fromNow());
  // Output: "in 2 days"

  console.log('Relative:', deadline.calendar());
  // Output: "Tomorrow at 12:00 AM"
};
```

---

## ⚠️ Common Issues & Solutions

### Issue 1: "Cannot parse date"

**Problem:**
```json
{
  "success": false,
  "error": "error decoding key deadline: parsing time..."
}
```

**Solution:**
Gunakan salah satu format yang didukung:
- ✅ `"2025-10-15"`
- ✅ `"2025-10-15T14:30:00Z"`
- ❌ `"15/10/2025"` (tidak didukung)
- ❌ `"Oct 15, 2025"` (tidak didukung)

### Issue 2: Null vs Empty String

**Problem:**
```javascript
// ❌ Wrong
{ "deadline": "" }  // Empty string akan error

// ✅ Correct
{ "deadline": null }  // Null untuk no deadline
// or
{}  // Omit field entirely
```

**Solution:**
```javascript
const deadline = formData.deadline || null;  // Convert empty to null
```

### Issue 3: Timezone Issues

**Problem:**
Date tampil berbeda di frontend dan backend

**Solution:**
```javascript
// Store in UTC
const deadline = new Date('2025-10-15');
const deadlineUTC = deadline.toISOString();  // "2025-10-15T00:00:00.000Z"

// Display in user's timezone
const displayDate = new Date(task.deadline).toLocaleString();
```

---

## 📝 Best Practices

### 1. **Use Date-Only for Deadlines**

```javascript
// ✅ Recommended for tasks/todos
{
  "title": "Complete report",
  "deadline": "2025-10-15"  // Simple date
}
```

### 2. **Use Full DateTime for Events**

```javascript
// ✅ For specific time requirements
{
  "title": "Team meeting",
  "deadline": "2025-10-15T14:30:00Z"  // Includes time
}
```

### 3. **Handle Null Gracefully**

```javascript
const formatDeadline = (deadline) => {
  if (!deadline) return 'No deadline';
  return new Date(deadline).toLocaleDateString();
};
```

### 4. **Validate Before Sending**

```javascript
const isValidDate = (dateString) => {
  if (!dateString) return true;  // null/undefined is valid
  const date = new Date(dateString);
  return !isNaN(date.getTime());
};

if (!isValidDate(formData.deadline)) {
  alert('Invalid date format');
  return;
}
```

### 5. **Use Timezone-Aware Libraries**

```javascript
// ✅ Using moment-timezone
import moment from 'moment-timezone';

const userTimezone = moment.tz.guess();  // "Asia/Jakarta"
const deadline = moment.tz('2025-10-15', userTimezone);
```

---

## 🧪 Testing Examples

### cURL Tests

**Create task with date:**
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Task",
    "priority": "high",
    "deadline": "2025-10-15",
    "tags": []
  }'
```

**Create task with datetime:**
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Task",
    "priority": "high",
    "deadline": "2025-10-15T14:30:00Z",
    "tags": []
  }'
```

**Update deadline:**
```bash
curl -X PATCH http://localhost:8080/api/v1/tasks/TASK_ID \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "deadline": "2025-10-20"
  }'
```

---

## 📚 Additional Resources

- **ISO 8601 Standard:** https://en.wikipedia.org/wiki/ISO_8601
- **JavaScript Date:** https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date
- **Moment.js:** https://momentjs.com/
- **Day.js:** https://day.js.org/
- **Date-fns:** https://date-fns.org/

---

**Happy Coding! 🚀**

*Last Updated: October 28, 2025*
