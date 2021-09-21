package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/graph/model"
	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/lti_utils"
	"github.com/adjust/rmq/v3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
	uuid "github.com/satori/go.uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	connectRetries int = 20
)

func UploadExam(minioClient *minio.Client, objectName string, fileReader io.Reader, fileSize int64, contentType string) error {
	bucketName := os.Getenv("MINIO_EXAM_BUCKET")

	info, err := minioClient.PutObject(
		context.Background(),
		bucketName,
		objectName,
		fileReader,
		fileSize,
		minio.PutObjectOptions{ContentType: contentType},
	)

	if err != nil {
		return err
	}
	log.Printf("Successfully uploaded %s of size %d with ContentType \"%s\"\n", objectName, info.Size, contentType)

	return nil
}

func InitDB(initialize bool) *gorm.DB {
	databaseConnectionString := fmt.Sprintf("host=%s port=5432 user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PASSWORD"))

	db, err := gorm.Open(
		postgres.Open(databaseConnectionString),
		&gorm.Config{},
	)

	for tries := 0; err != nil && tries <= connectRetries; tries++ {
		log.Println(err)
		db, err = gorm.Open(
			postgres.Open(databaseConnectionString),
			&gorm.Config{},
		)

		if tries >= connectRetries {
			log.Fatalln("reached maximium amount of connection tries to DB: ", connectRetries)
		}

		time.Sleep(1 * time.Second)
	}

	log.Print("connected successfully to the Database")

	deploymentEnv := os.Getenv("DEPLOYMENT_ENV")
	if deploymentEnv != "production" {
		log.Print("deployment environment: " + deploymentEnv)
		db.Debug()
	}

	if initialize {
		db.AutoMigrate(&model.Exam{})
		db.AutoMigrate(&lti_utils.LTIUserInfos{})
	}

	return db
}

func InitMinIO() *minio.Client {
	server := os.Getenv("MINIO_SERVER")
	port := os.Getenv("MINIO_PORT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	examBucket := os.Getenv("MINIO_EXAM_BUCKET")
	cacheBucket := os.Getenv("MINIO_CACHE_BUCKET")

	useSSL := false
	if os.Getenv("MINIO_SERVER_SSL") != "" {
		useSSLBool, err := strconv.ParseBool(os.Getenv("MINIO_SERVER_SSL"))
		if err != nil {
			log.Fatalln("MINIO_SERVER_SSL ", err)
		}
		useSSL = useSSLBool
	}

	// Initialize minio client object.
	minioClient, err := minio.New(server+":"+port, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	for tries := 0; err != nil && tries <= connectRetries; tries++ {
		log.Println(err)
		minioClient, err = minio.New(server+":"+port, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: useSSL,
		})

		time.Sleep(1 * time.Second)
	}
	if err != nil {
		log.Println("reached maximium amount of connection tries to S3 Storage Service: ", connectRetries)
		log.Fatalln(err)
	}

	setUpBucket(minioClient, examBucket)
	setUpBucket(minioClient, cacheBucket)

	// set the lifecycle policy for the cache bucket
	// this should auto delete the objects after a day
	config := lifecycle.NewConfiguration()
	config.Rules = []lifecycle.Rule{
		{
			ID:     cacheBucket,
			Status: "Enabled",
			Expiration: lifecycle.Expiration{
				Days: 1,
			},
		},
	}

	if err = minioClient.SetBucketLifecycle(context.Background(), cacheBucket, config); err != nil {
		log.Fatalln(err)
	}

	log.Println("MinIO client successfully set up!")
	return minioClient
}

func setUpBucket(minioClient *minio.Client, bucketName string) error {
	// set up the bucket to write the exams into
	contxt := context.Background()
	err := minioClient.MakeBucket(
		contxt,
		bucketName,
		minio.MakeBucketOptions{},
	)

	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(contxt, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			return err
		}
	} else {
		log.Printf("Successfully created the bucket \"%s\"\n", bucketName)
	}

	return nil
}

func InitRmq() rmq.Connection {
	// get job from queue
	errChan := make(chan error, 10)
	go logErrors(errChan)
	rmqClient, err := rmq.OpenConnection(
		os.Getenv("RMQ_QUEUE_NAME"),
		"tcp",
		os.Getenv("REDIS_CONNECTION_STRING"),
		1,
		errChan,
	)

	if err != nil {
		log.Fatalln(err)
	}
	return rmqClient
}

func logErrors(errChan <-chan error) {
	for err := range errChan {
		switch err := err.(type) {
		case *rmq.HeartbeatError:
			if err.Count == rmq.HeartbeatErrorLimit {
				log.Print("heartbeat error (limit): ", err)
			} else {
				log.Print("heartbeat error: ", err)
			}
		case *rmq.ConsumeError:
			log.Print("consume error: ", err)
		case *rmq.DeliveryError:
			log.Print("delivery error: ", err.Delivery, err)
		default:
			log.Print("other error: ", err)
		}
	}
}

type RMQMarkerTask struct {
	ExamUUID uuid.UUID `json:"examuuid"`
	Text     string    `json:"text"`
}
