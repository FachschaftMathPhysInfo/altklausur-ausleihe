// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Exam struct {
	ID         int     `json:"ID"`
	Subject    string  `json:"subject"`
	ModuleName string  `json:"moduleName"`
	URL        string  `json:"url"`
	Year       *int    `json:"year"`
	Examiners  *string `json:"examiners"`
	Semester   *string `json:"semester"`
}

type NewExam struct {
	Subject    string  `json:"subject"`
	ModuleName string  `json:"moduleName"`
	Year       *int    `json:"year"`
	Examiners  *string `json:"examiners"`
	Semester   *string `json:"semester"`
}
