# Mini-Zanzibar Authorization System

A complete implementation of Google's Zanzibar authorization model with web interface, featuring multi-user authentication, document management, and real-time permission testing.

## 🏗️ System Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Client    │    │   Auth Service  │    │ Mini-Zanzibar   │
│   (Port 3001)   │────│   (Port 8081)   │────│   (Port 8080)   │
│                 │    │                 │    │     Docker      │
│ • Document UI   │    │ • Authentication│    │ • ACL Engine    │
│ • ACL Manager   │    │ • Session Mgmt  │    │ • LevelDB       │
│ • Auth Testing  │    │ • API Proxy     │    │ • Consul Config │
└─────────────────┘    └─────────────────┘    │ • Redis Cache   │
                                              └─────────────────┘
```

## 🚀 Quick Start

### Prerequisites
- Node.js (for web client)
- Go (for auth service)
- Docker & Docker Compose (for Mini-Zanzibar)

### Starting the System

1. **Start Mini-Zanzibar (Docker)**:
   ```bash
   cd mini-zanzibar/deployments/docker
   docker-compose up -d
   ```

2. **Start Auth Service**:
   ```bash
   cd auth-service
   go run main.go
   ```

3. **Start Web Client**:
   ```bash
   cd web-client
   node server.js
   ```

4. **Access the System**:
   - Web Interface: http://localhost:3001
   - Auth Service: http://localhost:8081
   - Mini-Zanzibar: http://localhost:8080

## 👥 User Accounts

| Username | Password | Role | Default Permissions |
|----------|----------|------|-------------------|
| `alice` | `alice123` | Owner | Full access to all documents |
| `bob` | `bob123` | Editor | Limited access (configurable) |
| `charlie` | `charlie123` | Viewer | Read-only access (configurable) |

## 📋 Testing Guide

### Test 1: Basic Authentication & Document Access

1. **Login as Alice**:
   - Username: `alice`, Password: `alice123`
   - Should see: All 3 documents (document1, document2, document3)
   - Can: View, edit, save, and share documents

2. **Login as Bob** (new browser/incognito):
   - Username: `bob`, Password: `bob123`
   - Should see: Only documents with granted permissions
   - Can: Edit documents where he has editor access

3. **Login as Charlie** (new browser/incognito):
   - Username: `charlie`, Password: `charlie123`
   - Should see: Only documents with granted permissions
   - Can: Only view documents (read-only)

### Test 2: ACL Management

1. **As Alice**, navigate to **"Access Control"** tab
2. **Create permissions for Bob**:
   - Object: `doc:document2`
   - Relation: `editor`
   - User: `user:bob`
   - Click "Create ACL"

3. **Create permissions for Charlie**:
   - Object: `doc:document3`
   - Relation: `viewer`
   - User: `user:charlie`
   - Click "Create ACL"

4. **Verify**: Bob should now see document2, Charlie should see document3

### Test 3: Authorization Testing Panel

1. **Navigate to "Test Authorization" tab**
2. **Test specific permissions**:
   - User field auto-populates with current user
   - Select document and permission level
   - Click "Test Authorization"
   - Verify expected results

### Test 4: Document Management

1. **View documents**: Click on document names to open
2. **Edit content**: Use the text editor interface
3. **Save changes**: Click "Save" button
4. **Verify permissions**: Non-owners cannot edit certain documents

## 🎯 Features Implemented

### ✅ Completed Features

- **Multi-user authentication** with session management
- **Document CRUD operations** with authorization checks
- **Real-time ACL management** through web interface
- **Permission testing panel** for debugging authorization
- **Role-based access control** (Owner, Editor, Viewer)
- **Session-based API proxy** between services
- **Responsive web interface** with multiple tabs
- **Activity logging** for user actions
- **Docker containerization** for Mini-Zanzibar
- **Proper error handling** for most scenarios

### 🔧 Working Authorization Model

```
Alice (Owner)
├── document1: ✅ owner
├── document2: ✅ owner  
└── document3: ✅ owner

Bob (Editor)
└── document2: ✅ owner 
└── document1: ✅ editor 


Charlie (Viewer)  
└── document3: ✅ viewer 
└── document2: ✅ editor 
```

## ⚠️ Known Issues & Limitations

### 1. HTTP 500 Errors for False Permissions
**Issue**: When users check permissions they don't have, the system returns HTTP 500 errors instead of proper "false" responses.

**Affected**: Authorization testing panel shows errors instead of "Access Denied"

**Workaround**: Test only permissions that should return "true"

**Root Cause**: Error handling in auth service proxy needs refinement

### 2. Incomplete Permission Hierarchy
**Issue**: Permission inheritance not fully implemented (Owner should automatically have Editor+Viewer permissions)

**Expected Behavior**:
```
Owner → Should also have Editor + Viewer permissions
Editor → Should also have Viewer permissions  
Viewer → Only Viewer permissions
```

**Current Behavior**: Each permission must be explicitly granted

**Impact**: Need to create separate ACLs for each permission level

### 3. Session Cookie Handling Edge Cases
**Issue**: Some authorization checks may fail due to session cookie propagation between web client and auth service

**Symptoms**: Intermittent authorization failures, especially for Bob and Charlie

**Workaround**: Refresh browser or re-login if permissions seem incorrect

### 4. Computed Usersets Not Fully Active
**Issue**: Mini-Zanzibar namespace configuration includes computed usersets, but they're not being processed correctly

**Technical Details**: The `computed_userset` relations in namespace config exist but don't affect authorization results

**Impact**: Manual ACL creation required for all permission levels

## 🔧 Troubleshooting

### Service Status Check
```bash
# Check all services are running
curl http://localhost:3001        # Web Client
curl http://localhost:8081/health # Auth Service  
curl http://localhost:8080/health # Mini-Zanzibar
docker ps                        # Docker containers
```

### Manual ACL Creation
If web interface ACL creation fails, use direct API:
```bash
# Create ACL for Bob on document2
curl -X POST http://localhost:8080/api/v1/acl \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user:alice" \
  -d '{"object":"doc:document2","relation":"editor","user":"user:bob"}'
```

### Permission Testing
Direct authorization checks:
```bash
# Test Bob's permissions on document2
curl "http://localhost:8080/api/v1/acl/check?object=doc:document2&relation=editor&user=user:bob"
```

### Docker Issues
```bash
# Restart Mini-Zanzibar container
cd mini-zanzibar/deployments/docker
docker-compose down
docker-compose up -d

# Check logs
docker logs mini-zanzibar-app
```

## 📁 Project Structure

```
RBS-TEAM-10/
├── auth-service/           # Go-based authentication service
│   ├── main.go            # Auth service implementation
│   └── auth-service.exe   # Compiled binary
├── mini-zanzibar/         # Mini-Zanzibar authorization engine
│   ├── deployments/docker/   # Docker configuration
│   ├── internal/api/      # API handlers and middleware
│   └── server.exe         # Compiled binary
├── web-client/            # Frontend web application
│   ├── server.js          # Node.js backend server
│   ├── app.js             # Frontend JavaScript
│   ├── index.html         # Main web interface
│   ├── styles.css         # UI styling
│   └── documents/         # Sample documents
└── README.md              # This file
```

## 🚀 Future Improvements

### High Priority
1. **Fix HTTP 500 error handling** for false permission checks
2. **Implement permission hierarchy** properly in Mini-Zanzibar
3. **Namespace feature** properly set up on web-client
4. **Improve session cookie propagation** between services
5. **Add proper loading states** and error messages in UI
6. **Static testing**
7. **OWASP analysis**
9. **model treat**

## 📚 Technical Details

### Authorization Flow
1. User authenticates via auth service
2. Auth service issues session cookie
3. Web client makes requests with session
4. Auth service validates session and proxies to Mini-Zanzibar
5. Mini-Zanzibar checks ACLs and returns authorization result
6. Result propagated back through auth service to web client

### ACL Storage
- **Primary Storage**: LevelDB (in Docker container)
- **Configuration**: Consul KV store
- **Caching**: Redis for performance
- **Namespace**: `doc` for all document permissions

### Session Management
- **Type**: Server-side sessions with cookies
- **Storage**: In-memory (auth service)
- **Timeout**: Configurable (default: session-based)
- **Security**: HTTP-only cookies

