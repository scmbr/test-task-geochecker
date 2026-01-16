package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/scmbr/test-task-geochecker/pkg/hasher"
	"github.com/scmbr/test-task-geochecker/pkg/logger"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		user := os.Getenv("POSTGRES_USER")
		if user == "" {
			user = "scmbr"
		}
		pass := os.Getenv("POSTGRES_PASSWORD")
		if pass == "" {
			pass = "geochecker"
		}
		dbName := os.Getenv("POSTGRES_DB")
		if dbName == "" {
			dbName = "geochecker"
		}
		dsn = fmt.Sprintf(
			"postgres://%s:%s@db:5432/%s?sslmode=disable",
			user, pass, dbName,
		)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM operators`).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count > 0 {
		fmt.Println("operators already exist, skipping seed")
		return
	}

	operators := []struct {
		Name   string
		APIKey string
	}{
		{"Vasya", "secret1"},
		{"Sasha", "secret2"},
		{"Tanya", "secret3"},
	}
	secret := os.Getenv("API_SECRET_KEY")
	if secret == "" {
		secret = "geochecker"
	}
	logger.Info(secret, nil)
	for _, op := range operators {
		id := uuid.New()
		hash := hasher.HashAPIKey(
			secret,
			op.APIKey,
		)

		_, err := db.Exec(`
			INSERT INTO operators (operator_id, api_key_hash, name)
			VALUES ($1, $2, $3)
		`, id, hash, op.Name)

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Seeding complete")
}
