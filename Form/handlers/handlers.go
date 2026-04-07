package handlers

import (
	"DemoFormTutor/Form/structs"
	"DemoFormTutor/Form/usecase"

	"github.com/gofiber/fiber/v2"
)

type TutorHandler struct {
	tutorUC usecase.TutorUseCase
}

func NewTutorHandler(us usecase.TutorUseCase) *TutorHandler {
	return &TutorHandler{tutorUC: us}
}

func (h *TutorHandler) CreateTutor(ctx *fiber.Ctx) error {
	var tutor structs.Tutor
	if err := ctx.BodyParser(&tutor); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.tutorUC.CreateTutor(ctx, &tutor); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"Tutor added": tutor})
}
