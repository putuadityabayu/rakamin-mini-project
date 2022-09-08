package routes

import (
	"github.com/gofiber/fiber/v2"
	"rakamin.com/project/controllers"
)

func SetupConversation(r fiber.Router) {
	r.Post("/", controllers.NewMessage)
	r.Get("/", controllers.ListConversation)
	r.Delete("/:id", controllers.DeleteConversation)
}
