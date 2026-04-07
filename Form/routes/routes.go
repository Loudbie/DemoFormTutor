package routes

import (
	"DemoFormTutor/Form/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, handler *handlers.TutorHandler) {
	api := app.Group("/tutor")

	api.Get("/", func(c *fiber.Ctx) error { //Test
		return c.JSON(fiber.Map{})
	})

	api.Post("/save-tutor-data/:id", handler.CreateTutor) //Add Tutor's info
}
