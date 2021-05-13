package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"os"

	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/graph/generated"
	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/graph/model"
	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/utils"
	"github.com/minio/minio-go/v7"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
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

func (r *mutationResolver) RequestMarkedExam(ctx context.Context, uuidstring string) (*string, error) {
	// try to find the entry in cache
	_, e := r.MinIOClient.StatObject(context.Background(), os.Getenv("MINIO_CACHE_BUCKET"), uuidstring, minio.GetObjectOptions{})
	if e != nil {
		errResponse := minio.ToErrorResponse(e)
		if errResponse.Code != "NoSuchKey" {
			return nil, e
		}
	} else {
		return &uuidstring, nil
	}

	// check if we got a valid uuid and also prepare the DB search
	realUUID, err := uuid.FromString(uuidstring)
	if err != nil {
		return nil, err
	}

	// see if there is an registered exam for this uuid
	var exam model.Exam
	dbErr := r.DB.First(&exam, realUUID).Error
	if errors.Is(dbErr, gorm.ErrRecordNotFound) {
		return nil, dbErr
	}

	tagQueue, err := r.RmqClient.OpenQueue("tag-queue")
	if err != nil {
		return nil, err
	}

	// add the job to the workers
	if err := tagQueue.Publish(uuidstring); err != nil {
		return nil, err
	}

	//	// TODO(chris): is checking for the timeout even necessary?!
	//	timeout := 10
	//	for i := 0; i < timeout; i++ {
	//		// try to find the entry in cache
	//		_, e := r.MinIOClient.StatObject(context.Background(), os.Getenv("MINIO_CACHE_BUCKET"), uuidstring, minio.GetObjectOptions{})
	//		if e != nil {
	//			errResponse := minio.ToErrorResponse(e)
	//			if errResponse.Code != "NoSuchKey" {
	//				return nil, e
	//			}
	//		} else {
	//			return &uuidstring, nil
	//		}
	//		time.Sleep(500 * time.Millisecond)
	//	}
	//	return nil, fmt.Errorf("Timeout reached while marking exam \"%s\"!", uuidstring)

	return &uuidstring, nil
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
