package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
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

	// db.AutoMigrate(&model.Exam{})

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
	if err != nil {
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
