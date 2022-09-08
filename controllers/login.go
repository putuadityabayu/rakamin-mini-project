package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"rakamin.com/project/config"
	"rakamin.com/project/models"
)

func Login(c *fiber.Ctx) error {
	var post models.Login

	// Parse POST body
	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing username/password"})
	}

	var user models.Users
	db := config.DB

	// Get users from username
	if err := db.First(&user, "username = ?", post.UserName).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid username"})
	}

	// Compare md5 password
	hash := md5.Sum([]byte(post.Password))
	pass_str := hex.EncodeToString(hash[:])
	if pass_str != user.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid password"})
	}

	// Create JWT token
	claims := jwt.RegisteredClaims{
		Issuer:    config.APPLICATION_NAME,
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(config.LOGIN_EXPIRATION_DURATION)},
		Audience:  jwt.ClaimStrings{config.JWT_ISSUER},
		Subject:   fmt.Sprintf("%d", user.ID), // User ID
	}

	token := jwt.NewWithClaims(
		config.JWT_SIGNING_METHOD,
		claims,
	)

	signedToken, err := token.SignedString(config.JWT_SIGNATURE_KEY)
	if err != nil {
		panic(err)
	}

	// Sent JWT to client
	return c.Status(200).JSON(fiber.Map{"token": signedToken, "user": user})
}
