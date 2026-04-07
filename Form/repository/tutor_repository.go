package repository

import (
	"DemoFormTutor/Form/structs"
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

type TutorRepository interface {
	Save(ctx *fiber.Ctx, tutor *structs.Tutor) error
}
type PostgresTutorRepository struct {
	db *sql.DB
}

func NewPostgresTutorRepository(connStr string) (TutorRepository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return &PostgresTutorRepository{db: db}, nil
}

func (r *PostgresTutorRepository) Save(ctx *fiber.Ctx, t *structs.Tutor) error {
	const query = `INSERT INTO tutor 
    				(name, email, expworktime, expectation, needcourses, tutorbefore)
					VALUES ($1, $2, $3, $4, $5, $6)
					RETURNING id`
	_, err := r.db.Exec(query,
		t.Name,
		t.Email,
		t.ExpWorkTime,
		t.Expectation,
		t.NeedCourses,
		t.TutorBefore,
	)
	if err != nil {
		return ctx.Status(500).SendString("Ошибка вставки данных в базу")
	}
	return nil
}
