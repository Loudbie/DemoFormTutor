package database

import (
	"DemoFormTutor/Form/handlers"
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Fiber-приложение с хендлером

func setupApp() *fiber.App {
	app := fiber.New()
	app.Post("/tutor/save-tutor-data", handler.SaveTutor)
	return app
}

var handler = &handlers.TutorHandler{DB: localDB}

type TutorHandler struct {
	DB handlers.DatabaseExecutor
}

func TestSaveTutor_Success(t *testing.T) {
	// Используем правильную структуру хендлера

	h := &handlers.TutorHandler{DB: localDB}
	app := setupApp()
	// Исправленная структура JSON с английскими полями
	tutorData := map[string]interface{}{
		"name":        "Иван Иванов",
		"email":       "ivan@example.com",
		"expworktime": "5 лет",
		"expectation": "Хорошая зарплата",
		"needcourses": true,
		"tutorbefore": false,
	}

	body, err := json.Marshal(tutorData)
	require.NoError(t, err)
	t.Logf("Request body: %s", body)
	req := httptest.NewRequest(fiber.MethodPost, "/tutor/save-tutor-data", bytes.NewReader(body))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Проверяем статус код
	require.Equal(t, 201, resp.StatusCode, "Expected status code 201, got %d", resp.StatusCode)

	// Читаем и парсим ответ
	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	require.NoError(t, err)

	// Проверяем структуру ответа
	assert.Equal(t, "Репетитор успешно добавлен", result["message"])

	tutor, ok := result["tutor"].(map[string]interface{})
	require.True(t, ok)

	// Проверяем поля репетитора
	assert.Equal(t, "Иван Иванов", tutor["name"])
	assert.Equal(t, "ivan@example.com", tutor["email"])
	assert.Equal(t, "5 лет", tutor["expworktime"])
	assert.Equal(t, "Хорошая зарплата", tutor["expectation"])
	assert.Equal(t, 1, tutor["needcourses"])
	assert.Equal(t, 0, tutor["tutorbefore"])

	// Проверяем, что запись действительно добавилась в БД
	var count int
	err = h.DB.QueryRow("SELECT COUNT(*) FROM tutor WHERE name = $1 AND email = $2",
		"Иван Иванов", "ivan@example.com").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count, "Record should be inserted into database")
}
