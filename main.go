package main

import (
	"fmt"
	"os"
	"p2h-api/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	app := fiber.New()
	app.Use(cors.New())

	config.SetupRouter(app)
	err := app.Listen(":" + os.Getenv("APP_PORT"))

	if err != nil {
		panic(err)
	}
}
