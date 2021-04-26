// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Exam struct {
	ID       int     `json:"ID"`
	Subject  string  `json:"subject"`
	URL      string  `json:"url"`
	Semester *string `json:"semester"`
}

type Lecturer struct {
	ID      int    `json:"ID"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

type LinkInput struct {
	LecturerID int `json:"lecturer_id"`
	ExamID     int `json:"exam_id"`
}

type NewExam struct {
	Subject  string  `json:"subject"`
	Semester *string `json:"semester"`
}

type NewLecturer struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}
