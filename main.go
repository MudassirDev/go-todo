package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
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
}

func main() {}

func validateEnv(env, envName string) {
	if env == "" {
		log.Fatalf("cannot use empty variable: %v", envName)
	}
}
