package handlers

import (
	"github.com/AttFlederX/kanban_board_server/models"
	"github.com/AttFlederX/kanban_board_server/services"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTasks(c *fiber.Ctx) error {
	// Get authenticated user ID from context
	userID := c.Locals(contextKeyUserID).(string)

	// Convert string ID to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{jsonFieldError: errInvalidUserID})
	}

	// Find tasks belonging to the authenticated user
	tasks := []models.Task{}
	filter := bson.M{fieldUserID: userObjectID}
	if err := services.Find(collectionTasks, filter, &tasks); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{jsonFieldError: err.Error()})
	}
	return c.JSON(tasks)
}

func GetTask(c *fiber.Ctx) error {
	// Get authenticated user ID
	userID := c.Locals(contextKeyUserID).(string)
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{jsonFieldError: errInvalidUserID})
	}

	id, err := primitive.ObjectIDFromHex(c.Params(jsonFieldID))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{jsonFieldError: errInvalidID})
	}

	var task models.Task
	if err := services.FindByID(collectionTasks, id, &task); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{jsonFieldError: errTaskNotFound})
	}

	// Verify task belongs to authenticated user
	if task.UserID != userObjectID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{jsonFieldError: errAccessDenied})
	}

	return c.JSON(task)
}

func CreateTask(c *fiber.Ctx) error {
	// Get authenticated user ID
	userID := c.Locals(contextKeyUserID).(string)
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{jsonFieldError: errInvalidUserID})
	}

	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{jsonFieldError: err.Error()})
	}

	// Force task to belong to authenticated user
	task.UserID = userObjectID

	id, err := services.InsertOne(collectionTasks, task)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{jsonFieldError: err.Error()})
	}

	task.ID = id

	// Broadcast task creation to websocket clients
	BroadcastTaskChange("create", id, userObjectID, task)

	return c.Status(fiber.StatusCreated).JSON(task)
}

func UpdateTask(c *fiber.Ctx) error {
	// Get authenticated user ID
	userID := c.Locals(contextKeyUserID).(string)
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{jsonFieldError: errInvalidUserID})
	}

	id, err := primitive.ObjectIDFromHex(c.Params(jsonFieldID))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{jsonFieldError: errInvalidID})
	}

	// Check if task exists and belongs to user
	var existingTask models.Task
	if err := services.FindByID(collectionTasks, id, &existingTask); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{jsonFieldError: errTaskNotFound})
	}

	if existingTask.UserID != userObjectID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{jsonFieldError: errAccessDenied})
	}

	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{jsonFieldError: err.Error()})
	}

	update := bson.M{
		fieldName:        task.Name,
		fieldDescription: task.Description,
		fieldStatus:      task.Status,
		fieldUserID:      userObjectID,
	}
	if err := services.UpdateByID(collectionTasks, id, update); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{jsonFieldError: err.Error()})
	}

	task.ID = id
	task.UserID = userObjectID

	// Broadcast task update to websocket clients
	BroadcastTaskChange("update", id, userObjectID, task)

	return c.JSON(task)
}

func DeleteTask(c *fiber.Ctx) error {
	// Get authenticated user ID
	userID := c.Locals(contextKeyUserID).(string)
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{jsonFieldError: errInvalidUserID})
	}

	id, err := primitive.ObjectIDFromHex(c.Params(jsonFieldID))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{jsonFieldError: errInvalidID})
	}

	// Check if task exists and belongs to user
	var task models.Task
	if err := services.FindByID(collectionTasks, id, &task); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{jsonFieldError: errTaskNotFound})
	}

	if task.UserID != userObjectID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{jsonFieldError: errAccessDenied})
	}

	if err := services.DeleteByID(collectionTasks, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{jsonFieldError: err.Error()})
	}

	// Broadcast task deletion to websocket clients
	BroadcastTaskChange("delete", id, userObjectID, nil)

	return c.SendStatus(fiber.StatusNoContent)
}
