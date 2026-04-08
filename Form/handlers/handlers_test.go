package handlers

import (
	"DemoFormTutor/Form/structs"
	"DemoFormTutor/Form/usecase"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type FakeUseCase struct {
	err     error
	countID int
}

func (f FakeUseCase) CreateTutor(ctx context.Context, tutor *structs.Tutor) error {
	if f.err != nil {
		return f.err
	}
	f.countID++
	tutor.ID = f.countID
	return nil
}

func NewFakeUseCase(err error) FakeUseCase {
	return FakeUseCase{err: err}
}

func TestSaveTutor_Success(t *testing.T) {

	handler := TutorHandler{NewFakeUseCase(nil)}

	tutor := structs.Tutor{
		ID:          0,
		Name:        "Amogus",
		Email:       "blabla@bla.bla",
		ExpWorkTime: "Не хватает!!!",
		Expectation: "Деньга",
		NeedCourses: false,
		TutorBefore: true,
	}
	buf := bytes.NewBuffer(nil)

	err := json.NewEncoder(buf).Encode(tutor)
	require.NoError(t, err)

	app := fiber.New()
	app.Post("/tutor/save-tutor-data/:id", handler.CreateTutor)

	req := httptest.NewRequest("POST", "/tutor/save-tutor-data/:id", buf)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	AssertTutorResponse(t, tutor, UnmarshalTutor(t, resp.Body))
}

func UnmarshalTutor(t *testing.T, rc io.ReadCloser) structs.Tutor {
	defer rc.Close()
	response := struct {
		Result structs.Tutor `json:"Result"`
	}{}
	err := json.NewDecoder(rc).Decode(&response)
	require.NoError(t, err)

	return response.Result
}

func AssertTutorResponse(t *testing.T, expected structs.Tutor, got structs.Tutor) {
	assert.NotEmpty(t, got.ID)

	expected.ID = got.ID
	assert.Equal(t, expected, got)
}

func TestSaveTutor_Fail(t *testing.T) {
	t.Run("Name empty", func(t *testing.T) {
		handler := TutorHandler{usecase.NewTutorUseCase(nil)}
		tutor := structs.Tutor{}

		buf := bytes.NewBuffer(nil)

		err := json.NewEncoder(buf).Encode(tutor)
		require.NoError(t, err)

		app := fiber.New()
		app.Post("/tutor/save-tutor-data/:id", handler.CreateTutor)

		req := httptest.NewRequest("POST", "/tutor/save-tutor-data/:id", buf)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
	t.Run("Email empty", func(t *testing.T) {
		handler := TutorHandler{usecase.NewTutorUseCase(nil)}
		tutor := structs.Tutor{
			Name: "Amogus",
		}

		buf := bytes.NewBuffer(nil)

		err := json.NewEncoder(buf).Encode(tutor)
		require.NoError(t, err)

		app := fiber.New()
		app.Post("/tutor/save-tutor-data/:id", handler.CreateTutor)

		req := httptest.NewRequest("POST", "/tutor/save-tutor-data/:id", buf)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
	t.Run("Wrong struct empty", func(t *testing.T) {
		type fakeStruct struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Email       string `json:"email"`
			ExpWorkTime string `json:"expworktime"`
			Expectation string `json:"expectation"`
			NeedCourses bool   `json:"needcourses"`
			TutorBefore string `json:"tutorbefore"`
		}

		handler := TutorHandler{usecase.NewTutorUseCase(nil)}

		tutor := fakeStruct{
			Name:        "Amogus",
			Email:       "blabla@bla.bla",
			TutorBefore: "false",
		}

		buf := bytes.NewBuffer(nil)

		err := json.NewEncoder(buf).Encode(tutor)
		require.NoError(t, err)

		app := fiber.New()
		app.Post("/tutor/save-tutor-data/:id", handler.CreateTutor)

		req := httptest.NewRequest("POST", "/tutor/save-tutor-data/:id", buf)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}
