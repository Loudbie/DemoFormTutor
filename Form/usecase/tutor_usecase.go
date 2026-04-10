package usecase

import (
	"DemoFormTutor/Form/repository"
	"DemoFormTutor/Form/structs"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type TutorUseCase interface {
	CreateTutor(ctx context.Context, tutor *structs.Tutor) error
	FindTutorById(ctx context.Context, id string) (*structs.Tutor, error)
	GetAllTutors(ctx context.Context, sort string) ([]structs.Tutor, error)
}
type tutorUseCase struct {
	tutorRepo repository.TutorRepository
}

func NewTutorUseCase(repo repository.TutorRepository) TutorUseCase {
	return &tutorUseCase{tutorRepo: repo}
}

func (u *tutorUseCase) CreateTutor(ctx context.Context, tutor *structs.Tutor) error {
	switch {
	//Проверка на отсутствие ввода имени
	case tutor.Name == "":
		return fmt.Errorf("Имя нужно обязательно.")
	//Проверка на отстутствие ввода почты
	case tutor.Email == "":
		return fmt.Errorf("Email нужен обязательно")
	}

	_, err := u.tutorRepo.Save(ctx, tutor)
	if err != nil {
		return fmt.Errorf("failed save tutor: %w", err)
	}

	return nil
}

func (u *tutorUseCase) FindTutorById(ctx context.Context, id string) (*structs.Tutor, error) {
	t, err := u.tutorRepo.FindById(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, err
		default:
			return nil, fmt.Errorf("failed find tutor by id: %w", err)
		}
	}
	return t, nil
}

func (u *tutorUseCase) GetAllTutors(ctx context.Context, sort string) ([]structs.Tutor, error) {
	ts, err := u.tutorRepo.GetAll(ctx, sort)
	if err != nil {
		return nil, fmt.Errorf("failed get all tutors: %w", err)
	}
	return ts, nil
}
