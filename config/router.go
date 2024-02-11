package config

import (
	"p2h-api/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App) {
	db := ConnectToDB()
	repository := NewRepositoryRegistry(db)
	service := NewServiceRegistry(repository)
	publicRoutes(app, service)
	privateRoutes(app, service)
}

func publicRoutes(app *fiber.App, service *ServiceRegistry) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"code":    fiber.StatusOK,
			"message": "Testing API v1 Running",
		})
	})

	authController := controllers.NewAuthController(service.UserService)
	app.Post("/login", authController.Login)
}

func privateRoutes(app *fiber.App, service *ServiceRegistry) {
	userController := controllers.NewUserController(service.UserService)

	app.Use(Middleware)
	app.Get("/profile", userController.Profile)
	app.Put("/profile", userController.UpdateProfile)

	app.Get("/user", AdminCompanyMiddleware, userController.GetAllUser)
	app.Post("/user", AdminCompanyMiddleware, userController.Create)
	app.Put("/user", AdminCompanyMiddleware, userController.Update)
	app.Delete("/user/:uuid", AdminCompanyMiddleware, userController.Delete)
}
