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
	var idTest string
	tutor := &structs.Tutor{
		Name:        "Amogus",
		Email:       "blabla@bla.bla",
		ExpWorkTime: "Не хватает!!!",
		Expectation: "Деньга",
		NeedCourses: false,
		TutorBefore: true,
	}
	t.Run("Save", func(t *testing.T) {
		ctx := context.Background()

		tutorGetPointer, err := repo.Save(ctx, tutor)
		require.NoError(t, err)
		var tutorInDB structs.Tutor
		row := db.QueryRow(`SELECT id, name, email, expworktime, expectation, needcourses, tutorbefore
									FROM tutor WHERE id = $1`, tutorGetPointer.ID)

		err = row.Scan(&tutorInDB.ID, &tutorInDB.Name,
			&tutorInDB.Email, &tutorInDB.ExpWorkTime,
			&tutorInDB.Expectation, &tutorInDB.NeedCourses, &tutorInDB.TutorBefore)
		require.NoError(t, err)

		tutorGet := *tutorGetPointer
		assert.Equal(t, tutorInDB, tutorGet)
		idTest = strconv.Itoa(tutorInDB.ID)
	})
	t.Run("Get", func(t *testing.T) {
		ctx := context.Background()
		idGet, err := strconv.Atoi(idTest)
		require.NoError(t, err)
		tutor.ID = idGet

		tutorGetPointer, err := repo.FindById(ctx, idTest)
		require.NoError(t, err)
		var tutorInDB structs.Tutor
		row := db.QueryRow(`SELECT	name,
       						email,
       						expworktime,
       						expectation,
       						needcourses,
       						tutorbefore
				 	 FROM	tutor
		 			WHERE	id = $1`, tutor.ID)
		err = row.Scan(
			&tutorInDB.Name,
			&tutorInDB.Email,
			&tutorInDB.ExpWorkTime,
			&tutorInDB.Expectation,
			&tutorInDB.NeedCourses,
			&tutorInDB.TutorBefore,
		)
		require.NoError(t, err)
		tutorGet := *tutorGetPointer
		tutorGet.ID = tutorInDB.ID
		assert.Equal(t, tutorInDB, tutorGet)
	})
}
