package routes

import (
	"DemoFormTutor/Form/database"
	"DemoFormTutor/Form/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	api := app.Group("/tutor")
	localDB := database.Connect()

	api.Get("/", func(c *fiber.Ctx) error { //Test
		return c.JSON(fiber.Map{})
	})

	handler := &handlers.TutorHandler{DB: localDB}
	api.Post("/save-tutor-data/:id", handler.SaveTutor) //Add Tutor's info
}
