# Quick Start Guide - Mini-Zanzibar

## Prerequisites Check

Before starting, ensure you have:
- ✅ Go 1.21+ installed (`go version`)
- ✅ Consul installed (`consul version`) or Docker
- ✅ curl or PowerShell for testing

## Quick Test

### 1. Start the System

**Terminal 1 - Start Consul:**
```bash
# Option A: Using Consul directly
consul agent -dev -client=0.0.0.0

# Option B: Using Docker
docker run -d --name consul -p 8500:8500 hashicorp/consul:latest
```

**Terminal 2 - Start Mini-Zanzibar:**
```bash
cd mini-zanzibar
go mod tidy
go run cmd/server/main.go
```

### 2. Run Automated Tests

**Option A: PowerShell (Windows):**
```powershell
cd scripts
.\test-flow.ps1
```

**Option B: Bash (Linux/Mac):**
```bash
cd scripts
chmod +x test-flow.sh
./test-flow.sh
```

### 3. Manual Testing

**Test Health:**
```bash
curl http://localhost:8080/health
```

**Create Namespace:**
```bash
curl -X POST http://localhost:8080/namespace \
  -H "Content-Type: application/json" \
  -d '{
    "namespace": "doc",
    "relations": {
      "owner": {},
      "editor": {
        "union": [
          {"this": {}},
          {"computed_userset": {"relation": "owner"}}
        ]
      },
      "viewer": {
        "union": [
          {"this": {}},
          {"computed_userset": {"relation": "editor"}}
        ]
      }
    }
  }'
```

**Create ACL:**
```bash
curl -X POST http://localhost:8080/acl \
  -H "Content-Type: application/json" \
  -d '{
    "object": "doc:readme",
    "relation": "owner",
    "user": "user:alice"
  }'
```

**Check Authorization:**
```bash
curl "http://localhost:8080/acl/check?object=doc:readme&relation=owner&user=user:alice"
# Should return: {"authorized": true}
```

## What's Working vs TODO

### ✅ Currently Working:
- Health check endpoint
- Namespace creation and retrieval with versioning
- ACL tuple storage and retrieval
- Basic authorization checks (direct tuples only)
- Database persistence (LevelDB + Consul)
- All CRUD operations for ACLs and namespaces

### 🔧 TODO (Marked in Code):
- **Computed usersets evaluation** - Currently only direct tuple checks work
- **Union operations in authorization logic** - Relations like "editor includes owner" not fully implemented
- **Authentication/Authorization middleware** - All endpoints are currently public
- **Input validation and sanitization** - Basic gin validation only
- **Rate limiting** - Structure exists but not enforced
- **Comprehensive logging** - Basic logging only
- **TLS/HTTPS configuration**
- **Caching for performance**

### 🧪 Expected Test Results:

With current implementation:
- ✅ Direct ACL checks (user explicitly granted a relation)
- ❌ Computed userset checks (e.g., owner should have viewer access)
- ✅ Database operations (storage, retrieval, deletion)
- ✅ Namespace management with versioning
- ✅ Error handling for invalid requests

## Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐
│   REST API      │    │   Consul DB     │
│  (Port 8080)    │◄──►│ (Namespaces)    │
└─────────────────┘    └─────────────────┘
         │
         ▼
┌─────────────────┐
│   LevelDB       │
│  (ACL Tuples)   │
└─────────────────┘
```

## Next Steps

1. **Review the TODO items** in the codebase
2. **Implement computed usersets** for full Zanzibar functionality
3. **Add security features** (authentication, input validation)
4. **Write comprehensive tests**
5. **Set up CI/CD pipeline**

## Troubleshooting

**"Connection refused" errors:**
- Check if Consul is running on port 8500
- Check if Mini-Zanzibar is running on port 8080

**"Permission denied" errors:**
- Ensure write permissions to `./data/leveldb/` directory

**"Module not found" errors:**
- Run `go mod tidy` to install dependencies

**Tests failing:**
- Remember: Computed usersets are not implemented yet
- Direct tuple tests should pass
- Check server logs for detailed error messages