package repository_test

import (
	"DemoFormTutor/Form/repository"
	"DemoFormTutor/Form/structs"
	"context"
	"database/sql"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	db   *sql.DB
	repo repository.TutorRepository
)

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open("postgres", "user=postgres dbname=Demo sslmode=disable password=1234955 port=9090 host=localhost")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var maxID sql.NullInt64
	_ = db.QueryRow(`SELECT max(id) FROM tutor`).Scan(&maxID)

	repo = repository.NewPostgresTutorRepositoryFromDB(db)

	cleanUp := func() {
		if maxID.Valid {
			if _, err = db.Exec(`DELETE FROM tutor WHERE id > $1`, maxID.Int64); err != nil {
				log.Printf("cleanup delete: %v", err)
			}
		} else {
			if _, err = db.Exec(`DELETE FROM tutor`); err != nil {
				log.Printf("cleanup delete ALL: %v", err)
			}
		}

	}
	code := m.Run()
	cleanUp()
	os.Exit(code)
}

func TestPostgresTutorRepository(t *testing.T) {
	tutor := &structs.Tutor{
		ID:          13,
		Name:        "Amogus",
		Email:       "blabla@bla.bla",
		ExpWorkTime: "Не хватает!!!",
		Expectation: "Деньга",
		NeedCourses: false,
		TutorBefore: true,
		CreatedAt:   time.Now(),
	}
	t.Run("Save", func(t *testing.T) {
		ctx := context.Background()

		tutorGetPointer, err := repo.Save(ctx, tutor)
		require.NoError(t, err)

		tutor.ID = tutorGetPointer.ID
		assert.Equal(t, tutor, tutorGetPointer)
	})
	t.Run("Get", func(t *testing.T) {
		ctx := context.Background()
		tutorGetPointer, err := repo.FindById(ctx, strconv.Itoa(tutor.ID))
		require.NoError(t, err)

		tutor.CreatedAt = tutorGetPointer.CreatedAt
		assert.Equal(t, tutor, tutorGetPointer)
	})
	t.Run("Gets", func(t *testing.T) {
		ctx := context.Background()

		_, err := repo.GetAll(ctx, "status=false")
		require.NoError(t, err)

		_, err = repo.GetAll(ctx, "status=true")
		require.NoError(t, err)

		_, err = repo.GetAll(ctx, "date")
		require.NoError(t, err)

		_, err = repo.GetAll(ctx, "all")
	})
}
