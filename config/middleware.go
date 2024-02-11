package config

import (
	"fmt"
	"os"
	"p2h-api/app/models"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func Middleware(c *fiber.Ctx) error {
	jwtSecretKey := os.Getenv("JWT_SECRET")
	if jwtSecretKey == "" {
		fmt.Println("JWT Secret key not found in env variables")
	}

	authorizationHeader := c.Get("Authorization")
	tokenString := strings.Replace(authorizationHeader, "Bearer ", "", 1)

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"status":  false,
			"message": "Unauthorized",
		})
	}

	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"status":  false,
			"message": "Unauthorized",
		})
	}

	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		c.Locals("user", claims)
		return c.Next()
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"code":    fiber.StatusUnauthorized,
		"status":  false,
		"message": "Unauthorized",
	})
}

func AdminMiddleware(c *fiber.Ctx) error {
	user := c.Locals("user")
	userData, ok := user.(*models.CustomClaims)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":    fiber.StatusForbidden,
			"status":  false,
			"message": "Akses ditolak",
			"data":    fiber.Map{},
		})
	}

	if userData.Role != "adminsystem" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":    fiber.StatusForbidden,
			"status":  false,
			"message": "Akses ditolak",
			"data":    fiber.Map{},
		})
	}

	return c.Next()
}

func AdminCompanyMiddleware(c *fiber.Ctx) error {
	user := c.Locals("user")
	userData, ok := user.(*models.CustomClaims)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":    fiber.StatusForbidden,
			"status":  false,
			"message": "Akses ditolak",
			"data":    fiber.Map{},
		})
	}

	if userData.Role != "adminsystem" && userData.Role != "superadmin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":    fiber.StatusForbidden,
			"status":  false,
			"message": "Akses ditolak",
			"data":    fiber.Map{},
		})
	}

	return c.Next()
}
