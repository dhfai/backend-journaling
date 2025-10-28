# 📚 Documentation Index

Welcome to Backend Journaling API documentation! This index will help you find the information you need quickly.

---

## 🚀 Getting Started

**New to the project? Start here:**

1. **[Quick Start Guide](./QUICK_START.md)** ⭐
   - Setup in 5 minutes
   - Basic examples
   - Common troubleshooting

2. **[Environment Variables](./ENVIRONMENT_VARIABLES.md)**
   - Complete configuration guide
   - Security best practices
   - Environment-specific settings

---

## 📖 API Reference

### Complete Documentation
**[API Documentation](./API_DOCUMENTATION.md)** - 500+ lines of comprehensive documentation

**What's inside:**
- ✅ All endpoint specifications
- ✅ Request/Response formats
- ✅ Error handling
- ✅ Authentication flow
- ✅ cURL examples
- ✅ Best practices

### Quick Reference

#### 📝 Notes API
- Create, read, update, delete notes
- Add/remove/reorder blocks
- Support for paragraph, heading, and todo blocks
- [Full Notes API Docs →](./API_DOCUMENTATION.md#notes-api)

#### ✅ Todos API
- Simple todo management
- Priority levels
- Due dates
- [Full Todos API Docs →](./API_DOCUMENTATION.md#todos-api)

#### 📊 Tasks API
- Complex task tracking
- Status workflow
- Markdown descriptions
- [Full Tasks API Docs →](./API_DOCUMENTATION.md#tasks-api)

---

## 🎯 Feature Guides

### Reorder Blocks
**[Reorder Blocks Guide](./REORDER_BLOCKS_GUIDE.md)**

Detailed guide for understanding and implementing block reordering:
- Step-by-step workflow
- Common mistakes to avoid
- Error handling
- Complete examples

**When to read:** When implementing drag-and-drop or block reordering in your frontend.

### CORS Configuration
**[CORS Configuration Guide](./CORS_CONFIGURATION.md)**

Complete guide for Cross-Origin Resource Sharing:
- Development vs Production setup
- Frontend integration examples
- Troubleshooting CORS errors
- Security best practices

**When to read:** When integrating frontend or encountering CORS errors.

---

## 🏗️ Development

### Implementation Summary
**[Implementation Summary](./IMPLEMENTATION_SUMMARY.md)**

Technical deep-dive for developers:
- Architecture decisions
- Database schema
- Code organization
- Performance considerations
- Bug fixes and solutions

**When to read:**
- Contributing to the project
- Understanding the codebase
- Learning from implementation decisions

---

## 📋 Quick Access by Role

### 👨‍💻 Frontend Developer
**What you need:**
1. [Quick Start Guide](./QUICK_START.md) - Setup and test
2. [API Documentation](./API_DOCUMENTATION.md) - Endpoint reference
3. [CORS Configuration](./CORS_CONFIGURATION.md) - Fix CORS errors
4. [Reorder Blocks Guide](./REORDER_BLOCKS_GUIDE.md) - Feature implementation
5. Postman Collection: `../api_test/postman_collection.json`

### 🔧 Backend Developer
**What you need:**
1. [Implementation Summary](./IMPLEMENTATION_SUMMARY.md) - Architecture
2. [Environment Variables](./ENVIRONMENT_VARIABLES.md) - Configuration
3. Main README - Project structure
4. Source code in `internal/` directory

### 🚀 DevOps Engineer
**What you need:**
1. [Environment Variables](./ENVIRONMENT_VARIABLES.md) - Configuration
2. [Quick Start Guide](./QUICK_START.md) - Deployment basics
3. Docker configurations (if available)
4. Production checklist (in Environment Variables doc)

### 📱 Mobile Developer
**What you need:**
1. [API Documentation](./API_DOCUMENTATION.md) - Endpoint reference
2. [Quick Start Guide](./QUICK_START.md) - Testing
3. Authentication flow guide
4. Example JSON files in `../examples/`

### 🧪 QA Engineer
**What you need:**
1. [API Documentation](./API_DOCUMENTATION.md) - Test cases
2. Postman Collection: `../api_test/postman_collection.json`
3. [Reorder Blocks Guide](./REORDER_BLOCKS_GUIDE.md) - Edge cases
4. Error response reference

---

## 🔍 Find Information By Topic

### Authentication
- Setup: [Quick Start → Authentication Flow](./QUICK_START.md#authentication-flow)
- Endpoints: [API Docs → Authentication](./API_DOCUMENTATION.md#authentication)
- Configuration: [Environment Variables → JWT](./ENVIRONMENT_VARIABLES.md#jwt-configuration)

### CORS & Frontend Integration
- Configuration: [CORS Configuration](./CORS_CONFIGURATION.md)
- Troubleshooting: [CORS → Troubleshooting](./CORS_CONFIGURATION.md#troubleshooting)
- Frontend Examples: [CORS → Frontend Integration](./CORS_CONFIGURATION.md#frontend-integration)

### Database
- MongoDB Setup: [Quick Start → Setup](./QUICK_START.md#setup)
- Configuration: [Environment Variables → Database](./ENVIRONMENT_VARIABLES.md#database-configuration)
- Schema Design: [Implementation Summary → Data Models](./IMPLEMENTATION_SUMMARY.md#2-data-models)

### Error Handling
- Error Responses: [API Docs → Error Responses](./API_DOCUMENTATION.md#error-responses)
- Status Codes: [API Docs → Status Codes](./API_DOCUMENTATION.md#status-codes)
- Troubleshooting: [Quick Start → Common Issues](./QUICK_START.md#common-issues)

### Security
- JWT: [Environment Variables → JWT Configuration](./ENVIRONMENT_VARIABLES.md#jwt-configuration)
- Best Practices: [Environment Variables → Security](./ENVIRONMENT_VARIABLES.md#security-best-practices)
- User Isolation: [Implementation Summary → Security](./IMPLEMENTATION_SUMMARY.md#security-implementation)

---

## 📦 Additional Resources

### Example Files
Located in `../examples/` directory:
- Note examples (create, update, blocks)
- Todo examples (create, update)
- Task examples (create, update)

### Postman Collection
Located in `../api_test/`:
- `postman_collection.json` - All endpoints
- `postman_environment_local.json` - Local environment
- `postman_environment_production.json` - Production environment

### Source Code
- `../internal/handlers/` - HTTP handlers
- `../internal/service/` - Business logic
- `../internal/repository/` - Data access
- `../internal/models/` - Data models

---

## 🆘 Getting Help

### Can't find what you're looking for?

1. **Search the documentation**
   - Use Ctrl+F in your browser
   - Check the table of contents in each doc

2. **Check examples**
   - Postman collection has working examples
   - Example JSON files show request formats

3. **Common Issues**
   - [Quick Start → Common Issues](./QUICK_START.md#common-issues)
   - [API Docs → Error Responses](./API_DOCUMENTATION.md#error-responses)

4. **Report Issues**
   - Create issue on GitHub
   - Include error messages and steps to reproduce

---

## 📊 Documentation Stats

| Document | Lines | Topics Covered |
|----------|-------|----------------|
| API Documentation | 500+ | All endpoints, examples, errors |
| Quick Start Guide | 200+ | Setup, quick examples, troubleshooting |
| Environment Variables | 300+ | All config options, security |
| Implementation Summary | 400+ | Architecture, decisions, stats |
| Reorder Blocks Guide | 200+ | Workflow, examples, errors |
| CORS Configuration | 400+ | Setup, integration, troubleshooting |

**Total:** 2,000+ lines of documentation

---

## 🗺️ Documentation Roadmap

### Completed ✅
- ✅ API endpoint documentation
- ✅ Quick start guide
- ✅ Environment configuration
- ✅ Feature guides (Reorder Blocks)
- ✅ CORS configuration guide
- ✅ Implementation details
- ✅ Postman collection

### Planned 📋
- [ ] Video tutorials
- [ ] Migration guides
- [ ] Performance tuning guide
- [ ] Scaling guide
- [ ] API versioning strategy

---

## 📝 Contributing to Documentation

### Found an error or want to improve docs?

1. Documentation is in Markdown format
2. Follow existing structure and style
3. Include code examples where helpful
4. Update this index when adding new docs

### Style Guide
- Use clear, concise language
- Include examples for complex topics
- Use emojis for visual navigation
- Keep code blocks properly formatted
- Test all code examples

---

## 🔗 External Links

- **Go Documentation:** https://golang.org/doc/
- **MongoDB Docs:** https://docs.mongodb.com/
- **PostgreSQL Docs:** https://www.postgresql.org/docs/
- **Chi Router:** https://github.com/go-chi/chi
- **JWT Specification:** https://jwt.io/

---

**Happy Building! 🚀**

*Last Updated: October 28, 2025*
