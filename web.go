package main

import (
	"net/http"

	"github.com/MudassirDev/todo-go/db/database"
)

type APIConfig struct {
	DB *database.Queries
}

func CreateMux() *http.ServeMux {
	mux := http.NewServeMux()

	apiCfg := APIConfig{}

	queries := database.New(DB_CONN)
	apiCfg.DB = queries

	mux.HandleFunc("POST /api/users/create", apiCfg.CreateUser)

	return mux
}
