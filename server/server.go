package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/graph"
	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/graph/generated"
	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/lti_utils"
	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/utils"
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
	// we initialize the db
	db := utils.InitDB(true)
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET_KEY")), nil)

	authHelper := lti_utils.AuthHelper{
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

		router.Get("/testlogin", func(w http.ResponseWriter, r *http.Request) {
			tmpl := template.Must(template.ParseFiles("server/dummylogin.html"))
			tmpl.Execute(w, nil)
		})
	}

	router.Get("/distributor/lti_config", lti_utils.LTIConfigHandler)

	if os.Getenv("DEPLOYMENT_ENV") != "testing" {
		router.Post("/distributor/lti_launch", authHelper.LTILaunch)
	} else {
		router.Post("/distributor/lti_launch", authHelper.DummyLTILaunch)
	}

	router.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))

		// Handle valid / invalid tokens. In this example, we use
		// the provided authenticator middleware, but you can write your
		// own very easily, look at the Authenticator method in jwtauth.go
		// and tweak it, its not scary.
		r.Use(jwtauth.Authenticator)

		r.Handle("/query", srv)
		r.Get("/adminlogin", authHelper.AdminLoginHandler)
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
