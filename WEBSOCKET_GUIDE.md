# WebSocket Implementation Guide

## Overview

The Kanban Board Server now supports real-time task updates via WebSockets. When a user creates, updates, or deletes a task, all connected clients for that user will receive instant notifications.

## Connection

### Endpoint

```
ws://localhost:<PORT>/ws?userId=<USER_ID>
```

### Authentication

The WebSocket endpoint is protected by the same JWT authentication middleware as other routes. Include the JWT token in the connection headers:

```
Authorization: Bearer <JWT_TOKEN>
```

### Parameters

- `userId`: The MongoDB ObjectID of the authenticated user (required in query string)

## Message Format

All WebSocket messages follow this JSON structure:

```json
{
  "type": "create|update|delete",
  "taskId": "507f1f77bcf86cd799439011",
  "userId": "507f1f77bcf86cd799439012",
  "data": {
    // Task object (for create/update only)
  }
}
```

### Message Types

#### 1. Task Created

```json
{
  "type": "create",
  "taskId": "507f1f77bcf86cd799439011",
  "userId": "507f1f77bcf86cd799439012",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "New Task",
    "description": "Task description",
    "status": "todo",
    "userId": "507f1f77bcf86cd799439012"
  }
}
```

#### 2. Task Updated

```json
{
  "type": "update",
  "taskId": "507f1f77bcf86cd799439011",
  "userId": "507f1f77bcf86cd799439012",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "Updated Task",
    "description": "Updated description",
    "status": "in-progress",
    "userId": "507f1f77bcf86cd799439012"
  }
}
```

#### 3. Task Deleted

```json
{
  "type": "delete",
  "taskId": "507f1f77bcf86cd799439011",
  "userId": "507f1f77bcf86cd799439012",
  "data": null
}
```

## Client Implementation Examples

### JavaScript (Browser)

```javascript
const userId = "507f1f77bcf86cd799439012";
const token = "your_jwt_token";

const ws = new WebSocket(`ws://localhost:8080/ws?userId=${userId}`);

// Note: Most WebSocket implementations don't support custom headers
// You may need to pass the token via query parameter or use a library
// that supports headers like Socket.IO

ws.onopen = () => {
  console.log("WebSocket connected");
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);

  switch (message.type) {
    case "create":
      console.log("Task created:", message.data);
      // Update UI to show new task
      break;
    case "update":
      console.log("Task updated:", message.data);
      // Update UI to reflect task changes
      break;
    case "delete":
      console.log("Task deleted:", message.taskId);
      // Remove task from UI
      break;
  }
};

ws.onerror = (error) => {
  console.error("WebSocket error:", error);
};

ws.onclose = () => {
  console.log("WebSocket disconnected");
  // Implement reconnection logic here
};
```

### Flutter (Dart)

```dart
import 'package:web_socket_channel/web_socket_channel.dart';
import 'dart:convert';

class WebSocketService {
  WebSocketChannel? _channel;
  final String userId;
  final String token;

  WebSocketService(this.userId, this.token);

  void connect() {
    _channel = WebSocketChannel.connect(
      Uri.parse('ws://localhost:8080/ws?userId=$userId'),
    );

    _channel!.stream.listen(
      (message) {
        final data = json.decode(message);
        handleMessage(data);
      },
      onError: (error) {
        print('WebSocket error: $error');
      },
      onDone: () {
        print('WebSocket closed');
        // Implement reconnection logic
      },
    );
  }

  void handleMessage(Map<String, dynamic> message) {
    switch (message['type']) {
      case 'create':
        print('Task created: ${message['data']}');
        // Update state/UI
        break;
      case 'update':
        print('Task updated: ${message['data']}');
        // Update state/UI
        break;
      case 'delete':
        print('Task deleted: ${message['taskId']}');
        // Update state/UI
        break;
    }
  }

  void disconnect() {
    _channel?.sink.close();
  }
}
```

## Features

- **User Isolation**: Each user only receives updates for their own tasks
- **Multiple Connections**: A user can have multiple WebSocket connections (e.g., from different devices)
- **Automatic Cleanup**: Connections are automatically cleaned up when clients disconnect
- **Thread-Safe**: The hub uses mutex locks to ensure thread-safe operations

## Testing

You can test the WebSocket connection using tools like:

- **wscat**: `wscat -c "ws://localhost:8080/ws?userId=YOUR_USER_ID"`
- **Postman**: Has built-in WebSocket support
- Browser Developer Tools: Use the browser console with the JavaScript example above

## Notes

- The WebSocket endpoint requires authentication like all other protected routes
- Clients should implement reconnection logic for handling disconnections
- The `userId` parameter must match an authenticated user's ID
- Keep-alive pings are handled automatically by reading client messages
