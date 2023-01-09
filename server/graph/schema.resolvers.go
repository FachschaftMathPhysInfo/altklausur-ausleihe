package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net/url"
	"os"
	"time"

	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/graph/generated"
	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/graph/model"
	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/lti_utils"
	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/utils"
	"github.com/dustin/go-humanize"
	"github.com/gabriel-vasile/mimetype"
	minio "github.com/minio/minio-go/v7"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// UUID is the resolver for the UUID field.
func (r *examResolver) UUID(ctx context.Context, obj *model.Exam) (string, error) {
	return obj.UUID.String(), nil
}

// CreateExam is the resolver for the createExam field.
func (r *mutationResolver) CreateExam(ctx context.Context, input model.NewExam) (*model.Exam, error) {
	user, err := lti_utils.GetUserInfosFromContext(&ctx)
	if err != nil {
		return nil, err
	}

	// check if the user is an admin
	if !user.IsAdmin {
		return nil, fmt.Errorf("You are not an Admin lol, nice try!" +
			" Please hand in your exam via mail to fachschaft@mathphys.info ..." +
			" BTW, since you seem to read code, wanna contribute? ;)",
		)
	}

	if input.Semester != nil && !(*input.Semester == "SoSe" || *input.Semester == "WiSe") {
		return nil, fmt.Errorf("Input \"%s\" is not a valid input for field input.Semester", *input.Semester)
	}

	// we use a TeeReader to hash and copy to the buffer at the same time
	fileBuf := &bytes.Buffer{}
	tee := io.TeeReader(input.File.File, fileBuf)

	// generate the hash of the input file
	fileHash := sha256.New()
	if _, err := io.Copy(fileHash, tee); err != nil {
		log.Fatal(err)
	}
	encodedHash := hex.EncodeToString(fileHash.Sum(nil))

	dbExam := model.Exam{}
	// check if the exam already exists
	r.DB.Where("hash = ?", encodedHash).Find(&dbExam)

	// map the GraphQL input to the Model
	exam := model.Exam{
		UUID:          dbExam.UUID,
		Subject:       input.Subject,
		ModuleName:    input.ModuleName,
		ModuleAltName: input.ModuleAltName,
		Year:          input.Year,
		Examiners:     input.Examiners,
		Semester:      input.Semester,
		Hash:          encodedHash,
	}

	//
	// An existing exam was found
	//
	if !uuid.Equal(dbExam.UUID, uuid.Nil) {
		// update the exam in the database and abort
		r.DB.Save(&exam)
		if r.DB.Error != nil {
			return nil, r.DB.Error
		}
		return &exam, nil
	}

	//
	// no existing exam was found
	//

	// create the exam in the database
	r.DB.Create(&exam)
	if r.DB.Error != nil {
		return nil, r.DB.Error
	}

	// check file size
	if input.File.Size < 512 && fileBuf.Len() < 512 {
		// TODO: implement DB rollback here!
		return nil, fmt.Errorf("File is not valid: size of %s too small, buffer size %s", humanize.Bytes(uint64(input.File.Size)), humanize.Bytes(uint64(fileBuf.Len())))
	}

	// check file MIME type
	// Only the first 512 bytes are used to sniff the content type.
	fileReader := bufio.NewReader(fileBuf)
	buffer, err := fileReader.Peek(512)
	if err != nil {
		return nil, err
	}

	mtype := mimetype.Detect(buffer)
	allowedMIMETypes := []string{"application/pdf"}
	if !mimetype.EqualsAny(mtype.String(), allowedMIMETypes...) {
		// TODO: implement DB rollback here!
		return nil, fmt.Errorf("File is not valid: mimetype \"%s\" forbidden", mtype.String())
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

	// update the TotalExams metric
	utils.UpdateTotalExamsMetric(r.DB)

	return &exam, nil
}

// RequestMarkedExam is the resolver for the requestMarkedExam field.
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

	// get all the user infos
	userInfos, err := lti_utils.GetUserInfosFromContext(&ctx)
	if err != nil {
		return nil, err
	}

	// try to find the entry in cache
	_, e := r.MinIOClient.StatObject(
		context.Background(),
		os.Getenv("MINIO_CACHE_BUCKET"),
		utils.GetExamCachePath(userInfos.ID, realUUID),
		minio.GetObjectOptions{})

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

	// create the task for the exam marker
	task, err := json.Marshal(
		utils.RMQMarkerTask{
			ExamUUID:     realUUID,
			UserID:       userInfos.ID,
			TextLeft:     userInfos.PersonFamilyName + " - " + userInfos.PersonFamilyName,
			TextDiagonal: userInfos.PersonPrimaryEmail,
			SubmitTime:   time.Now(),
		},
	)

	if err != nil {
		return nil, err
	}

	// communicate it to the marker
	if err := tagQueue.Publish(string(task)); err != nil {
		return nil, err
	}

	utils.ExamsMarkedMetric.Inc()

	return &stringUUID, nil
}

// Exams is the resolver for the exams field.
func (r *queryResolver) Exams(ctx context.Context) ([]*model.Exam, error) {
	var exam []*model.Exam
	r.DB.Find(&exam)

	if r.DB.Error != nil {
		return nil, r.DB.Error
	}

	return exam, nil
}

// GetExam is the resolver for the getExam field.
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

	// get all the user infos
	userInfos, err := lti_utils.GetUserInfosFromContext(&ctx)
	if err != nil {
		return nil, err
	}

	// try to find the entry in cache
	objectInfo, e := r.MinIOClient.StatObject(
		context.Background(),
		os.Getenv("MINIO_CACHE_BUCKET"),
		utils.GetExamCachePath(userInfos.ID, realUUID),
		minio.GetObjectOptions{})

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
	presignedViewURL, err := r.MinIOClient.PresignedGetObject(
		context.Background(),
		os.Getenv("MINIO_CACHE_BUCKET"),
		utils.GetExamCachePath(userInfos.ID, realUUID),
		15*time.Minute,
		reqParams)

	if err != nil {
		return nil, err
	}

	filename := exam.ToFilename()

	if extensions, err := mime.ExtensionsByType(objectInfo.ContentType); extensions != nil && err == nil {
		filename += extensions[0]
	}

	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	// Generates a presigned url to download the pdf which expires in 15 min.
	presignedDownloadURL, err := r.MinIOClient.PresignedGetObject(
		context.Background(),
		os.Getenv("MINIO_CACHE_BUCKET"),
		utils.GetExamCachePath(userInfos.ID, realUUID),
		15*time.Minute,
		reqParams)

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
