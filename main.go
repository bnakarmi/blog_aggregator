package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bnakarmi/blog_aggregator/repository"
	"github.com/bnakarmi/blog_aggregator/service"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_CONN")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("Cannot ping database because %s", err)
	}

	repository := repository.NewQueryRepository(db)
	routingService := service.NewRoutingService(repository)
    feedWorkerService := service.NewFeedWorkerService(repository)

    router := routingService.Initialize()
    feedWorkerService.Initialize()

	fmt.Printf("Server started at localhost:%s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}
