package structs

import "time"

type Tutor struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	ExpWorkTime string    `json:"expworktime"`
	Expectation string    `json:"expectation"`
	NeedCourses bool      `json:"needcourses"`
	TutorBefore bool      `json:"tutorbefore"`
	CreatedAt   time.Time `json:"createdat"`
	Status      bool      `json:"status"`
}
