## Tiger Backend API Project

**Limitations or Future Scope:**

- Basic data validation/sanitization
- Basic authentication (access token, single route login)
  - Future exploration of full OAuth2 & 2FA for enhanced security
- Custom logger Implimentation
- ☁️ Image upload to cloud/NAS storage

## Live APIs and Documentation

Visit the live API documentation here: [https://tiger-backend-api-app.onrender.com/api/v1/docs/index.html](https://tiger-backend-api-app.onrender.com/api/v1/docs/index.html)

### Prerequisites to RUN:

| Requirement                   | Version  | Icon |
| ----------------------------- | -------- | ---- |
| Golang                        | ~v1.21   |      |
| Postgres DB                   | ~v15     |      |
| Brevo Email Platform API Keys | Required |      |

## How to RUN:

1. Clone the repository to your system.
2. Fix all dependencies: `go mod tidy`
3. Create a `.env` file in the root directory (see sample below).
4. Start the server: `go run ./cmd/`
5. Access API Docs: http://localhost:3000/api/v1/docs/index.html

**Sample .env:**

```
# Database Connection String
DB_STR="postgres://user:password@localhost/db_name"

# Server Configuration
SERVER_PORT="3000"
ALLOWED_ORIGINS="http://localhost:3000"

# Authentication
TOKEN_SECRET="Secret for signing JWT Token"

# Brevo Email Configuration (DO NOT change service name)
EMAIL_SERVICE="brevo"
EMAIL_FROM_ADDRESS="your-email@example.com"
EMAIL_API_KEY="Your Brevo API Key"
EMAIL_API_ENDPOINT="https://api.brevo.com/v3/smtp/email"

# AWS Configuration (currently unused)
AWS_REGION="test"
AWS_ACCESS_KEY_ID="test"
AWS_SECRET_ACCESS_KEY="test"

# Image Storage Path
IMAGE_STORAGE_PATH="/tiger_images/"
```

**Generate Own Documentation:**

1. Install `swag`: `go install github.com/swaggo/swag/cmd/swag@latest`
2. Generate documents: `swag init -g api/router.go`
3. Start the server and access docs: http://localhost:3000/api/v1/docs/index.html
