package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"net/url"
	"os"
	"time"

	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/graph/generated"
	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/graph/model"
	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/utils"
	"github.com/gabriel-vasile/mimetype"
	minio "github.com/minio/minio-go/v7"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func (r *examResolver) UUID(ctx context.Context, obj *model.Exam) (string, error) {
	return obj.UUID.String(), nil
}

func (r *mutationResolver) CreateExam(ctx context.Context, input model.NewExam) (*model.Exam, error) {
	if input.Semester != nil && !(*input.Semester == "SoSe" || *input.Semester == "WiSe") {
		return nil, fmt.Errorf("Input \"%s\" is not a valid input for field input.Semester", *input.Semester)
	}

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

	// check file size
	if input.File.Size < 512 {
		// TODO: implement DB rollback here!
		return nil, fmt.Errorf("File is not valid: size of %d too small!", input.File.Size)
	}

	// check file MIME type
	// Only the first 512 bytes are used to sniff the content type.
	fileReader := bufio.NewReader(input.File.File)
	buffer, err := fileReader.Peek(512)
	if err != nil {
		return nil, err
	}

	mtype := mimetype.Detect(buffer)
	allowedMIMETypes := []string{"application/pdf"}
	if !mimetype.EqualsAny(mtype.String(), allowedMIMETypes...) {
		// TODO: implement DB rollback here!
		return nil, fmt.Errorf("File is not valid: mimetype \"%s\" forbidden!", mtype.String())
	}

	// upload the file to the storage server
	// this assumes that the database sets an exams' UUID
	uploadErr := utils.UploadExam(
		r.MinIOClient,
		exam.UUID.String(),
		fileReader,
		input.File.Size,
		input.File.ContentType)

	if uploadErr != nil {
		// do we need to reroll the inserted db entry?
		return nil, uploadErr
	}

	return &exam, nil
}

func (r *mutationResolver) RequestMarkedExam(ctx context.Context, stringUUID string) (*string, error) {
	// check if we got a valid uuid and also prepare the DB search
	realUUID, err := uuid.FromString(stringUUID)
	if err != nil {
		return nil, err
	}

	// see if there is an registered exam for this uuid
	var exam model.Exam
	dbErr := r.DB.First(&exam, realUUID).Error
	if errors.Is(dbErr, gorm.ErrRecordNotFound) {
		return nil, dbErr
	}

	// try to find the entry in cache
	_, e := r.MinIOClient.StatObject(context.Background(), os.Getenv("MINIO_CACHE_BUCKET"), stringUUID, minio.GetObjectOptions{})
	if e != nil {
		errResponse := minio.ToErrorResponse(e)
		if errResponse.Code != "NoSuchKey" {
			return nil, e
		}
	} else {
		return &stringUUID, nil
	}

	tagQueue, err := r.RmqClient.OpenQueue("tag-queue")
	if err != nil {
		return nil, err
	}

	rawTask := utils.RMQMarkerTask{
		ExamUUID: realUUID,
		Text:     "test",
	}

	task, err := json.Marshal(rawTask)
	if err != nil {
		return nil, err
	}

	// add the job to the workers
	if err := tagQueue.Publish(string(task)); err != nil {
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

	return &stringUUID, nil
}

func (r *queryResolver) Exams(ctx context.Context) ([]*model.Exam, error) {
	var exam []*model.Exam
	r.DB.Find(&exam)

	if r.DB.Error != nil {
		return nil, r.DB.Error
	}

	return exam, nil
}

func (r *queryResolver) GetExam(ctx context.Context, stringUUID string) (*model.PresignedReturn, error) {
	// check if we got a valid uuid and also prepare the DB search
	realUUID, err := uuid.FromString(stringUUID)
	if err != nil {
		return nil, err
	}

	// see if there is an registered exam for this uuid
	var exam model.Exam
	dbErr := r.DB.First(&exam, realUUID).Error
	if errors.Is(dbErr, gorm.ErrRecordNotFound) {
		return nil, dbErr
	}

	// try to find the entry in cache
	_, e := r.MinIOClient.StatObject(context.Background(), os.Getenv("MINIO_CACHE_BUCKET"), stringUUID, minio.GetObjectOptions{})
	if e != nil {
		errResponse := minio.ToErrorResponse(e)
		if errResponse.Code != "NoSuchKey" {
			return nil, e
		}
		return nil, nil
	}

	// Set request parameters for content-disposition.
	// Beware of this issue: https://github.com/minio/minio/issues/7936
	reqParams := make(url.Values)

	// Generates a presigned url to view the pdf which expires in 15 min.
	presignedViewURL, err := r.MinIOClient.PresignedGetObject(context.Background(), os.Getenv("MINIO_CACHE_BUCKET"), stringUUID, 15*time.Minute, reqParams)
	if err != nil {
		return nil, err
	}

	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", stringUUID))
	// Generates a presigned url to download the pdf which expires in 15 min.
	presignedDownloadURL, err := r.MinIOClient.PresignedGetObject(context.Background(), os.Getenv("MINIO_CACHE_BUCKET"), stringUUID, 15*time.Minute, reqParams)
	if err != nil {
		return nil, err
	}

	return &model.PresignedReturn{
			ViewURL:     presignedViewURL.String(),
			DownloadURL: presignedDownloadURL.String(),
		},
		nil
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
