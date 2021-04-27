package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/graph/generated"
	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/graph/model"
)

func (r *mutationResolver) CreateExam(ctx context.Context, input model.NewExam) (*model.Exam, error) {
	exam := model.Exam{
		Subject:       input.Subject,
		ModuleName:    input.ModuleName,
		ModuleAltName: input.ModuleAltName,
		Year:          input.Year,
		Examiners:     input.Examiners,
		Semester:      input.Semester,
	}

	r.DB.Create(&exam)

	if r.DB.Error != nil {
		return nil, r.DB.Error
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

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
