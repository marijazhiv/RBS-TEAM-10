# Mini-Zanzibar API Documentation

## Overview

The Mini-Zanzibar API provides endpoints for managing Access Control Lists (ACLs) and namespace configurations for authorization decisions.

Base URL: `http://localhost:8080`

## Authentication

**TODO**: Authentication is not yet implemented. All endpoints are currently public.

Future implementation will use JWT-based authentication:
```
Authorization: Bearer <jwt-token>
```

## Endpoints

### Health Check

#### GET /health
Check the health status of the service.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-09-28T10:00:00Z",
  "service": "mini-zanzibar",
  "version": "1.0.0"
}
```

### ACL Management

#### POST /acl
Create or update an ACL tuple.

**Request Body:**
```json
{
  "object": "doc:readme",
  "relation": "viewer",
  "user": "user:alice"
}
```

**Response:**
```json
{
  "message": "ACL created successfully"
}
```

#### GET /acl/check
Check if a user has a specific relation to an object.

**Query Parameters:**
- `object` (required): The object to check access for
- `relation` (required): The relation to check
- `user` (required): The user to check access for

**Example:**
```
GET /acl/check?object=doc:readme&relation=viewer&user=user:alice
```

**Response:**
```json
{
  "authorized": true
}
```

#### DELETE /acl
Delete an ACL tuple.

**Request Body:**
```json
{
  "object": "doc:readme",
  "relation": "viewer",
  "user": "user:alice"
}
```

**Response:**
```json
{
  "message": "ACL deleted successfully"
}
```

#### GET /acl/object/{object}
List all ACL tuples for a specific object.

**Response:**
```json
{
  "tuples": [
    {
      "object": "doc:readme",
      "relation": "viewer",
      "user": "user:alice"
    },
    {
      "object": "doc:readme",
      "relation": "editor",
      "user": "user:bob"
    }
  ]
}
```

#### GET /acl/user/{user}
List all ACL tuples for a specific user.

**Response:**
```json
{
  "tuples": [
    {
      "object": "doc:readme",
      "relation": "viewer",
      "user": "user:alice"
    },
    {
      "object": "doc:report",
      "relation": "owner",
      "user": "user:alice"
    }
  ]
}
```

### Namespace Management

#### POST /namespace
Create or update a namespace configuration.

**Request Body:**
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
  }
}
```

**Response:**
```json
{
  "message": "Namespace created successfully",
  "namespace": "doc",
  "version": 1
}
```

#### GET /namespace/{namespace}
Get the latest version of a namespace configuration.

**Response:**
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

#### GET /namespace/{namespace}/version/{version}
Get a specific version of a namespace configuration.

**Response:** Same as above, but for the specified version.

#### GET /namespaces
List all available namespaces.

**Response:**
```json
{
  "namespaces": ["doc", "file", "folder"]
}
```

#### DELETE /namespace/{namespace}
Delete a namespace and all its versions.

**Response:**
```json
{
  "message": "Namespace deleted successfully"
}
```

## Error Responses

All endpoints may return error responses in the following format:

```json
{
  "error": "Error description"
}
```

Common HTTP status codes:
- `200 OK`: Success
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Authentication required (TODO)
- `403 Forbidden`: Access denied (TODO)
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## Data Formats

### ACL Tuple Format
ACL tuples follow the format: `object#relation@user`

Examples:
- `doc:readme#viewer@user:alice`
- `folder:public#editor@user:bob`
- `file:config#owner@user:admin`

### Namespace Relations
Namespaces define relations between objects and users. Relations can be:

1. **Direct relations**: Users directly assigned to objects
2. **Computed usersets**: Relations computed from other relations using union operations

**TODO**: Full implementation of computed usersets and union operations is pending.

## Rate Limiting

**TODO**: Rate limiting is not yet implemented.

Future implementation will limit requests per IP/user:
- Default: 100 requests per minute
- Configurable via environment variables

## CORS

Cross-Origin Resource Sharing (CORS) is enabled by default for all origins.

**TODO**: Configure CORS for production use with specific allowed origins.