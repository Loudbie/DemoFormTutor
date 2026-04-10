package repository

import (
	"DemoFormTutor/Form/structs"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type TutorRepository interface {
	Save(ctx context.Context, tutor *structs.Tutor) (*structs.Tutor, error)
	FindById(ctx context.Context, id string) (*structs.Tutor, error)
	GetAll(ctx context.Context, sort string) ([]structs.Tutor, error)
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
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, err
		default:
			return &structs.Tutor{}, fmt.Errorf("find tutor by id: %w", err)
		}
	}
	t.ID = idInt
	log.Println(t)
	return t, nil
}

func (r *PostgresTutorRepository) GetAll(ctx context.Context, sort string) ([]structs.Tutor, error) {
	const (
		query = `SELECT 	id,
       						name,
       						email,
       						expworktime,
       						expectation,
       						needcourses,
       						tutorbefore,
       						createdat,
       						status
					FROM 	tutor`
		queryStatus = `SELECT 	id,
       						name,
       						email,
       						expworktime,
       						expectation,
       						needcourses,
       						tutorbefore,
       						createdat,
       						status
					FROM 	tutor
					WHERE	status = $1`
		queryDate = `SELECT id,
       						name,
       						email,
       						expworktime,
       						expectation,
       						needcourses,
       						tutorbefore,
       						createdat,
       						status
					FROM 	tutor
					ORDER BY createdat DESC`
	)
	var tutors []structs.Tutor
	var err error
	var rows *sql.Rows

	switch sort {
	case "status=false", "status=true":
		sort = strings.TrimPrefix(sort, "status=")
		parseBool, err := strconv.ParseBool(sort)
		if err != nil {
			return nil, fmt.Errorf("str to boolean Parse false: %w", err)
		}
		rows, err = r.db.QueryContext(ctx, queryStatus, parseBool)
	case "date":
		rows, err = r.db.QueryContext(ctx, queryDate)
	case "all":
		rows, err = r.db.QueryContext(ctx, query)
	default:
		return []structs.Tutor{}, errors.New("invalid sort parameter")
	}

	if err != nil {
		return tutors, fmt.Errorf("get tutors: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t structs.Tutor
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Email,
			&t.ExpWorkTime,
			&t.Expectation,
			&t.NeedCourses,
			&t.TutorBefore,
			&t.CreatedAt,
			&t.Status,
		)
		if err != nil {
			return tutors, fmt.Errorf("scan to get tutors: %w", err)
		}
		tutors = append(tutors, t)
	}
	if err = rows.Err(); err != nil {
		return tutors, fmt.Errorf("rows error: %w", err)
	}
	return tutors, nil
}
