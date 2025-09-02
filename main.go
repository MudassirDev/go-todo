package main

import (
	"database/sql"
	"embed"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

var (
	DB_CONN *sql.DB
	//go:embed db/schema/*.sql
	embedMigrations embed.FS
)

const (
	DB_PATH string = "app.db"
)

// initialize everything
func init() {
	godotenv.Load()

	log.Println("loading env variables")

	port := os.Getenv("PORT")
	validateEnv(port, "PORT")

	log.Println("env variables loaded")

	log.Println("making a connection with db")

	conn, err := sql.Open("sqlite3", DB_PATH)
	if err != nil {
		log.Fatalf("failed to form a connection: %v", err)
	}
	DB_CONN = conn

	log.Println("connection DB formed!")
}

// running migrations
func init() {
	log.Println("running migrations")

	goose.SetDialect("sqlite3")
	goose.SetBaseFS(embedMigrations)

	if err := goose.Up(DB_CONN, "db/schema"); err != nil { //
		log.Fatalf("failed to run migrations: %v", err)
	}

	log.Println("migrations ran successfully!")
}

func main() {
	defer DB_CONN.Close()
}

func validateEnv(env, envName string) {
	if env == "" {
		log.Fatalf("cannot use empty variable: %v", envName)
	}
}
