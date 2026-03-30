package main

import (
	"DemoFormTutor/Form/database"
	"DemoFormTutor/Form/routes"
	"DemoFormTutor/Form/viperConfig"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

/*
1. Сделать ручку для сохранения заявки
(поля формы: опыт работы, ожидания от работы,
занимался ли репетиторством до этого, нужны ли обучающие курсы).
2. Сделать ручку для получения всех заявок с возможностью фильтрации
(по статусу, по дате отправки).
3. При подтверждении заявки нужно отправить пользователю
на email сообщение о том, что его приняли.
4. Сделать Тесты
*/
func main() {
	viperConfig.CheckSetConfig()

	database.Connect()
	app := fiber.New(fiber.Config{
		Prefork: false,
	})

	routes.RegisterRoutes(app)

	log.Fatal(app.Listen(viper.GetString("port")))
}
