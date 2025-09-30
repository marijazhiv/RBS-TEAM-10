# Auth Service

Authentication middleware service that sits between the web client and Mini-Zanzibar.

## Architecture

```
Web Client (Port 3000) 
    ↓
Auth Service (Port 8081) 
    ↓
Mini-Zanzibar (Port 8080)
```

## Features

- **User Authentication**: Login/logout with session management
- **JWT Tokens**: Secure token-based authentication
- **Authorization Proxy**: Forwards authorization requests to Mini-Zanzibar with user context
- **Role-Based Access**: Admin/Editor/Viewer/User roles
- **CORS Support**: Cross-origin requests from web client

## Demo Users

| Username | Password   | Role   | Access Level |
|----------|------------|--------|--------------|
| alice    | alice123   | admin  | Full access  |
| bob      | bob123     | editor | Edit docs    |
| charlie  | charlie123 | viewer | View only    |
| david    | david123   | user   | Limited      |

## API Endpoints

### Authentication
- `POST /auth/login` - User login
- `POST /auth/logout` - User logout  
- `GET /auth/me` - Get current user info

### Protected API (requires authentication)
- `GET /api/acl/check` - Check authorization (proxied to Zanzibar)
- `POST /api/acl` - Create ACL (admin only, proxied to Zanzibar)
- `DELETE /api/acl` - Delete ACL (admin only, proxied to Zanzibar)
- `POST /api/namespace` - Create namespace (admin only, proxied to Zanzibar)
- `GET /api/namespace/:id` - Get namespace (proxied to Zanzibar)

### Documents
- `GET /documents` - List available documents based on user role
- `POST /documents/:id/access` - Check document access

## Usage

1. **Start the service:**
   ```bash
   go run main.go
   ```

2. **Login:**
   ```bash
   curl -X POST http://localhost:8081/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username": "alice", "password": "alice123"}'
   ```

3. **Check authorization:**
   ```bash
   curl -X GET "http://localhost:8081/api/acl/check?object=doc:report1&relation=viewer&user=user:alice" \
     -H "Cookie: auth-session=your-session-cookie"
   ```

## Session Management

- Sessions are stored in HTTP-only cookies
- Session expires after 1 hour of inactivity
- JWT tokens are provided for stateless authentication (optional)

## Security Features

- Password hashing with bcrypt
- Secure session cookies
- CORS protection
- Role-based access control
- Request proxying with user context injection

## Development

The service automatically forwards authenticated requests to Mini-Zanzibar at `http://localhost:8080` with the user context added in headers.

Make sure Mini-Zanzibar is running before starting the auth service.