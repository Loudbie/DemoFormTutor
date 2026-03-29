package routes

import (
	"DemoFormTutor/Form/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	api := app.Group("/tutor")

	api.Get("/", func(c *fiber.Ctx) error { //Test
		return c.JSON(fiber.Map{})
	})

	api.Post("/save-tutor-data/:id", handlers.SaveTutor) //Add Tutor's info
}
