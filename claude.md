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
- **api_keys**: API key management with key, name, created_at, last_used, and revoked status
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
- Tracks `last_used` timestamp on each valid request
- Supports key revocation without deletion

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

### Retrieve All Flora Observations
```bash
curl http://localhost:8080/images \
  -H "X-API-Key: your-64-character-hex-key"
```

### Get Specific Observation
```bash
curl http://localhost:8080/images/123 \
  -H "X-API-Key: your-64-character-hex-key"
```

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

## Security Considerations

### Current State
- API key management endpoints are unprotected (temporary for development)
- API keys are stored in plaintext (consider hashing for production)
- No rate limiting implemented yet
- CORS not configured

### Production Recommendations
1. Hash API keys before storage (store hash, return key only on creation)
2. Implement rate limiting middleware
3. Add CORS configuration for frontend domain
4. Protect API key endpoints with user session authentication
5. Add API key scopes/permissions system
6. Implement request logging and monitoring
7. Add input validation and sanitization

## Environment Variables
- `DATABASE_URL`: PostgreSQL connection string (required)
  - Format: `postgresql://username:password@host:port/database?sslmode=require`

## Related Projects
- **Frontend**: `/home/andrew/Code/2025/natureDopes2.0` (Next.js app with mapping, games, multi-language support)
- Uses shared PostgreSQL database
- Frontend handles user authentication via NextAuth
- API provides programmatic data access for external tools and advanced users

## Recent Changes
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
