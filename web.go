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

	mux.Handle("POST /api/tasks/create", apiCfg.AuthMiddleware(apiCfg.CreateTask()))
	mux.Handle("POST /api/tasks/delete", apiCfg.AuthMiddleware(apiCfg.DeleteTask()))

	mux.HandleFunc("POST /api/users/create", apiCfg.CreateUser)
	mux.HandleFunc("POST /api/users/login", apiCfg.LoginUser)

	return mux
}
