package handlers

import (
	"layers/Form/database"
	"layers/Form/structs"

	"github.com/gofiber/fiber/v2"
)

func SaveTutor(c *fiber.Ctx) error {
	tutor := new(structs.Tutor)
	if err := c.BodyParser(tutor); err != nil {
		return c.Status(400).SendString("Неверный формат запрос")
	}
	_, err := database.DB.Exec("INSERT INTO tutor (name, email, expworktime, expectation, needcourses, tutorbefore) VALUES ($1, $2, $3, $4, $5, $6)", tutor.Name, tutor.Email, tutor.ExpWorkTime, tutor.Expectation, tutor.NeedCourses, tutor.TutorBefore)
	if err != nil {
		return c.Status(500).SendString("Ошибка вставки данных в базу")
	}
	return c.Status(201).SendString("Репетитор успешно добавлен")
}
