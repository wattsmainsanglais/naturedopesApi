# NatureDopes API

## Project Overview
NatureDopes API is a Go-based REST API backend for cataloging locations of wild flora (plants). It serves the NatureDopes 2.0 Next.js application (located at `/home/andrew/Code/2025/natureDopes2.0`), providing secure access to plant observation data with GPS coordinates.

## Tech Stack
- **Language**: Go 1.24.5
- **Web Framework**: Gorilla Mux (routing)
- **Database**: PostgreSQL with pgx/v4 driver
- **Image Storage**: Iagon (decentralized storage)
- **Deployment**: Railway
- **Schema Management**: Prisma (schema definition)

## Architecture

### Project Structure
```
naturedopesApi/
├── main.go              # Entry point, DB connection, route setup
├── routes.go            # Route handlers for images and API keys
├── endpoints/
│   ├── image.go         # Image data retrieval logic
│   └── apikey.go        # API key generation, validation, management
├── middleware/
│   └── apikey.go        # API key authentication middleware
└── prisma/
    └── schema.prisma    # Database schema definition
```

### Database Schema
The API uses the following main tables:
- **images**: Stores flora observations with species name, GPS coordinates (lat/long), image path, and user ID
- **users**: User accounts with username, email, password, and auth tokens
- **api_keys**: API key management with key, name, created_at, expires_at, last_used, revoked status, and created_ip
- **passResetToken**: Password reset token management

## Current Features

### Image Endpoints (Protected by API Key)
All image endpoints require an `X-API-Key` header with a valid API key.

- `GET /images` - Retrieve all flora observations
- `GET /images/{id}` - Retrieve a specific observation by ID

### API Key Management Endpoints (Unprotected)
Currently unprotected to allow frontend integration. Will be protected once user authentication is integrated.

- `POST /api/keys` - Generate a new API key
  - Body: `{"name": "key-name"}`
  - Returns: API key object with 64-character hex key

- `GET /api/keys` - List all API keys (ordered by created_at DESC)
  - Returns: Array of API key objects with usage metadata

- `DELETE /api/keys/{id}` - Revoke an API key
  - Returns: 204 No Content on success

### Authentication System
- API key middleware validates the `X-API-Key` header on protected routes
- Keys are stored as 64-character hex strings (32 random bytes)
- Keys expire after 90 days from creation
- Tracks `last_used` timestamp on each valid request
- Logs IP address on key creation for abuse tracking
- Supports key revocation without deletion
- Validation checks: exists, not revoked, not expired

## API Integration with Frontend

The Next.js app (natureDopes2.0) will integrate with this API:
- Users authenticate via NextAuth on the frontend
- A new page will be added allowing authenticated users to generate API keys
- API keys enable programmatic access to flora observation data
- Frontend uses the same PostgreSQL database (via Prisma client)

## Development Setup

### Prerequisites
- Go 1.24.5+
- PostgreSQL database
- Environment variable: `DATABASE_URL` (PostgreSQL connection string)

### Running Locally
```bash
# Install dependencies
go mod download

# Set environment variable
export DATABASE_URL="postgresql://user:password@host:port/dbname"

# Run the server
go run .
```

The server will start on port 8080.

### Database Migrations
Use Prisma for schema management:
```bash
cd prisma
npx prisma migrate dev
npx prisma generate
```

## API Usage Examples

### Generate an API Key
```bash
curl -X POST http://localhost:8080/api/keys \
  -H "Content-Type: application/json" \
  -d '{"name": "My Flora Research Key"}'
```
Response includes the key (save it!), expiration date (90 days), and creation timestamp.

### Retrieve All Flora Observations
```bash
curl http://localhost:8080/images \
  -H "X-API-Key: your-64-character-hex-key"
```
Rate limit: 100 requests/hour per key, 1000 requests/day per IP.

### Get Specific Observation
```bash
curl http://localhost:8080/images/123 \
  -H "X-API-Key: your-64-character-hex-key"
```

### Rate Limit Headers
When rate limited, you'll receive:
- HTTP 429 Too Many Requests
- Error message indicating which limit was exceeded

## Planned Features

### Short-term
1. Add page to Next.js app for API key generation (in progress)
2. Protect API key management endpoints with user authentication
3. Associate API keys with user IDs for access control
4. Add image upload endpoint to store flora photos on Iagon

### Medium-term
1. User authentication endpoints (register, login, profile)
2. Filtering and search capabilities (by species, location, date)
3. Geographic queries (find flora within radius of coordinates)
4. Rate limiting and usage quotas per API key
5. Image upload and metadata extraction

### Long-term
1. Species identification API integration
2. Community contributions and verification system
3. Public vs. private observation settings
4. Export functionality (GeoJSON, CSV)

## Security Measures

### Implemented (Production Ready)
✅ **Rate Limiting**
  - Per-IP limit: 1000 requests/day
  - Per-key limit: 100 requests/hour
  - Automatic cleanup of old rate limit entries
  - Returns 429 Too Many Requests when exceeded

✅ **API Key Security**
  - 64-character cryptographically random hex keys
  - 90-day automatic expiration
  - IP address logging on creation
  - Revocation support
  - Validation checks: exists, not revoked, not expired

✅ **CORS Configuration**
  - Configured for open access (all origins)
  - Allows: GET, POST, DELETE, OPTIONS
  - Permits Content-Type and X-API-Key headers

✅ **Request Tracking**
  - Logs IP addresses for abuse investigation
  - Tracks last_used timestamp per key
  - Created_at timestamp for all keys

### Future Enhancements
1. Hash API keys before storage (store hash, return key only on creation)
2. Add CAPTCHA to key generation endpoint
3. Protect API key management endpoints with user authentication
4. Add API key scopes/permissions system
5. Implement structured logging and monitoring
6. Add usage analytics dashboard
7. Webhook alerts for abuse patterns

## Environment Variables
- `DATABASE_URL`: PostgreSQL connection string (required)
  - Format: `postgresql://username:password@host:port/database?sslmode=require`

## Related Projects
- **Frontend**: `/home/andrew/Code/2025/natureDopes2.0` (Next.js app with mapping, games, multi-language support)
- Uses shared PostgreSQL database
- Frontend handles user authentication via NextAuth
- API provides programmatic data access for external tools and advanced users

## Recent Changes
- **Security Hardening (2025-10-27)**:
  - Added rate limiting: 100 req/hour per key, 1000 req/day per IP
  - Implemented 90-day key expiration
  - Added IP address logging on key creation
  - Configured CORS for open access
  - Updated validation to check key expiration
- **Initial API Key System**:
  - Migrated database schema to add `api_keys` table
  - Added API key generation and management endpoints
  - Implemented middleware for API key validation
  - Updated route handlers to use API key authentication
  - Set up route organization with `SetupRoutes` function

## Notes
- The `firstApi` binary in the project root is an older compiled version
- `.env` file exists but godotenv is commented out for Railway deployment
- Prisma schema uses JavaScript client generator but API is written in Go
- API uses raw SQL queries via pgx instead of Prisma client
