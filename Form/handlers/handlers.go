package handlers

import (
	"DemoFormTutor/Form/database"
	"DemoFormTutor/Form/structs"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

/*
TODO Возможность протестировать
TODO Переиспользование кода (разделение на части)
TODO Возврат без строки, вместо неё нужна ошибка
TODO Проверку на дубликат вводных данных

Выполнено:
Возврат записанных данных
*/

func SaveTutor(c *fiber.Ctx) error {
	tutor := new(structs.Tutor)
	if err := c.BodyParser(tutor); err != nil {
		return c.Status(400).SendString("Неверный формат запрос")
	}
	_, err := database.Connect().Exec("INSERT INTO tutor (name, email, expworktime, expectation, needcourses, tutorbefore) VALUES ($1, $2, $3, $4, $5, $6)", tutor.Name, tutor.Email, tutor.ExpWorkTime, tutor.Expectation, tutor.NeedCourses, tutor.TutorBefore)
	if err != nil {
		return c.Status(500).SendString("Ошибка вставки данных в базу")
	}
	return c.Status(201).SendString(
		"Репетитор успешно добавлен:" +
			"\nИмя: " + tutor.Name +
			"\nEmail:" + tutor.Email +
			"\nОпыт работы: " + tutor.ExpWorkTime +
			"\nОжидания от работы: " + tutor.Expectation +
			"\nНеобходимость курсов: " + strconv.FormatBool(tutor.NeedCourses) +
			"\nРаботал ли до этого: " + strconv.FormatBool(tutor.TutorBefore),
	)
}
