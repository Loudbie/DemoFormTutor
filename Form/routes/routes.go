package routes

import (
	"layers/Form/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error { //Test
		return c.JSON(fiber.Map{})
	})

	api.Post("/Tutor/:id", handlers.SaveTutor) //Add Tutor's info
}
