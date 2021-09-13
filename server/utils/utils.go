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
	"github.com/adjust/rmq/v3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

func InitDB() *gorm.DB {
	databaseConnectionString := fmt.Sprintf("host=%s port=5432 user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PASSWORD"))

	db, err := gorm.Open(
		postgres.Open(databaseConnectionString),
		&gorm.Config{},
	)

	for err != nil {
		log.Println(err)
		db, err = gorm.Open(
			postgres.Open(databaseConnectionString),
			&gorm.Config{},
		)
		time.Sleep(1 * time.Second)
	}

	log.Print("connected successfully to the Database")

	deploymentEnv := os.Getenv("DEPLOYMENT_ENV")
	if deploymentEnv != "production" {
		log.Print("deployment environment: " + deploymentEnv)
		db.Debug()
	}

	db.AutoMigrate(&model.Exam{})
	db.AutoMigrate(&LTIUserInfos{})

	return db
}

func InitMinIO() *minio.Client {
	server := os.Getenv("MINIO_SERVER")
	port := os.Getenv("MINIO_PORT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	bucketName := os.Getenv("MINIO_EXAM_BUCKET")

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
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("MinIO client successfully set up!")

	// set up the bucket to write the exams into
	contxt := context.Background()
	err = minioClient.MakeBucket(
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
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created the bucket \"%s\"\n", bucketName)
	}

	return minioClient
}

func InitRmq() rmq.Connection {
	// get job from queue
	errChan := make(chan error, 10)
	go logErrors(errChan)
	rmqClient, err := rmq.OpenConnection("exam-message-passing", "tcp", "altklausur_ausleihe-redis:6379", 1, errChan)
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
