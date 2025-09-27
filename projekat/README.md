# Mini Zanzibar Authorization System

A simplified implementation of Google's Zanzibar authorization system in Go.

## Getting Started

### Prerequisites
- Go 1.21 or higher
- Git

### Installation

1. Navigate to the project directory:
```bash
cd projekat
```

2. Download dependencies:
```bash
go mod tidy
```

3. Run the server:
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

### API Endpoints

#### 1. Write Relations
**POST** `/write`

Write one or more relations to the system.

```json
{
  "relations": [
    {
      "object": {"type": "document", "id": "doc1"},
      "relation": "viewer",
      "subject": {"type": "user", "id": "alice"}
    }
  ]
}
```

#### 2. Check Permissions
**POST** `/check`

Check if a subject has a specific relation to an object.

```json
{
  "object": {"type": "document", "id": "doc1"},
  "relation": "viewer",
  "subject": {"type": "user", "id": "alice"}
}
```

Response:
```json
{
  "allowed": true
}
```

#### 3. Read Relations
**POST** `/read`

Read all relations for a specific object.

```json
{
  "type": "document",
  "id": "doc1"
}
```

#### 4. Health Check
**GET** `/health`

Check if the service is running.

### Example Usage

```bash
# Write a relation
curl -X POST http://localhost:8080/write \
  -H "Content-Type: application/json" \
  -d '{
    "relations": [
      {
        "object": {"type": "document", "id": "doc1"},
        "relation": "viewer",
        "subject": {"type": "user", "id": "alice"}
      }
    ]
  }'

# Check permission
curl -X POST http://localhost:8080/check \
  -H "Content-Type: application/json" \
  -d '{
    "object": {"type": "document", "id": "doc1"},
    "relation": "viewer",
    "subject": {"type": "user", "id": "alice"}
  }'

# Read relations
curl -X POST http://localhost:8080/read \
  -H "Content-Type: application/json" \
  -d '{
    "type": "document",
    "id": "doc1"
  }'
```

### Project Structure

```
projekat/
├── cmd/server/          # Application entry point
├── internal/
│   ├── api/            # HTTP handlers
│   ├── auth/           # Authorization logic
│   └── storage/        # Data storage interfaces and implementations
├── pkg/models/         # Shared data types
├── tests/              # Integration tests
├── go.mod              # Go module definition
└── README.md           # This file
```

### Development

```bash
# Run tests
go test ./...

# Build binary
go build -o bin/mini-zanzibar cmd/server/main.go

# Run binary
./bin/mini-zanzibar
```

## Features

- ✅ Basic relation storage (in-memory)
- ✅ Permission checking
- ✅ REST API
- ✅ CORS support
- 🔄 Transitive relations (planned)
- 🔄 Persistent storage (planned)
- 🔄 Watch API (planned)