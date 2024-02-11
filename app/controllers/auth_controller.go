package controllers

import (
	"p2h-api/app/requests"
	"p2h-api/app/services"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	UserService *services.UserService
}

func NewAuthController(userService *services.UserService) *AuthController {
	return &AuthController{UserService: userService}
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
	var request requests.LoginRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"status":  false,
			"message": "Payload kosong",
		})
	}

	if errors, err := request.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":    fiber.StatusUnprocessableEntity,
			"status":  false,
			"message": "Inputan form belum sesuai",
			"data":    errors,
		})
	}

	result, _ := ac.UserService.UserAuthentication(request.Email, request.Password)
	return c.Status(result.Code).JSON(result)
}
