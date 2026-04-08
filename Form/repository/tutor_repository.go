package repository

import (
	"DemoFormTutor/Form/structs"
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type TutorRepository interface {
	Save(ctx context.Context, tutor *structs.Tutor) (*structs.Tutor, error)
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

func (r *PostgresTutorRepository) Save(ctx context.Context, t *structs.Tutor) (*structs.Tutor, error) {
	const query = `INSERT INTO tutor 
    				(name, email, expworktime, expectation, needcourses, tutorbefore)
					VALUES ($1, $2, $3, $4, $5, $6)
					RETURNING id`
	err := r.db.QueryRowContext(ctx, query,
		t.Name,
		t.Email,
		t.ExpWorkTime,
		t.Expectation,
		t.NeedCourses,
		t.TutorBefore,
	).Scan(&t.ID)
	if err != nil {
		return &structs.Tutor{}, fmt.Errorf("insert tutor: %w", err)
	}
	return t, nil
}
