package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/motemen/example-gqlgen-dataloader/db"
	"github.com/motemen/example-gqlgen-dataloader/db/loaders"
	"github.com/motemen/example-gqlgen-dataloader/graph"
	"github.com/motemen/example-gqlgen-dataloader/graph/generated"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, err := db.Init()
	if err != nil {
		panic(err)
	}

	var h http.Handler = handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{Resolvers: &graph.Resolver{DB: db}},
		),
	)
	h = loaders.Middleware(db, h)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", h)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
