# Mini-Zanzibar Demo with Authentication Middleware

## Architecture

This demo implements a proper authentication/authorization architecture with three layers:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Client    │───▶│  Auth Service   │───▶│ Mini-Zanzibar   │
│  (Frontend UI)  │    │ (Authentication)│    │ (Authorization) │
│   Port 3000     │    │   Port 8081     │    │   Port 8080     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Components

1. **Web Client** (Port 3000)
   - HTML/CSS/JavaScript frontend
   - Login/logout functionality
   - Document management interface
   - Authorization testing tools

2. **Auth Service** (Port 8081)
   - User authentication (login/logout)
   - Session management with cookies
   - JWT token support
   - Request proxying to Mini-Zanzibar
   - Role-based access control

3. **Mini-Zanzibar** (Port 8080)
   - Authorization decisions
   - ACL management
   - Namespace configuration
   - Hierarchical permission evaluation

## Quick Start

### Option 1: Use the startup script
```powershell
# Make sure Mini-Zanzibar is running first
cd mini-zanzibar
go run cmd/server/main.go

# In another terminal, start the demo
.\start-demo.ps1
```

### Option 2: Manual startup
```powershell
# Terminal 1: Mini-Zanzibar
cd mini-zanzibar
go run cmd/server/main.go

# Terminal 2: Auth Service  
cd auth-service
go run main.go

# Terminal 3: Web Client
cd web-client
python -m http.server 3000
# OR with Node.js: npx http-server -p 3000
```

Then open http://localhost:3000 in your browser.

## Demo Users

| Username | Password   | Role   | Permissions |
|----------|------------|--------|-------------|
| alice    | alice123   | admin  | Full access to all documents and ACL management |
| bob      | bob123     | editor | Can edit specific documents |
| charlie  | charlie123 | viewer | Can view specific documents only |
| david    | david123   | user   | Limited access |

## Features Demonstrated

### Authentication Flow
1. User logs in via web interface
2. Auth service validates credentials
3. Session cookie is created
4. Subsequent requests include session for authentication

### Authorization Flow  
1. User attempts document action
2. Auth service forwards request to Mini-Zanzibar with user context
3. Mini-Zanzibar checks ACL rules
4. Decision is returned and enforced

### Key Capabilities
- **Secure Authentication**: Session-based with HTTP-only cookies
- **Role-Based Access**: Different capabilities per user role
- **ACL Management**: Admin users can create/modify access rules
- **Document Access Control**: Fine-grained permissions per document
- **Authorization Testing**: Built-in tools to test permission scenarios

## API Endpoints

### Auth Service
- `POST /auth/login` - User authentication
- `POST /auth/logout` - User logout
- `GET /auth/me` - Current user info
- `GET /documents` - List documents by role
- `POST /documents/:id/access` - Check document access
- `GET /api/acl/check` - Authorization check (proxied)
- `POST /api/acl` - Create ACL (admin only, proxied)

### Usage Examples

#### Login
```bash
curl -X POST http://localhost:8081/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "alice", "password": "alice123"}' \
  -c cookies.txt
```

#### Check Authorization
```bash
curl -X GET "http://localhost:8081/api/acl/check?object=doc:report1&relation=viewer&user=user:alice" \
  -b cookies.txt
```

#### Create ACL (Admin Only)
```bash
curl -X POST http://localhost:8081/api/acl \
  -H "Content-Type: application/json" \
  -d '{"object": "doc:report1", "relation": "editor", "user": "user:bob"}' \
  -b cookies.txt
```

## Security Features

- **Session Management**: Secure HTTP-only cookies
- **CORS Protection**: Configured for local development
- **Role-Based Access**: Admin/Editor/Viewer/User roles
- **Request Proxying**: User context injection
- **Authentication Middleware**: All API calls require valid session

## Development Notes

- Auth service acts as a secure proxy to Mini-Zanzibar
- Sessions expire after 1 hour of inactivity
- All document operations require proper authorization
- ACL management is restricted to admin users only
- The system demonstrates separation of authentication and authorization concerns