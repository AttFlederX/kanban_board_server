package main

import (
	"log"

	"github.com/AttFlederX/kanban_board_server/config"
	"github.com/AttFlederX/kanban_board_server/database"
	"github.com/AttFlederX/kanban_board_server/handlers"
	"github.com/AttFlederX/kanban_board_server/middleware"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	cfg := config.Load()

	log.Println("Starting server on port", cfg.Port)
	log.Println("Connecting to MongoDB at", cfg.MongoURI)

	if err := database.Connect(cfg.MongoURI, cfg.DBName); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Initialize websocket hub
	handlers.SetJWTSecret(cfg.JWTSecret)
	handlers.InitHub()

	app := fiber.New()

	// Enable CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Auth routes (public)
	app.Post("/auth/google", func(c *fiber.Ctx) error {
		return handlers.GoogleSignIn(c, cfg.JWTSecret)
	})

	// WebSocket route (handles auth via token query param)
	app.Get("/ws", websocket.New(handlers.HandleWebSocket))

	// Protected routes
	authApp := app.Group("", middleware.AuthRequired(cfg.JWTSecret))

	// User routes (protected)
	authApp.Get("/users/:id", handlers.GetUser)
	authApp.Post("/users", handlers.CreateUser)
	authApp.Put("/users/:id", handlers.UpdateUser)
	authApp.Delete("/users/:id", handlers.DeleteUser)

	// Task routes (protected)
	authApp.Get("/tasks", handlers.GetTasks)
	authApp.Get("/tasks/:id", handlers.GetTask)
	authApp.Post("/tasks", handlers.CreateTask)
	authApp.Put("/tasks/:id", handlers.UpdateTask)
	authApp.Delete("/tasks/:id", handlers.DeleteTask)

	log.Fatal(app.Listen(":" + cfg.Port))
}
