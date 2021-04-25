package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/graph/generated"
	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/graph/model"
)

func (r *mutationResolver) CreateExam(ctx context.Context, input model.NewExam) (*model.Exam, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateLecturer(ctx context.Context, input model.NewLecturer) (*model.Lecturer, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) LinkExamAndLecturer(ctx context.Context, input *model.LinkInput) (*model.Exam, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Exams(ctx context.Context) ([]*model.Exam, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
