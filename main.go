package main

import (
	"log"

	"github.com/AttFlederX/kanban_board_server/config"
	"github.com/AttFlederX/kanban_board_server/database"
	"github.com/AttFlederX/kanban_board_server/handlers"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.Load()

	log.Println("Starting server on port", cfg.Port)
	log.Println("Connecting to MongoDB at", cfg.MongoURI)

	if err := database.Connect(cfg.MongoURI, cfg.DBName); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	app := fiber.New()

	// User routes
	app.Get("/users", handlers.GetUsers)
	app.Get("/users/:id", handlers.GetUser)
	app.Post("/users", handlers.CreateUser)
	app.Put("/users/:id", handlers.UpdateUser)
	app.Delete("/users/:id", handlers.DeleteUser)

	// Task routes
	app.Get("/tasks", handlers.GetTasks)
	app.Get("/tasks/:id", handlers.GetTask)
	app.Post("/tasks", handlers.CreateTask)
	app.Put("/tasks/:id", handlers.UpdateTask)
	app.Delete("/tasks/:id", handlers.DeleteTask)

	log.Fatal(app.Listen(":" + cfg.Port))
}
