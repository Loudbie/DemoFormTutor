package structs

type Tutor struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	ExpWorkTime string `json:"expworktime"`
	Expectation string `json:"expectation"`
	NeedCourses bool   `json:"needcourses"`
	TutorBefore bool   `json:"tutorbefore"`
}
