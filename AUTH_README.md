# Google Authentication Implementation

## Overview

This server now supports Google Sign-In authentication. Users authenticate via Google OAuth on the client side, then send the Google ID token to the server for verification and JWT token generation.

## Authentication Flow

1. **Client Side**: User signs in with Google and receives an ID token
2. **Server Side**:
   - Client sends ID token to `/auth/google` endpoint
   - Server verifies the Google ID token
   - Server creates/updates user in database
   - Server generates JWT token
   - Server returns JWT token and user info to client
3. **Protected Routes**: Client includes JWT token in `Authorization` header as `Bearer <token>`

## API Endpoints

### Public Endpoints

#### POST `/auth/google`

Authenticate with Google Sign-In ID token.

**Request Body:**

```json
{
  "id_token": "google_id_token_from_client"
}
```

**Response (200 OK):**

```json
{
  "token": "jwt_token",
  "user": {
    "id": "user_object_id",
    "google_id": "google_user_id",
    "email": "user@example.com",
    "name": "User Name",
    "photourl": "https://profile.photo.url"
  }
}
```

**Error Responses:**

- `400 Bad Request`: Invalid request body or missing ID token
- `401 Unauthorized`: Invalid Google ID token
- `500 Internal Server Error`: Database or token generation error

### Protected Endpoints

All user and task endpoints now require authentication. Include the JWT token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

**User Routes:**

- `GET /users` - Get all users
- `GET /users/:id` - Get user by ID
- `POST /users` - Create new user
- `PUT /users/:id` - Update user
- `DELETE /users/:id` - Delete user

**Task Routes:**

- `GET /tasks` - Get all tasks
- `GET /tasks/:id` - Get task by ID
- `POST /tasks` - Create new task
- `PUT /tasks/:id` - Update task
- `DELETE /tasks/:id` - Delete task

## Configuration

The server requires a JWT secret for signing tokens. Set the `JWT_SECRET` environment variable:

```bash
export JWT_SECRET="your-secure-secret-key"
```

Default value (for development only): `"your-secret-key-change-this-in-production"`

## User Model

The User model has been updated to include Google authentication fields:

```go
type User struct {
    ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    GoogleID string             `json:"google_id" bson:"google_id"`
    Email    string             `json:"email" bson:"email"`
    Name     string             `json:"name" bson:"name"`
    PhotoURL string             `json:"photourl" bson:"photourl"`
}
```

## Security Features

- **Google ID Token Verification**: Server validates tokens directly with Google
- **JWT Token Authentication**: Stateless authentication using JWT tokens
- **Token Expiration**: JWT tokens expire after 24 hours
- **Protected Routes**: All user and task endpoints require valid JWT token
- **CORS Enabled**: Configured for cross-origin requests

## Client Implementation Example

```javascript
// 1. Get Google ID token from Google Sign-In
const googleUser = await googleSignIn();
const idToken = googleUser.credential;

// 2. Send to server for verification
const response = await fetch("http://localhost:3000/auth/google", {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
  },
  body: JSON.stringify({ id_token: idToken }),
});

const { token, user } = await response.json();

// 3. Store JWT token for subsequent requests
localStorage.setItem("token", token);

// 4. Use token for authenticated requests
const tasksResponse = await fetch("http://localhost:3000/tasks", {
  headers: {
    Authorization: `Bearer ${token}`,
  },
});
```

## Environment Variables

- `MONGO_URI` - MongoDB connection string (default: `mongodb://admin:admin@localhost:27017/?directConnection=true&serverSelectionTimeoutMS=2000`)
- `DB_NAME` - Database name (default: `kanban_board`)
- `PORT` - Server port (default: `3000`)
- `JWT_SECRET` - Secret key for JWT signing (default: `your-secret-key-change-this-in-production`)

## Dependencies

- `github.com/golang-jwt/jwt/v5` - JWT token generation and validation
- `google.golang.org/api/idtoken` - Google ID token verification
- `github.com/gofiber/fiber/v2` - Web framework
- `go.mongodb.org/mongo-driver` - MongoDB driver
