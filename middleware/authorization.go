package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"rakamin.com/project/config"
	"rakamin.com/project/models"
)

func Authorization(c *fiber.Ctx) error {
	// Get gofiber context
	ctx := c.Locals("ctx").(*models.Context)

	// Get Authorization header , bearer ....
	auth := c.Get("authorization", "")

	if auth != "" {
		auth_splice := strings.Split(auth, " ")
		auth_type := strings.ToLower(auth_splice[0])
		token_string := auth_splice[1]

		if auth_type == "bearer" {
			if token_string == "" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
			}

			// Validating JWT
			token, err := jwt.Parse(token_string, func(t *jwt.Token) (interface{}, error) {
				if method, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("signing method invalid")
				} else if method != config.JWT_SIGNING_METHOD {
					return nil, fmt.Errorf("signing method invalid")
				}
				return config.JWT_SIGNATURE_KEY, nil
			})

			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": fmt.Sprintf("Invalid token: %s", err.Error())})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
			}

			var user models.Users

			// Get users from jwt sub
			if err := config.DB.First(&user, "id = ?", claims["sub"]).Error; err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
			}

			// Set gofiber context
			ctx.User = user
			c.Locals("ctx", ctx)

			// Next
			return c.Next()
		}
	}

	// Unauthorized
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
}
