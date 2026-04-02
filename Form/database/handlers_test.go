package database

import (
	"DemoFormTutor/Form/handlers"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Фейк для DatabaseExecutor

// fakeResult фейковый sql.Result
type fakeResult struct{}

func (f fakeResult) LastInsertId() (int64, error) {
	return 1, nil
}
func (f fakeResult) RowsAffected() (int64, error) {
	return 1, nil
}

// fakeDB - фейковый DatabaseExecutor
type fakeDB struct {
	ShouldFail bool
	LastQuery  string
	LastArgs   []interface{}
}

func (db *fakeDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	db.LastQuery = query
	db.LastArgs = args

	if db.ShouldFail {
		return nil, errors.New("db error")
	}
	return fakeResult{}, nil
}

// Fiber-приложение с хендлером
func setupApp(db handlers.DatabaseExecutor) *fiber.App {
	app := fiber.New()
	handler := &handlers.TutorHandler{DB: db}
	app.Post("/tutor", handler.SaveTutor)
	return app
}

// Тест №1:Успешное сохранение репетитора (Код ошибки: 201)
func TestSaveTutor_Success(t *testing.T) {
	db := &fakeDB{}
	app := setupApp(db)

	body := `{
		"name":         "Иван Иванов",
		"email":        "ivan@example.com",
		"expworktime":  "5 лет",
		"expectation":  "Хорошая зарплата",
		"needcourses":  "Да",
		"tutorbefore":  "Нет"
	}`

	req := httptest.NewRequest("POST", "/tutor", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	require.NoError(t, err)

	assert.Equal(t, "Репетитор успешно добавлен:", result["message"])

	tutor, ok := result["tutor"].(map[string]interface{})
	require.True(t, ok)

	assert.Equal(t, "Иван Иванов", tutor["Имя"])
	assert.Equal(t, "ivan@example.com", tutor["Почта"])
	assert.Equal(t, "5 лет", tutor["Опыт работы"])
	assert.Equal(t, "Хорошая зарплата", tutor["Ожидания от работы"])
	assert.Equal(t, "Да", tutor["Нужны ли курсы"])
	assert.Equal(t, "Нет", tutor["Преподавал ли до этого"])

	assert.Contains(t, db.LastQuery, "INSERT INTO tutor")
	assert.Equal(t, "Иван Иванов", db.LastArgs[0])
	assert.Equal(t, "ivan@example.com", db.LastArgs[1])
}

// Тест №2 Ошибка базы данных (Код ошибки:500)
func TestSaveTutor_DatabaseError(t *testing.T) {
	db := &fakeDB{ShouldFail: true}
	app := setupApp(db)

	body := `{
		"name": "Test",
		"email": "test@test.com"
	}`
	req := httptest.NewRequest("POST", "/tutor", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	require.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "Ошибка вставки данных в базу", string(respBody))
}

// Тест №3 Пустое тело (Код ошибки:400)
func TestSaveTutor_EmptyBody(t *testing.T) {
	db := &fakeDB{}
	app := setupApp(db)

	req := httptest.NewRequest("POST", "/tutor", nil)
	resp, err := app.Test(req, -1)

	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}
