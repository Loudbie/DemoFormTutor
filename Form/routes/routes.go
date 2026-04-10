package routes

import (
	"DemoFormTutor/Form/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, handler *handlers.TutorHandler) {
	api := app.Group("/tutor")

	api.Get("/find-tutors/:sort", handler.GetAllTutors) //Get all Tutors and their info
	api.Get("/find-tutor-data/:id", handler.GetTutor)   //Find Tutor's info for Tutor's email

	api.Post("/save-tutor-data/:id", handler.CreateTutor) //Add Tutor's info
}
