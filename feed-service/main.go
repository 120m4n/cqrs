package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/120m4n/cqrs/database"
	"github.com/120m4n/cqrs/events"
	"github.com/120m4n/cqrs/repository"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PostgresDB       string `envconfig:"POSTGRES_DB" default:"postgres"`
	PostgresUser     string `envconfig:"POSTGRES_USER" default:"postgres"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD" default:"postgres"`
	NatsAddress      string `envconfig:"NATS_ADDRESS" default:"localhost:4222"`
}

func NewRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/feeds", createdFeedHandler).Methods("POST")
	return
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)

	if err != nil {
		log.Fatal(err)
	}

	addr := fmt.Sprintf("poastgres://%s:%s@postgres/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)
	fmt.Println(addr)

	repo, err := database.NewPostgresRepository(addr)
	if err != nil {
		log.Fatal(err)
	}

	defer repo.Close()

	repository.SetRepository(repo)

	n, err := events.NewNatsEventStore(cfg.NatsAddress)
	if err != nil {
		log.Fatal(err)
	}

	events.SetEventStore(n)
	defer events.Close()

	router := NewRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}

}
