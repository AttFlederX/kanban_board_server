package handlers

import "github.com/AttFlederX/kanban_board_server/services"

var (
	TaskService = services.NewMongoService("tasks")
	UserService = services.NewMongoService("users")
)

const (
	// Context keys
	contextKeyUserID = "userID"

	// BSON field names
	fieldName        = "name"
	fieldDescription = "description"
	fieldStatus      = "status"
	fieldUserID      = "userId"
	fieldPhotoURL    = "photourl"
	fieldGoogleID    = "google_id"
	fieldEmail       = "email"

	// JSON field names
	jsonFieldID      = "id"
	jsonFieldError   = "error"

	// Token payload claim keys
	claimEmail   = "email"
	claimName    = "name"
	claimPicture = "picture"

	// Error messages
	errInvalidUserID       = "Invalid user ID"
	errInvalidID           = "Invalid ID"
	errTaskNotFound        = "Task not found"
	errUserNotFound        = "User not found"
	errAccessDenied        = "Access denied"
	errInvalidRequestBody  = "Invalid request body"
	errIDTokenRequired     = "ID token is required"
	errInvalidGoogleToken  = "Invalid Google ID token"
	errFailedCreateUser    = "Failed to create user"
	errFailedUpdateUser    = "Failed to update user"
	errFailedGenerateToken = "Failed to generate token"
)
