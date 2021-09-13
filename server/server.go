package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/graph"
	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/graph/generated"
	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/utils"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
)

const defaultPort = "8081"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()
	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	}).Handler)

	var mb int64 = 1 << 20

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &graph.Resolver{
					DB:          utils.InitDB(),
					MinIOClient: utils.InitMinIO(),
					RmqClient:   utils.InitRmq(),
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

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Get("/distributor/lti_config", utils.LTIConfigHandler)
	router.Post("/distributor/lti_launch", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	router.Handle("/query", srv)
	fmt.Print(
		"==========================================\n",
		"Started the backend listening on Port "+port+"\n",
		"==========================================\n",
	)
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		panic(err)
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
