package handlers

import (
	"github.com/AttFlederX/kanban_board_server/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUser(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params(jsonFieldID))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{jsonFieldError: errInvalidID})
	}

	var user models.User
	if err := UserService.FindByID(id, &user); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{jsonFieldError: errUserNotFound})
	}

	return c.JSON(user)
}

func CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{jsonFieldError: err.Error()})
	}

	id, err := UserService.InsertOne(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{jsonFieldError: err.Error()})
	}

	user.ID = id
	return c.Status(fiber.StatusCreated).JSON(user)
}

func UpdateUser(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params(jsonFieldID))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{jsonFieldError: errInvalidID})
	}

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{jsonFieldError: err.Error()})
	}

	update := bson.M{fieldName: user.Name, fieldPhotoURL: user.PhotoURL}
	if err := UserService.UpdateByID(id, update); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{jsonFieldError: err.Error()})
	}

	user.ID = id
	return c.JSON(user)
}

func DeleteUser(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params(jsonFieldID))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{jsonFieldError: errInvalidID})
	}

	if err := UserService.DeleteByID(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{jsonFieldError: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
