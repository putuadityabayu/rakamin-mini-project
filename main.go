package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"rakamin.com/project/config"
	"rakamin.com/project/models"
	"rakamin.com/project/routes"
)

func newApp() *fiber.App {
	// Setup gofiber
	r := routes.SetupRouters()
	return r
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err.Error())
	}
	// Setup Database
	config.Initialization()
	models.SetupModels()
	err = newApp().Listen(fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		panic(err)
	}
}
