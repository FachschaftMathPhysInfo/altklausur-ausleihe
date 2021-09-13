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
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
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
	db := utils.InitDB()
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET_KEY")), nil)

	ltiConnector := utils.LTIConnector{
		DB:        db,
		TokenAuth: tokenAuth,
	}

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &graph.Resolver{
					DB:          db,
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

	if os.Getenv("DEPLOYMENT_ENV") == "testing" {
		router.Handle("/", playground.Handler("GraphQL playground", "/query"))
		log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	}

	router.Get("/distributor/lti_config", utils.LTIConfigHandler)
	router.Post("/distributor/lti_launch", ltiConnector.LTILaunch)

	router.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))

		// Handle valid / invalid tokens. In this example, we use
		// the provided authenticator middleware, but you can write your
		// own very easily, look at the Authenticator method in jwtauth.go
		// and tweak it, its not scary.
		r.Use(jwtauth.Authenticator)

		r.Handle("/query", srv)
	})

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
