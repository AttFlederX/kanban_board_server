# Flutter Integration Guide

## Overview

This Kanban Board API uses Google OAuth authentication and provides RESTful endpoints for task management. The server runs on **port 3000** by default.

**Base URL:** `http://localhost:3000`

---

## API Endpoints

### Authentication

**`POST /auth/google`** - Authenticate with Google

- **Public endpoint** (no auth required)
- **Request:** JSON body with `id_token` (from Google Sign-In)
- **Response:** JWT token and user object
- **Purpose:** Exchange Google ID token for application JWT

### Tasks

All task endpoints require **JWT authentication** via Bearer token in the `Authorization` header.

- **`GET /tasks`** - Retrieve all tasks
- **`GET /tasks/:id`** - Get a specific task by ID
- **`POST /tasks`** - Create a new task
- **`PUT /tasks/:id`** - Update an existing task
- **`DELETE /tasks/:id`** - Delete a task

---

## Authentication Flow

### 1. Client-Side Google Sign-In

- User initiates Google Sign-In in the Flutter app using `google_sign_in` package
- Google authentication completes and returns a Google ID token
- This token contains verified user information (email, name, photo, Google ID)

### 2. Exchange Google Token for JWT

- Flutter app sends the Google ID token to `POST /auth/google`
- Server verifies the Google token with Google's servers
- Server extracts user info (Google ID, email, name, photo URL)
- Server checks if user exists in MongoDB:
  - **New user:** Creates user record in database
  - **Existing user:** Updates user information (name, email, photo)
- Server generates a JWT token (expires in 24 hours)
- Server responds with JWT token and user object

### 3. Authenticated Requests

- Flutter app stores JWT token securely (using `flutter_secure_storage`)
- All subsequent API requests include JWT in the `Authorization` header as `Bearer <token>`
- Server validates JWT on each request to protected endpoints
- JWT contains user ID, email, and Google ID as claims

### 4. Session Management

- JWT tokens expire after 24 hours
- On app launch, restore saved JWT and validate it's still valid
- On token expiration, redirect user to sign in again
- Sign out: Delete stored JWT and clear Google Sign-In session

---

## Data Models

### User

- **`id`** - MongoDB ObjectID (string)
- **`google_id`** - Google user identifier
- **`email`** - User's email address
- **`name`** - User's display name
- **`photourl`** - Profile picture URL

### Task

- **`id`** - MongoDB ObjectID (string, optional on create)
- **`name`** - Task title
- **`description`** - Task details
- **`status`** - Task state: `'todo'`, `'in_progress'`, or `'done'`
- **`userId`** - ID of user who owns the task

---

## Request/Response Examples

### Authentication

**Request:**

```
POST /auth/google
Content-Type: application/json

{
  "id_token": "eyJhbGciOiJSUzI1NiIsImtpZCI6..."
}
```

**Response:**

```
200 OK

{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6...",
  "user": {
    "id": "674f4c8e9b8c123456789abc",
    "google_id": "1234567890",
    "email": "user@example.com",
    "name": "John Doe",
    "photourl": "https://lh3.googleusercontent.com/..."
  }
}
```

### Create Task

**Request:**

```
POST /tasks
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6...
Content-Type: application/json

{
  "name": "Implement login screen",
  "description": "Create UI for user authentication",
  "status": "todo",
  "userId": "674f4c8e9b8c123456789abc"
}
```

**Response:**

```
201 Created

{
  "id": "674f5d1a2c3d456789012def",
  "name": "Implement login screen",
  "description": "Create UI for user authentication",
  "status": "todo",
  "userId": "674f4c8e9b8c123456789abc"
}
```

### Update Task

**Request:**

```
PUT /tasks/674f5d1a2c3d456789012def
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6...
Content-Type: application/json

{
  "name": "Implement login screen",
  "description": "Create UI for user authentication",
  "status": "in_progress",
  "userId": "674f4c8e9b8c123456789abc"
}
```

**Response:**

```
200 OK

{
  "id": "674f5d1a2c3d456789012def",
  "name": "Implement login screen",
  "description": "Create UI for user authentication",
  "status": "in_progress",
  "userId": "674f4c8e9b8c123456789abc"
}
```

### Get All Tasks

**Request:**

```
GET /tasks
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6...
```

**Response:**

```
200 OK

[
  {
    "id": "674f5d1a2c3d456789012def",
    "name": "Implement login screen",
    "description": "Create UI for user authentication",
    "status": "in_progress",
    "userId": "674f4c8e9b8c123456789abc"
  },
  {
    "id": "674f5d1a2c3d456789012df0",
    "name": "Setup database",
    "description": "Configure MongoDB connection",
    "status": "done",
    "userId": "674f4c8e9b8c123456789abc"
  }
]
```

### Delete Task

**Request:**

```
DELETE /tasks/674f5d1a2c3d456789012def
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6...
```

**Response:**

```
204 No Content
```

---

## Error Responses

All endpoints may return these errors:

- **`400 Bad Request`** - Invalid request body or parameters
- **`401 Unauthorized`** - Missing or invalid JWT token
- **`404 Not Found`** - Resource doesn't exist
- **`500 Internal Server Error`** - Server-side error

Error response format:

```json
{
  "error": "Error message description"
}
```

---

## Implementation Notes

### Flutter Dependencies Required

- **`google_sign_in`** - Google OAuth flow
- **`http`** or **`dio`** - HTTP client for API calls
- **`flutter_secure_storage`** - Secure token storage
- **`provider`**, **`riverpod`**, or **`bloc`** - State management

### Connection URLs by Platform

- **Android Emulator:** `http://10.0.2.2:3000`
- **iOS Simulator:** `http://localhost:3000`
- **Physical Device:** `http://YOUR_LOCAL_IP:3000`
- **Production:** `https://your-domain.com`

### Security Best Practices

- Store JWT tokens in secure storage (never SharedPreferences)
- Use HTTPS in production
- Handle token expiration gracefully
- Clear tokens on sign out
- Validate token presence before making authenticated requests

### CORS Configuration

The server is configured to accept requests from any origin (`*`) for development. Update CORS settings for production to restrict to your Flutter app's domain.
