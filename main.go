package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/mailwilliams/go-ambassador/src/database"
	"github.com/mailwilliams/go-ambassador/src/routes"
)

func main() {
	database.Connect()
	database.AutoMigrate()
	database.SetUpRedis()
	database.SetupCacheChannel()

	app := fiber.New()

	app.Use(cors.New(cors.Config{AllowCredentials: true}))

	routes.SetUp(app)
	app.Listen(":8000")
}
