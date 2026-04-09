package repository_test

import (
	"DemoFormTutor/Form/repository"
	"DemoFormTutor/Form/structs"
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

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

func TestPostgresTutorRepository_Save(t *testing.T) {
	t.Run("TestPostgresTutorRepository_Save", func(t *testing.T) {
		ctx := context.Background()

		tutor := &structs.Tutor{
			Name:        "Amogus",
			Email:       "blabla@bla.bla",
			ExpWorkTime: "Не хватает!!!",
			Expectation: "Деньга",
			NeedCourses: false,
			TutorBefore: true,
		}

		tutorGet, err := repo.Save(ctx, tutor)
		require.NoError(t, err)

		var expectedTrue bool
		err = db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM tutor WHERE id = $1)`, tutorGet.ID).Scan(&expectedTrue)
		require.NoError(t, err)
		require.True(t, expectedTrue)
	})
}
