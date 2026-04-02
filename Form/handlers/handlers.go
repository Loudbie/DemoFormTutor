package handlers

import (
	"DemoFormTutor/Form/structs"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

/*
TODO Переиспользование кода (разделение на части)

Выполнено:
Возврат записанных данных
Возможность протестировать
Возврат без строки, вместо неё ошибка
Проверку на дубликат вводных данных
*/
type DatabaseExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type TutorHandler struct {
	DB DatabaseExecutor
}

func (h *TutorHandler) SaveTutor(c *fiber.Ctx) error {
	tutor := new(structs.Tutor)
	//Получение данных с запроса
	if err := c.BodyParser(tutor); err != nil {
		return c.Status(400).SendString("Неверный формат запрос")
	}
	//============================================================================
	//ВАЛИДАЦИЯ ДАННЫХ
	//============================================================================

	//Проверка на повторное регистрирование почты
	_, err := h.DB.Exec("SELECT email FROM tutor WHERE email = ?", tutor.Email)

	switch {
	//Проверка на отсутствие ввода имени
	case tutor.Name == "":
		return c.Status(400).SendString("Имя нужно обязательно.")
	//Проверка на отстутствие ввода почты
	case tutor.Email == "":
		return c.Status(400).SendString("Email нужен обязательно.")
	//Завершение проверки на повторное регистрирование почты
	case err != nil:
		return c.Status(400).SendString("Email уже занят")
	}
	//============================================================================

	_, err = h.DB.Exec("INSERT INTO tutor (name, email, expworktime, expectation, needcourses, tutorbefore) VALUES (?, ?, ?, ?, ?, ?)", tutor.Name, tutor.Email, tutor.ExpWorkTime, tutor.Expectation, tutor.NeedCourses, tutor.TutorBefore)
	if err != nil {
		return c.Status(500).SendString("Ошибка вставки данных в базу")
	}
	return c.Status(201).JSON(fiber.Map{
		"message": "Репетитор успешно добавлен:",
		"tutor": fiber.Map{
			"Имя":                    tutor.Name,
			"Почта":                  tutor.Email,
			"Опыт работы":            tutor.ExpWorkTime,
			"Ожидания от работы":     tutor.Expectation,
			"Нужны ли курсы":         tutor.NeedCourses,
			"Преподавал ли до этого": tutor.TutorBefore,
		},
	})
}
