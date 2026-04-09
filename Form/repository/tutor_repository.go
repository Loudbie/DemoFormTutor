package repository

import (
	"DemoFormTutor/Form/structs"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type TutorRepository interface {
	Save(ctx context.Context, tutor *structs.Tutor) (*structs.Tutor, error)
	FindById(ctx context.Context, id string) (*structs.Tutor, error)
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
func NewPostgresTutorRepositoryFromDB(db *sql.DB) TutorRepository {
	return &PostgresTutorRepository{db: db}
}

func (r *PostgresTutorRepository) Save(ctx context.Context, t *structs.Tutor) (*structs.Tutor, error) {
	const query = `INSERT INTO tutor 
    				(name, email, expworktime, expectation, needcourses, tutorbefore, createdat)
					VALUES ($1, $2, $3, $4, $5, $6, $7)
					RETURNING id`
	err := r.db.QueryRowContext(ctx, query,
		t.Name,
		t.Email,
		t.ExpWorkTime,
		t.Expectation,
		t.NeedCourses,
		t.TutorBefore,
		time.Now(),
	).Scan(&t.ID)
	t.CreatedAt = time.Now()
	if err != nil {
		return &structs.Tutor{}, fmt.Errorf("insert tutor: %w", err)
	}
	return t, nil
}

func (r *PostgresTutorRepository) FindById(ctx context.Context, id string) (*structs.Tutor, error) {
	const query = `SELECT	name,
       						email,
       						expworktime,
       						expectation,
       						needcourses,
       						tutorbefore,
       						createdat
				 	 FROM	tutor
		 			WHERE	id = $1`
	t := &structs.Tutor{}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("parse tutor: %w", err)
	}

	var timeFromDB sql.NullTime

	tutorRow := r.db.QueryRowContext(ctx, query, idInt)
	err = tutorRow.Scan(
		&t.Name,
		&t.Email,
		&t.ExpWorkTime,
		&t.Expectation,
		&t.NeedCourses,
		&t.TutorBefore,
		&timeFromDB,
	)
	if timeFromDB.Valid {
		t.CreatedAt = timeFromDB.Time
	}
	log.Println(t.Name, t.Email, t.ExpWorkTime, t.Expectation, t.NeedCourses, t.TutorBefore, t.CreatedAt, err)
	if err != nil {
		return &structs.Tutor{}, fmt.Errorf("find tutor by id: %w", err)
	}
	t.ID = idInt
	log.Println(t)
	return t, nil
}
