# Mini-Zanzibar Authorization System

A simplified implementation of Google's Zanzibar global authorization system.

## Features

- Flexible configuration language for access control policies
- ACL storage and evaluation using relational tuples
- Consistent and scalable authorization decisions
- Low latency and high availability for authorization checks

## Architecture

- **LevelDB**: For storing ACL tuples (object#relation@user format)
- **ConsulDB**: For namespace configuration with versioning
- **REST API**: HTTP endpoints for ACL and namespace management

## API Endpoints

- `POST /acl` - Create/Update ACL
- `GET /acl/check` - Check authorization
- `POST /namespace` - Create/Update namespace configuration

## Running the Application

```bash
# Install dependencies
go mod tidy

# Run the server
go run cmd/server/main.go
```

## Environment Variables

Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
```

## Development

This project follows secure development practices and includes:
- Threat modeling documentation
- Security requirements analysis
- Static code analysis integration
- Secure code review process