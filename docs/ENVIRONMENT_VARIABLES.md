# Environment Variables Documentation

Complete guide for all environment variables used in Backend Journaling API.

## 🗂️ Categories

- [Database Configuration](#database-configuration)
- [MongoDB Configuration](#mongodb-configuration)
- [JWT Configuration](#jwt-configuration)
- [OTP Configuration](#otp-configuration)
- [SMTP Configuration](#smtp-configuration)
- [Server Configuration](#server-configuration)

---

## Database Configuration

PostgreSQL database settings for authentication and user management.

### `DB_HOST`
- **Type:** String
- **Default:** `localhost`
- **Description:** PostgreSQL server hostname
- **Example:** `localhost`, `db.example.com`, `192.168.1.100`

### `DB_PORT`
- **Type:** Integer
- **Default:** `5432`
- **Description:** PostgreSQL server port
- **Example:** `5432`, `5433`

### `DB_USER`
- **Type:** String
- **Default:** `postgres`
- **Description:** PostgreSQL username
- **Example:** `postgres`, `journaling_user`

### `DB_PASSWORD`
- **Type:** String
- **Default:** `postgres`
- **Description:** PostgreSQL user password
- **Example:** `your_secure_password`
- **Security:** ⚠️ Keep this secret!

### `DB_NAME`
- **Type:** String
- **Default:** `journaling_auth`
- **Description:** PostgreSQL database name
- **Example:** `journaling_auth`, `auth_db`

### `DB_SSLMODE`
- **Type:** String
- **Default:** `disable`
- **Options:** `disable`, `require`, `verify-ca`, `verify-full`
- **Description:** PostgreSQL SSL connection mode
- **Production:** Use `require` or higher

---

## MongoDB Configuration

MongoDB settings for notes, todos, and tasks storage.

### `MONGO_URI`
- **Type:** String (URI)
- **Default:** `mongodb://localhost:27017`
- **Description:** MongoDB connection string
- **Examples:**
  ```
  mongodb://localhost:27017
  mongodb://username:password@localhost:27017
  mongodb://host1:27017,host2:27017/?replicaSet=mySet
  mongodb+srv://cluster.mongodb.net/mydb
  ```
- **Formats:**
  - Local: `mongodb://localhost:27017`
  - Auth: `mongodb://user:pass@host:27017`
  - Atlas: `mongodb+srv://user:pass@cluster.mongodb.net`
  - Replica Set: `mongodb://host1,host2,host3/?replicaSet=rs0`

### `MONGO_DATABASE`
- **Type:** String
- **Default:** `journaling`
- **Description:** MongoDB database name for storing notes, todos, tasks
- **Example:** `journaling`, `notes_db`, `production_data`

---

## JWT Configuration

JSON Web Token settings for authentication.

### `JWT_PRIVATE_KEY_PATH`
- **Type:** String (File Path)
- **Default:** `./keys/jwt_private.pem`
- **Description:** Path to RSA private key for signing JWTs
- **Example:** `./keys/jwt_private.pem`, `/etc/secrets/private.pem`
- **Security:** ⚠️ Keep this file secure and never commit to git!

### `JWT_PUBLIC_KEY_PATH`
- **Type:** String (File Path)
- **Default:** `./keys/jwt_public.pem`
- **Description:** Path to RSA public key for verifying JWTs
- **Example:** `./keys/jwt_public.pem`, `/etc/secrets/public.pem`

### `JWT_ACCESS_TOKEN_DURATION`
- **Type:** Duration String
- **Default:** `15m`
- **Description:** Access token expiration time
- **Examples:**
  - `15m` - 15 minutes
  - `1h` - 1 hour
  - `30s` - 30 seconds
- **Recommended:** `15m` to `30m` for production

### `JWT_REFRESH_TOKEN_DURATION`
- **Type:** Duration String
- **Default:** `168h` (7 days)
- **Description:** Refresh token expiration time
- **Examples:**
  - `168h` - 7 days
  - `720h` - 30 days
  - `2160h` - 90 days
- **Recommended:** `168h` to `720h` for production

---

## OTP Configuration

One-Time Password settings for email verification and password reset.

### `OTP_PEPPER`
- **Type:** String
- **Default:** `default-pepper-change-me`
- **Description:** Secret pepper for OTP hashing
- **Example:** `my-super-secret-pepper-string-2024`
- **Security:** ⚠️ MUST change in production! Keep secret!
- **Length:** Minimum 32 characters recommended

### `OTP_TTL_MINUTES`
- **Type:** Integer
- **Default:** `5`
- **Description:** OTP expiration time in minutes
- **Examples:** `5`, `10`, `15`
- **Recommended:** `5` to `10` minutes

### `OTP_MAX_ATTEMPTS`
- **Type:** Integer
- **Default:** `5`
- **Description:** Maximum OTP verification attempts
- **Examples:** `3`, `5`, `10`
- **Recommended:** `3` to `5` attempts

---

## SMTP Configuration

Email server settings for sending OTP and notifications.

### `SMTP_HOST`
- **Type:** String
- **Default:** `smtp.gmail.com`
- **Description:** SMTP server hostname
- **Examples:**
  - Gmail: `smtp.gmail.com`
  - SendGrid: `smtp.sendgrid.net`
  - Mailgun: `smtp.mailgun.org`
  - Custom: `mail.yourdomain.com`

### `SMTP_PORT`
- **Type:** Integer
- **Default:** `587`
- **Description:** SMTP server port
- **Common Ports:**
  - `587` - TLS (recommended)
  - `465` - SSL
  - `25` - Unencrypted (not recommended)

### `SMTP_USERNAME`
- **Type:** String
- **Default:** `` (empty)
- **Description:** SMTP authentication username
- **Example:** `your-email@gmail.com`, `apikey` (SendGrid)

### `SMTP_PASSWORD`
- **Type:** String
- **Default:** `` (empty)
- **Description:** SMTP authentication password or API key
- **Example:** Gmail app password, SendGrid API key
- **Security:** ⚠️ Keep this secret!
- **Gmail:** Use [App Passwords](https://support.google.com/accounts/answer/185833)

### `SMTP_FROM_EMAIL`
- **Type:** String (Email)
- **Default:** `` (empty)
- **Description:** Sender email address
- **Example:** `noreply@yourdomain.com`
- **Note:** Must be authorized to send from SMTP server

### `SMTP_FROM_NAME`
- **Type:** String
- **Default:** `Journaling App`
- **Description:** Sender name displayed in emails
- **Example:** `My Journaling App`, `YourCompany`

---

## Server Configuration

HTTP server settings.

### `SERVER_PORT`
- **Type:** String
- **Default:** `8080`
- **Description:** HTTP server port
- **Examples:** `8080`, `3000`, `80`, `443`
- **Note:** Ports below 1024 require root privileges

### `SERVER_HOST`
- **Type:** String
- **Default:** `0.0.0.0`
- **Description:** Server bind address
- **Options:**
  - `0.0.0.0` - Listen on all interfaces (recommended)
  - `localhost` - Listen only on localhost
  - `192.168.1.100` - Listen on specific IP

### `ENVIRONMENT`
- **Type:** String
- **Default:** `development`
- **Options:** `development`, `staging`, `production`
- **Description:** Application environment
- **Effects:**
  - Logging verbosity
  - Error detail in responses
  - Debug features

---

## 📋 Complete .env Example

```env
# PostgreSQL Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_secure_password
DB_NAME=journaling_auth
DB_SSLMODE=disable

# MongoDB
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=journaling

# JWT Keys
JWT_PRIVATE_KEY_PATH=./keys/jwt_private.pem
JWT_PUBLIC_KEY_PATH=./keys/jwt_public.pem
JWT_ACCESS_TOKEN_DURATION=15m
JWT_REFRESH_TOKEN_DURATION=168h

# OTP Settings
OTP_PEPPER=your-super-secret-pepper-change-this-in-production
OTP_TTL_MINUTES=5
OTP_MAX_ATTEMPTS=5

# SMTP Email
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM_EMAIL=noreply@yourdomain.com
SMTP_FROM_NAME=Journaling App

# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
ENVIRONMENT=development
```

---

## 🔧 Environment-Specific Configurations

### Development
```env
DB_HOST=localhost
MONGO_URI=mongodb://localhost:27017
ENVIRONMENT=development
JWT_ACCESS_TOKEN_DURATION=1h
DB_SSLMODE=disable
```

### Staging
```env
DB_HOST=staging-db.example.com
MONGO_URI=mongodb://staging-mongo:27017
ENVIRONMENT=staging
JWT_ACCESS_TOKEN_DURATION=30m
DB_SSLMODE=require
```

### Production
```env
DB_HOST=prod-db.example.com
MONGO_URI=mongodb+srv://user:pass@cluster.mongodb.net
ENVIRONMENT=production
JWT_ACCESS_TOKEN_DURATION=15m
DB_SSLMODE=verify-full
OTP_TTL_MINUTES=5
```

---

## 🔐 Security Best Practices

### 1. Never Commit Secrets
```bash
# Add to .gitignore
.env
.env.local
.env.production
*.pem
```

### 2. Use Strong Values
- **Passwords:** Minimum 16 characters, mix of letters, numbers, symbols
- **Pepper:** Minimum 32 characters, cryptographically random
- **Keys:** Use proper RSA key generation (2048-bit minimum)

### 3. Rotate Regularly
- JWT keys: Every 90 days
- Database passwords: Every 180 days
- OTP pepper: Every year
- SMTP credentials: As recommended by provider

### 4. Use Environment-Specific Values
- Different credentials per environment
- Never use production credentials in development
- Separate MongoDB databases per environment

### 5. Restrict Access
- Limit who can access .env files
- Use secret management services (AWS Secrets Manager, HashiCorp Vault)
- Set proper file permissions: `chmod 600 .env`

---

## 🚀 Docker Deployment

### Using Environment File
```bash
docker run -d \
  --name backend-journaling \
  --env-file .env.production \
  -p 8080:8080 \
  backend-journaling:latest
```

### Using Individual Variables
```bash
docker run -d \
  --name backend-journaling \
  -e DB_HOST=db \
  -e MONGO_URI=mongodb://mongo:27017 \
  -e JWT_PRIVATE_KEY_PATH=/secrets/private.pem \
  -p 8080:8080 \
  backend-journaling:latest
```

---

## 🧪 Testing Configuration

### Unit Tests
```env
DB_HOST=localhost
DB_NAME=journaling_test
MONGO_DATABASE=journaling_test
ENVIRONMENT=test
```

### Integration Tests
```env
DB_HOST=testdb
MONGO_URI=mongodb://testmongo:27017
SMTP_HOST=mailhog
SMTP_PORT=1025
```

---

## ❗ Common Issues

### Issue: "Failed to connect to database"
**Check:**
- DB_HOST is correct
- DB_PORT is correct
- PostgreSQL is running
- Network connectivity

### Issue: "Failed to connect to MongoDB"
**Check:**
- MONGO_URI format is correct
- MongoDB is running
- Network connectivity
- Authentication credentials

### Issue: "Failed to read private key"
**Check:**
- JWT_PRIVATE_KEY_PATH exists
- File permissions allow reading
- File is valid PEM format

### Issue: "Failed to send email"
**Check:**
- SMTP credentials are correct
- SMTP_HOST and SMTP_PORT are correct
- Less secure app access enabled (Gmail)
- Firewall not blocking SMTP port

---

## 📚 Additional Resources

- [Go Environment Variables](https://pkg.go.dev/os#Getenv)
- [MongoDB Connection String](https://docs.mongodb.com/manual/reference/connection-string/)
- [PostgreSQL Connection Strings](https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)
- [SMTP Configuration Guide](https://nodemailer.com/smtp/)

---

**Last Updated:** October 28, 2025
