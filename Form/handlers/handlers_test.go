package handlers

import (
	"DemoFormTutor/Form/structs"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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

func (f FakeUseCase) GetAllTutors(ctx context.Context, sort string) ([]structs.Tutor, error) {
	if f.err != nil {
		return nil, f.err
	}
	return []structs.Tutor{}, nil
}

func (f FakeUseCase) FindTutorById(ctx context.Context, id string) (*structs.Tutor, error) {
	switch {
	case errors.Is(f.err, sql.ErrNoRows):
		return nil, f.err
	case f.err != nil:
		return nil, f.err
	}
	return &structs.Tutor{}, nil
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

func TestTutorHandler(t *testing.T) {
	t.Run("Success_Save", func(t *testing.T) {
		handler := TutorHandler{tutorUC: NewFakeUseCase(nil)}

		testApp := fiber.New()
		testApp.Post("/tutor/save-tutor-data/:id", handler.CreateTutor)

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

		req := httptest.NewRequest("POST", "/tutor/save-tutor-data/:id", buf)
		req.Header.Set("Content-Type", "application/json")

		resp, err := testApp.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		AssertTutorResponse(t, tutor, UnmarshalTutor(t, resp.Body))
	})
	t.Run("Wrong_struct_empty_Save", func(t *testing.T) {
		handler := TutorHandler{tutorUC: NewFakeUseCase(nil)}

		testApp := fiber.New()
		testApp.Post("/tutor/save-tutor-data/:id", handler.CreateTutor)

		type fakeStruct struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Email       string `json:"email"`
			ExpWorkTime string `json:"expworktime"`
			Expectation string `json:"expectation"`
			NeedCourses bool   `json:"needcourses"`
			TutorBefore string `json:"tutorbefore"`
		}

		tutor := fakeStruct{
			Name:        "Amogus",
			Email:       "blabla@bla.bla",
			TutorBefore: "false",
		}

		buf := bytes.NewBuffer(nil)
		err := json.NewEncoder(buf).Encode(tutor)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/tutor/save-tutor-data/:id", buf)
		req.Header.Set("Content-Type", "application/json")

		resp, err := testApp.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
	t.Run("Success_FindById", func(t *testing.T) {
		handler := TutorHandler{tutorUC: NewFakeUseCase(nil)}

		testApp := fiber.New()
		testApp.Get("/tutor/find-tutor-data/:id", handler.GetTutor)

		req := httptest.NewRequest("GET", "/tutor/find-tutor-data/some_id", nil)

		resp, err := testApp.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
	t.Run("NoRows_Error_FindById", func(t *testing.T) {
		handler := TutorHandler{tutorUC: NewFakeUseCase(sql.ErrNoRows)}

		testApp := fiber.New()
		testApp.Get("/tutor/find-tutor-data/:id", handler.GetTutor)

		req := httptest.NewRequest("GET", "/tutor/find-tutor-data/give_me_error", nil)

		resp, err := testApp.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
	t.Run("Random_Error_FindById", func(t *testing.T) {
		handler := TutorHandler{tutorUC: NewFakeUseCase(errors.New("Error Here!"))}

		testApp := fiber.New()
		testApp.Get("/tutor/find-tutor-data/:id", handler.GetTutor)

		req := httptest.NewRequest("GET", "/tutor/find-tutor-data/give_me_error", nil)

		resp, err := testApp.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
	t.Run("Success_GetAllTutors", func(t *testing.T) {
		handler := TutorHandler{tutorUC: NewFakeUseCase(nil)}

		testApp := fiber.New()
		testApp.Get("/tutor/find-tutors/:sort", handler.GetAllTutors)

		req := httptest.NewRequest("GET", "/tutor/find-tutors/some_sort_args", nil)

		resp, err := testApp.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
	t.Run("Error_GetAllTutors", func(t *testing.T) {
		handler := TutorHandler{tutorUC: NewFakeUseCase(errors.New("Error Here!"))}

		testApp := fiber.New()
		testApp.Get("/tutor/find-tutors/:sort", handler.GetAllTutors)

		req := httptest.NewRequest("GET", "/tutor/find-tutors/give_me_error", nil)

		resp, err := testApp.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
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
