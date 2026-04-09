package main

import (
	"DemoFormTutor/Form/handlers"
	"DemoFormTutor/Form/repository"
	"DemoFormTutor/Form/routes"
	"DemoFormTutor/Form/usecase"
	"DemoFormTutor/Form/viperConfig"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

//TODO 2."Добавить возможность получить обратно данные по заявке"

//TODO 3. Сделать ручку для получения всех заявок с возможностью фильтрации
//TODO (по статусу, по дате отправки).

//TODO 4. При подтверждении заявки нужно отправить пользователю
//TODO на email сообщение о том, что его приняли.

//TODO 5. Сделать Тесты
/*
Выполнено:
1. Сделать ручку для сохранения заявки
(поля формы: опыт работы, ожидания от работы,
занимался ли репетиторством до этого, нужны ли обучающие курсы).
*/
func main() {
	if err := viperConfig.CheckSetConfig(); err != nil {
		log.Fatalln(err)
	}
	app := fiber.New(fiber.Config{
		Prefork: false,
	})
	connstr := viper.GetString("db")
	tutorRepo, err := repository.NewPostgresTutorRepository(connstr)
	if err != nil {
		log.Fatalf("Cannot connect to db: %v", err)
	}

	tutorUC := usecase.NewTutorUseCase(tutorRepo)
	handler := handlers.NewTutorHandler(tutorUC)

	routes.RegisterRoutes(app, handler)

	log.Fatal(app.Listen(viper.GetString("port")))
}
