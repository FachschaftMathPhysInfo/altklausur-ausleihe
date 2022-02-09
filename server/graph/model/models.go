// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"strconv"
	"strings"

	"gorm.io/gorm"

	"github.com/99designs/gqlgen/graphql"
	uuid "github.com/satori/go.uuid"
)

type Exam struct {
	UUID          uuid.UUID `gorm:"type:uuid;primary_key;" json:"ID,omitempty"`
	Subject       string    `json:"subject"`
	Hash          string    `json:"hash"`
	ModuleName    string    `json:"moduleName"`
	ModuleAltName *string   `json:"moduleAltName"`
	Year          *int      `json:"year"`
	Examiners     *string   `json:"examiners"`
	Semester      *string   `json:"semester"`
}

// BeforeCreate will set a UUID rather than numeric ID.
// taken from
// https://github.com/FachschaftMathPhysInfo/ostseee/blob/master/server/go/model_base.go
func (exam *Exam) BeforeCreate(db *gorm.DB) error {
	// Check if the UUID is already set (i.e. db.Save(...))
	if uuid.Equal(exam.UUID, uuid.Nil) {
		uuid := uuid.NewV4()
		exam.UUID = uuid
	}
	return nil
}

// ToFilename returns a normalized string version of the metadata for this exam
func (exam *Exam) ToFilename() string {
	filename := exam.ModuleName

	// add year and semester if both are present
	if exam.Year != nil && exam.Semester != nil {
		filename += "_" + *exam.Semester + "_" + strconv.Itoa(*exam.Year)
	}

	filename = strings.ReplaceAll(strings.ToLower(filename), " ", "_")

	return filename
}

type NewExam struct {
	Subject       string         `json:"subject"`
	ModuleName    string         `json:"moduleName"`
	File          graphql.Upload `json:"file"`
	ModuleAltName *string        `json:"moduleAltName"`
	Year          *int           `json:"year"`
	Examiners     *string        `json:"examiners"`
	Semester      *string        `json:"semester"`
}
