package main

import (
	"log"
	"os"

	"github.com/stageddat/shelter-node/internal/db"
	"github.com/stageddat/shelter-node/internal/db/postgres"
	"github.com/stageddat/shelter-node/internal/db/sqlite"
	"github.com/stageddat/shelter-node/internal/server"
)

func main() {
	var store db.Store
	var err error

	// database can be sqlite or postgres
	switch os.Getenv("DB_DRIVER") {
	case "postgres":
		store, err = postgres.New(os.Getenv("DATABASE_URL"))
	case "sqlite":
		store, err = sqlite.New(os.Getenv("DATABASE_PATH"))
	default:
		store, err = sqlite.New(os.Getenv("DATABASE_PATH"))
	}

	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer store.Close()

	srv := server.New(store)
	log.Fatal(srv.Start(":4123"))
}
