package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/graph"
	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/graph/generated"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/graph/model"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const defaultPort = "8081"

func initDB() *gorm.DB {
	databaseConnectionString := os.Getenv("DB_CONNECTION_STRING")

	if databaseConnectionString != "" {
		log.Println(databaseConnectionString)
	} else {
		log.Fatal("DB_CONNECTION_STRING empty!")
	}
	err := fmt.Errorf("initial connect failed")

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

	return db
}

func initMinIO() *minio.Client {
	server := os.Getenv("MINIO_SERVER")
	port := os.Getenv("MINIO_PORT")
	accessKeyID := os.Getenv("MINIO_ROOT_USER")
	secretAccessKey := os.Getenv("MINIO_ROOT_PASSWORD")
	bucketName := os.Getenv("MINIO_BUCKET_NAME")
	useSSL := false

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
		log.Printf("Successfully created %s\n", bucketName)
	}

	return minioClient
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	var mb int64 = 1 << 20

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &graph.Resolver{
					DB:          initDB(),
					MinIOClient: initMinIO(),
				},
			},
		),
	)
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{
		MaxMemory:     32 * mb,
		MaxUploadSize: 50 * mb,
	})
	srv.Use(extension.Introspection{})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
