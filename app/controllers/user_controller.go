package controllers

import (
	"p2h-api/app/models"
	"p2h-api/app/requests/reqUser"
	"p2h-api/app/services"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{UserService: userService}
}

func (uc *UserController) Create(c *fiber.Ctx) error {
	user := c.Locals("user")
	userData, ok := user.(*models.CustomClaims)

	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"status":  false,
			"message": "Data user",
			"data":    fiber.Map{},
		})
	}

	var request reqUser.Save

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"status":  false,
			"message": "Payload kosong",
		})
	}

	file, err := c.FormFile("tanda_tangan")

	if err == nil {
		request.Tanda_Tangan = file
	}

	if errors, err := request.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":    fiber.StatusUnprocessableEntity,
			"status":  false,
			"message": "Inputan form belum sesuai",
			"data":    errors,
		})
	}

	result := uc.UserService.Create(userData.Role, request)
	return c.Status(result.Code).JSON(result)
}

func (uc *UserController) Update(c *fiber.Ctx) error {
	user := c.Locals("user")
	userData, ok := user.(*models.CustomClaims)

	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"status":  false,
			"message": "Data user",
			"data":    fiber.Map{},
		})
	}

	var request reqUser.Update

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"status":  false,
			"message": "Payload kosong",
		})
	}

	file, err := c.FormFile("tanda_tangan")

	if err == nil {
		request.Tanda_Tangan = file
	}

	if errors, err := request.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":    fiber.StatusUnprocessableEntity,
			"status":  false,
			"message": "Inputan form belum sesuai",
			"data":    errors,
		})
	}

	result := uc.UserService.Update(userData.Role, request)
	return c.Status(result.Code).JSON(result)
}

func (uc *UserController) GetAllUser(c *fiber.Ctx) error {
	var request reqUser.List
	user := c.Locals("user")
	userData, ok := user.(*models.CustomClaims)

	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"status":  false,
			"message": "Data user",
			"data":    fiber.Map{},
		})
	}

	if err := c.QueryParser(&request); err != nil {
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
			"total":   0,
		})
	}

	result := uc.UserService.GetAllUser(userData, request)
	return c.Status(result.Code).JSON(result)
}

func (uc *UserController) Delete(c *fiber.Ctx) error {
	param := c.Params("uuid")
	user := c.Locals("user")
	userData, ok := user.(*models.CustomClaims)

	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"status":  false,
			"message": "Data user",
			"data":    fiber.Map{},
		})
	}

	result := uc.UserService.Delete(userData, param)
	return c.Status(result.Code).JSON(result)
}

func (uc *UserController) UpdateProfile(c *fiber.Ctx) error {
	user := c.Locals("user")
	userData, ok := user.(*models.CustomClaims)

	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"status":  false,
			"message": "Data user",
			"data":    fiber.Map{},
		})
	}

	var request reqUser.UpdateProfile

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"status":  false,
			"message": "Payload kosong",
		})
	}

	file, err := c.FormFile("tanda_tangan")

	if err == nil {
		request.Tanda_Tangan = file
	}

	if errors, err := request.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":    fiber.StatusUnprocessableEntity,
			"status":  false,
			"message": "Inputan form belum sesuai",
			"data":    errors,
		})
	}

	result := uc.UserService.UpdateProfile(userData.ID, request)
	return c.Status(result.Code).JSON(result)
}

func (uc *UserController) Profile(c *fiber.Ctx) error {
	user := c.Locals("user")
	userData, ok := user.(*models.CustomClaims)

	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"status":  false,
			"message": "Data user",
			"data":    fiber.Map{},
		})
	}

	result := uc.UserService.GetUserByID(userData.ID)
	return c.Status(result.Code).JSON(result)
}
