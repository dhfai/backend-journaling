# 🌐 CORS Configuration Guide

## Overview

CORS (Cross-Origin Resource Sharing) adalah security feature yang mengatur akses API dari domain berbeda. Backend journaling API menggunakan CORS middleware yang dapat dikonfigurasi untuk development dan production.

---

## 🔧 Configuration

### Environment Variables

```bash
# For Development - Allow all origins
CORS_ALLOWED_ORIGINS=*
CORS_ALLOWED_METHODS=GET, POST, PUT, PATCH, DELETE, OPTIONS
CORS_ALLOWED_HEADERS=Content-Type, Authorization, X-Requested-With

# For Production - Specific domains
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://app.yourdomain.com
CORS_ALLOWED_METHODS=GET, POST, PUT, PATCH, DELETE, OPTIONS
CORS_ALLOWED_HEADERS=Content-Type, Authorization, X-Requested-With
```

### Default Values

Jika tidak di-set, aplikasi menggunakan default:

```go
AllowedOrigins: "*"  // Allow all origins
AllowedMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS"
AllowedHeaders: "Content-Type, Authorization, X-Requested-With"
MaxAge: 86400  // 24 hours cache
```

---

## 🎯 HTTP Methods Supported

| Method | Usage | Example Endpoints |
|--------|-------|-------------------|
| `GET` | Retrieve data | GET /api/v1/notes |
| `POST` | Create new resource | POST /api/v1/notes |
| `PUT` | Update entire resource | PUT /api/v1/profile |
| `PATCH` | Partial update | PATCH /api/v1/notes/{id} |
| `DELETE` | Delete resource | DELETE /api/v1/notes/{id} |
| `OPTIONS` | Preflight request | OPTIONS /api/v1/notes |

---

## 🔐 Security Recommendations

### Development Environment

**Permissive Configuration:**
```bash
CORS_ALLOWED_ORIGINS=*
```

✅ **Pros:**
- Easy testing from localhost
- Works with any development tool
- No configuration needed

⚠️ **Cons:**
- Less secure
- Not for production use

### Production Environment

**Strict Configuration:**
```bash
CORS_ALLOWED_ORIGINS=https://app.example.com
```

✅ **Pros:**
- More secure
- Prevents unauthorized access
- Complies with security best practices

⚠️ **Configuration Required:**
- Must specify exact domains
- No wildcards in production

### Multiple Origins

For multiple frontend domains:

```bash
CORS_ALLOWED_ORIGINS=https://app.example.com,https://dashboard.example.com,https://mobile.example.com
```

**Format:** Comma-separated list, no spaces

---

## 🚀 Frontend Integration

### JavaScript / Fetch API

```javascript
// GET Request
const response = await fetch('http://localhost:8080/api/v1/notes', {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${accessToken}`
  }
});

// POST Request
const response = await fetch('http://localhost:8080/api/v1/notes', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${accessToken}`
  },
  body: JSON.stringify({
    title: 'My Note',
    content: 'Note content'
  })
});

// PATCH Request
const response = await fetch(`http://localhost:8080/api/v1/notes/${noteId}`, {
  method: 'PATCH',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${accessToken}`
  },
  body: JSON.stringify({
    title: 'Updated Title'
  })
});
```

### Axios

```javascript
import axios from 'axios';

// Configure base URL
const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
  headers: {
    'Content-Type': 'application/json'
  }
});

// Add auth token to all requests
api.interceptors.request.use(config => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// GET Request
const response = await api.get('/notes');

// POST Request
const response = await api.post('/notes', {
  title: 'My Note',
  content: 'Note content'
});

// PATCH Request
const response = await api.patch(`/notes/${noteId}`, {
  title: 'Updated Title'
});

// DELETE Request
const response = await api.delete(`/notes/${noteId}`);
```

### React Example

```javascript
import { useState, useEffect } from 'react';

function NotesPage() {
  const [notes, setNotes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetchNotes();
  }, []);

  const fetchNotes = async () => {
    try {
      const token = localStorage.getItem('access_token');
      const response = await fetch('http://localhost:8080/api/v1/notes', {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (!response.ok) {
        throw new Error('Failed to fetch notes');
      }

      const data = await response.json();
      setNotes(data.data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const createNote = async (noteData) => {
    try {
      const token = localStorage.getItem('access_token');
      const response = await fetch('http://localhost:8080/api/v1/notes', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(noteData)
      });

      if (!response.ok) {
        throw new Error('Failed to create note');
      }

      const data = await response.json();
      setNotes([...notes, data.data]);
      return data.data;
    } catch (err) {
      setError(err.message);
      throw err;
    }
  };

  const updateNote = async (noteId, updates) => {
    try {
      const token = localStorage.getItem('access_token');
      const response = await fetch(`http://localhost:8080/api/v1/notes/${noteId}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(updates)
      });

      if (!response.ok) {
        throw new Error('Failed to update note');
      }

      const data = await response.json();
      setNotes(notes.map(note =>
        note.id === noteId ? data.data : note
      ));
      return data.data;
    } catch (err) {
      setError(err.message);
      throw err;
    }
  };

  const deleteNote = async (noteId) => {
    try {
      const token = localStorage.getItem('access_token');
      const response = await fetch(`http://localhost:8080/api/v1/notes/${noteId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (!response.ok) {
        throw new Error('Failed to delete note');
      }

      setNotes(notes.filter(note => note.id !== noteId));
    } catch (err) {
      setError(err.message);
      throw err;
    }
  };

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div>
      <h1>Notes</h1>
      {notes.map(note => (
        <div key={note.id}>
          <h2>{note.title}</h2>
          <p>{note.content}</p>
          <button onClick={() => updateNote(note.id, { title: 'Updated' })}>
            Update
          </button>
          <button onClick={() => deleteNote(note.id)}>
            Delete
          </button>
        </div>
      ))}
    </div>
  );
}
```

---

## 🔍 Troubleshooting

### Error: "No 'Access-Control-Allow-Origin' header"

**Problem:** Browser blocks request karena CORS policy

**Solution:**
```bash
# Check server is running
curl http://localhost:8080/api/v1/health

# Check CORS headers
curl -i -X OPTIONS http://localhost:8080/api/v1/notes \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: GET"

# Should return:
# Access-Control-Allow-Origin: *
# Access-Control-Allow-Methods: GET, POST, PUT, PATCH, DELETE, OPTIONS
```

### Error: "Method not allowed in CORS policy"

**Problem:** HTTP method tidak ada di `CORS_ALLOWED_METHODS`

**Solution:**
```bash
# Add missing method to .env
CORS_ALLOWED_METHODS=GET, POST, PUT, PATCH, DELETE, OPTIONS
```

### Error: "Request header not allowed"

**Problem:** Custom header tidak ada di `CORS_ALLOWED_HEADERS`

**Solution:**
```bash
# Add missing header to .env
CORS_ALLOWED_HEADERS=Content-Type, Authorization, X-Requested-With, X-Custom-Header
```

### Preflight Request Fails

**Problem:** OPTIONS request returns error

**Solution:**
1. Check server logs
2. Verify OPTIONS method is allowed
3. Check preflight request headers:

```bash
curl -i -X OPTIONS http://localhost:8080/api/v1/notes \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: PATCH" \
  -H "Access-Control-Request-Headers: Content-Type, Authorization"
```

Expected response:
- Status: `204 No Content`
- Headers: All CORS headers present

### Browser Cache Issues

**Problem:** Old CORS policy cached by browser

**Solution:**
1. Clear browser cache
2. Hard reload (Ctrl+Shift+R)
3. Use incognito/private mode
4. Check `Access-Control-Max-Age` value

---

## 📊 CORS Flow

### Simple Request (GET, POST with simple headers)

```
1. Browser → Server: GET /api/v1/notes
   Headers: Authorization: Bearer token

2. Server → Browser: 200 OK
   Headers: Access-Control-Allow-Origin: *
   Body: Notes data
```

### Preflight Request (PUT, PATCH, DELETE or custom headers)

```
1. Browser → Server: OPTIONS /api/v1/notes/123
   Headers:
     Origin: http://localhost:3000
     Access-Control-Request-Method: PATCH
     Access-Control-Request-Headers: Content-Type, Authorization

2. Server → Browser: 204 No Content
   Headers:
     Access-Control-Allow-Origin: *
     Access-Control-Allow-Methods: GET, POST, PUT, PATCH, DELETE, OPTIONS
     Access-Control-Allow-Headers: Content-Type, Authorization, X-Requested-With
     Access-Control-Max-Age: 86400

3. Browser → Server: PATCH /api/v1/notes/123
   Headers:
     Content-Type: application/json
     Authorization: Bearer token
   Body: { "title": "Updated" }

4. Server → Browser: 200 OK
   Headers: Access-Control-Allow-Origin: *
   Body: Updated note data
```

---

## 🧪 Testing CORS

### Using cURL

**Test GET request:**
```bash
curl -i http://localhost:8080/api/v1/notes \
  -H "Origin: http://localhost:3000" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Test OPTIONS (preflight):**
```bash
curl -i -X OPTIONS http://localhost:8080/api/v1/notes \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: PATCH" \
  -H "Access-Control-Request-Headers: Content-Type, Authorization"
```

**Test PATCH request:**
```bash
curl -i -X PATCH http://localhost:8080/api/v1/notes/123 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Origin: http://localhost:3000" \
  -d '{"title":"Updated"}'
```

### Using Browser DevTools

1. Open DevTools (F12)
2. Go to **Network** tab
3. Make API request
4. Click on request
5. Check **Headers** section:
   - Request Headers: Origin, Access-Control-Request-Method
   - Response Headers: Access-Control-Allow-*

### Expected Response Headers

```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, PATCH, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization, X-Requested-With
Access-Control-Expose-Headers: Content-Length, Content-Type
Access-Control-Max-Age: 86400
Access-Control-Allow-Credentials: true
```

---

## 🔒 Production Checklist

### Before Deployment

- [ ] Set specific `CORS_ALLOWED_ORIGINS` (no wildcards)
- [ ] List only required methods in `CORS_ALLOWED_METHODS`
- [ ] Limit `CORS_ALLOWED_HEADERS` to necessary ones
- [ ] Set appropriate `MaxAge` for caching
- [ ] Test all endpoints from production domain
- [ ] Verify preflight requests work
- [ ] Check browser console for CORS errors
- [ ] Test with different browsers
- [ ] Document allowed origins for team

### Example Production Config

```bash
# .env.production
CORS_ALLOWED_ORIGINS=https://app.example.com
CORS_ALLOWED_METHODS=GET, POST, PUT, PATCH, DELETE, OPTIONS
CORS_ALLOWED_HEADERS=Content-Type, Authorization
ENVIRONMENT=production
```

---

## 📚 Additional Resources

- **MDN CORS Guide:** https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
- **CORS Specification:** https://www.w3.org/TR/cors/
- **CORS Tester:** https://www.test-cors.org/
- **Chrome DevTools Network:** https://developer.chrome.com/docs/devtools/network/

---

## 💡 Best Practices

1. **Development:**
   - Use `*` for quick testing
   - Enable detailed logging
   - Test with actual frontend URL

2. **Staging:**
   - Use specific staging domain
   - Test with production-like config
   - Verify all HTTP methods

3. **Production:**
   - Never use `*` wildcard
   - List specific domains
   - Monitor CORS errors
   - Keep config in secrets/vault

4. **Security:**
   - Validate origin on server side
   - Use HTTPS in production
   - Implement rate limiting
   - Log suspicious requests

---

**Happy Coding! 🚀**

*Last Updated: October 28, 2025*
