# Complete Testing Guide for Mini-Zanzibar

This guide will walk you through testing the entire Mini-Zanzibar authorization system flow.

## Prerequisites

1. **Go 1.21+** installed
2. **Consul** installed and running
3. **curl** or **Postman** for API testing

## Step 1: Environment Setup

### 1.1 Install Dependencies
```bash
cd mini-zanzibar
go mod tidy
```

### 1.2 Setup Environment
```bash
# Copy environment configuration
cp .env.example .env
```

### 1.3 Start Consul
```bash
# Option 1: Using Consul directly
consul agent -dev -client=0.0.0.0

# Option 2: Using Docker
docker run -d --name consul -p 8500:8500 consul:latest
```

## Step 2: Start Mini-Zanzibar Server

```bash
# From the mini-zanzibar directory
go run cmd/server/main.go
```

You should see output like:
```
INFO  Starting Mini-Zanzibar server  host=localhost port=8080
```

## Step 3: Test Health Check

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "timestamp": "2025-09-28T10:00:00Z",
  "service": "mini-zanzibar",
  "version": "1.0.0"
}
```

## Step 4: Test Complete Authorization Flow

### 4.1 Create a Namespace Configuration

First, define the authorization rules:

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

Expected response:
```json
{
  "message": "Namespace created successfully",
  "namespace": "doc",
  "version": 1
}
```

### 4.2 Verify Namespace Creation

```bash
curl http://localhost:8080/namespace/doc
```

Expected response:
```json
{
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
  },
  "version": 1
}
```

### 4.3 Create ACL Tuples

Create some authorization relationships:

```bash
# Alice is owner of readme document
curl -X POST http://localhost:8080/acl \
  -H "Content-Type: application/json" \
  -d '{
    "object": "doc:readme",
    "relation": "owner",
    "user": "user:alice"
  }'

# Bob is editor of readme document
curl -X POST http://localhost:8080/acl \
  -H "Content-Type: application/json" \
  -d '{
    "object": "doc:readme",
    "relation": "editor",
    "user": "user:bob"
  }'

# Charlie is viewer of readme document
curl -X POST http://localhost:8080/acl \
  -H "Content-Type: application/json" \
  -d '{
    "object": "doc:readme",
    "relation": "viewer",
    "user": "user:charlie"
  }'
```

Each should return:
```json
{
  "message": "ACL created successfully"
}
```

### 4.4 Test Authorization Checks

Now test various authorization scenarios:

```bash
# Test 1: Alice (owner) should have viewer access
curl "http://localhost:8080/acl/check?object=doc:readme&relation=viewer&user=user:alice"

# Test 2: Bob (editor) should have viewer access
curl "http://localhost:8080/acl/check?object=doc:readme&relation=viewer&user=user:bob"

# Test 3: Charlie (viewer) should have viewer access
curl "http://localhost:8080/acl/check?object=doc:readme&relation=viewer&user=user:charlie"

# Test 4: Alice (owner) should have editor access
curl "http://localhost:8080/acl/check?object=doc:readme&relation=editor&user=user:alice"

# Test 5: Bob (editor) should have editor access
curl "http://localhost:8080/acl/check?object=doc:readme&relation=editor&user=user:bob"

# Test 6: Charlie (viewer) should NOT have editor access
curl "http://localhost:8080/acl/check?object=doc:readme&relation=editor&user=user:charlie"

# Test 7: Alice (owner) should have owner access
curl "http://localhost:8080/acl/check?object=doc:readme&relation=owner&user=user:alice"

# Test 8: Bob should NOT have owner access
curl "http://localhost:8080/acl/check?object=doc:readme&relation=owner&user=user:bob"

# Test 9: Unknown user should have no access
curl "http://localhost:8080/acl/check?object=doc:readme&relation=viewer&user=user:unknown"
```

**Expected Results:**
- Tests 1-5, 7: `{"authorized": true}`
- Tests 6, 8, 9: `{"authorized": false}`

### 4.5 List ACLs

```bash
# List all ACLs for the readme document
curl http://localhost:8080/acl/object/doc:readme

# List all ACLs for Alice
curl http://localhost:8080/acl/user/user:alice
```

### 4.6 Test ACL Deletion

```bash
# Remove Bob's editor access
curl -X DELETE http://localhost:8080/acl \
  -H "Content-Type: application/json" \
  -d '{
    "object": "doc:readme",
    "relation": "editor",
    "user": "user:bob"
  }'

# Verify Bob no longer has editor access
curl "http://localhost:8080/acl/check?object=doc:readme&relation=editor&user=user:bob"
# Should return: {"authorized": false}

# But Bob should still have viewer access (if there's a direct viewer tuple)
curl "http://localhost:8080/acl/check?object=doc:readme&relation=viewer&user=user:bob"
```

## Step 5: Advanced Testing Scenarios

### 5.1 Multiple Namespaces

```bash
# Create a file namespace
curl -X POST http://localhost:8080/namespace \
  -H "Content-Type: application/json" \
  -d '{
    "namespace": "file",
    "relations": {
      "owner": {},
      "reader": {}
    }
  }'

# Create file ACLs
curl -X POST http://localhost:8080/acl \
  -H "Content-Type: application/json" \
  -d '{
    "object": "file:config.txt",
    "relation": "owner",
    "user": "user:admin"
  }'

# Test file access
curl "http://localhost:8080/acl/check?object=file:config.txt&relation=owner&user=user:admin"
```

### 5.2 Namespace Versioning

```bash
# Update the doc namespace
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
      },
      "commenter": {}
    }
  }'

# Check version increased
curl http://localhost:8080/namespace/doc
# Should show version: 2

# Access previous version
curl http://localhost:8080/namespace/doc/version/1
```

### 5.3 List All Namespaces

```bash
curl http://localhost:8080/namespaces
```

Expected response:
```json
{
  "namespaces": ["doc", "file"]
}
```

## Step 6: Error Testing

### 6.1 Invalid Requests

```bash
# Invalid ACL creation (missing fields)
curl -X POST http://localhost:8080/acl \
  -H "Content-Type: application/json" \
  -d '{"object": "doc:readme"}'
# Should return 400 Bad Request

# Invalid namespace
curl -X POST http://localhost:8080/namespace \
  -H "Content-Type: application/json" \
  -d '{"invalid": "data"}'
# Should return 400 Bad Request

# Non-existent namespace
curl http://localhost:8080/namespace/nonexistent
# Should return 404 Not Found
```

## Step 7: Performance Testing

### 7.1 Load Testing with Multiple Requests

```bash
# Create a simple load test script
for i in {1..100}; do
  curl -s "http://localhost:8080/acl/check?object=doc:readme&relation=viewer&user=user:alice" &
done
wait
```

## Step 8: Database Verification

### 8.1 Check LevelDB Data

The LevelDB data is stored in `./data/leveldb/`. You can verify the data persists by:
1. Stopping the server
2. Restarting it
3. Running the same authorization checks

### 8.2 Check Consul Data

Open Consul UI at `http://localhost:8500` and navigate to Key/Value store. Look for keys under `zanzibar/namespaces/`.

## Step 9: Integration Testing with Code

Run the integration tests:

```bash
# Run all tests
go test ./test/integration/...

# Run with verbose output
go test -v ./test/integration/...
```

## Troubleshooting

### Common Issues:

1. **Server won't start**
   - Check if port 8080 is available
   - Verify Consul is running
   - Check environment variables

2. **Consul connection failed**
   - Ensure Consul is running on localhost:8500
   - Check Consul logs for errors

3. **LevelDB errors**
   - Ensure write permissions to `./data/leveldb/`
   - Check disk space

4. **API returns 500 errors**
   - Check server logs
   - Verify database connections

### Logs Location:
- Server logs: Console output
- Consul logs: Check Consul agent output
- LevelDB: No separate logs (embedded)

## Next Steps

After basic testing:
1. Implement the TODO items in the code
2. Add comprehensive unit tests
3. Set up CI/CD pipeline
4. Implement security features (authentication, TLS)
5. Add monitoring and metrics