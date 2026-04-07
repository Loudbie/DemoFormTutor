package usecase

import (
	"DemoFormTutor/Form/repository"
	"DemoFormTutor/Form/structs"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type TutorUseCase interface {
	CreateTutor(ctx *fiber.Ctx, tutor *structs.Tutor) error
}
type tutorUseCase struct {
	tutorRepo repository.TutorRepository
}

func NewTutorUseCase(repo repository.TutorRepository) TutorUseCase {
	return &tutorUseCase{tutorRepo: repo}
}

func (u *tutorUseCase) CreateTutor(ctx *fiber.Ctx, tutor *structs.Tutor) error {

	switch {
	//Проверка на отсутствие ввода имени
	case tutor.Name == "":
		return fmt.Errorf("Имя нужно обязательно.")
	//Проверка на отстутствие ввода почты
	case tutor.Email == "":
		return fmt.Errorf("Email нужен обязательно")
	}
	return u.tutorRepo.Save(ctx, tutor)
}
