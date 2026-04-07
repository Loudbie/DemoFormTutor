package handlers

import (
	"DemoFormTutor/Form/repository"
	"DemoFormTutor/Form/structs"
	"DemoFormTutor/Form/usecase"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/stretchr/testify/require"
)

func TestSaveTutor_Success(t *testing.T) {
	tutorRepo, err := repository.NewPostgresTutorRepository("user=postgres dbname=Demo sslmode=disable password=1234955 port=9090 host=localhost")
	if err != nil {
		log.Fatalf("Cannot connect to db: %v", err)
	}
	tutorUC := usecase.NewTutorUseCase(tutorRepo)
	handler := NewTutorHandler(tutorUC)
	tutor := structs.Tutor{
		ID:          0,
		Name:        "Amogus",
		Email:       "blabla@bla.bla",
		ExpWorkTime: "Не хватает!!!",
		Expectation: "Деньга",
		NeedCourses: false,
		TutorBefore: true,
	}
	tutorString, err := json.Marshal(tutor)
	require.NoError(t, err)

	app := fiber.New()
	app.Post("/tutor/save-tutor-data/:id", handler.CreateTutor)

	req := httptest.NewRequest("POST", "/tutor/save-tutor-data/:id", bytes.NewReader(tutorString))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	require.Equal(t, fiber.StatusCreated, resp.StatusCode)
}
