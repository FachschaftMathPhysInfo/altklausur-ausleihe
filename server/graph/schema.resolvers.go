package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"log"

	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/graph/generated"
	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/graph/model"
	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/utils"
)

func (r *examResolver) UUID(ctx context.Context, obj *model.Exam) (string, error) {
	return obj.UUID.String(), nil
}

func (r *mutationResolver) CreateExam(ctx context.Context, input model.NewExam) (*model.Exam, error) {
	// map the GraphQL input to the Model
	exam := model.Exam{
		Subject:       input.Subject,
		ModuleName:    input.ModuleName,
		ModuleAltName: input.ModuleAltName,
		Year:          input.Year,
		Examiners:     input.Examiners,
		Semester:      input.Semester,
	}

	// create the exam in the database
	r.DB.Create(&exam)
	if r.DB.Error != nil {
		return nil, r.DB.Error
	}

	// upload the file to the storage server
	// this assumes that the database sets an exams' UUID
	uploadErr := utils.UploadExam(
		r.MinIOClient,
		exam.UUID.String(),
		input.File.File,
		input.File.Size,
		input.File.ContentType)

	if uploadErr != nil {
		// do we need to reroll the inserted db entry?
		return nil, uploadErr
	}

	return &exam, nil
}

func (r *queryResolver) Exams(ctx context.Context) ([]*model.Exam, error) {
	var exam []*model.Exam
	r.DB.Find(&exam)

	if r.DB.Error != nil {
		return nil, r.DB.Error
	}

	return exam, nil
}

func (r *mutationResolver) RequestMarkedExam(ctx context.Context, uuid string) (*string, error) {
	log.Println(uuid)

	tagQueue, err := r.RmqClient.OpenQueue("tag-queue")
	if err != nil {
		return nil, err
	}

	tagQueue.Publish(uuid)
	return &uuid, nil
}

// Exam returns generated.ExamResolver implementation.
func (r *Resolver) Exam() generated.ExamResolver { return &examResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type examResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
