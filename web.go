package main

import (
	"net/http"
	"text/template"

	"github.com/MudassirDev/todo-go/db/database"
)

var (
	Templates *template.Template = template.New("")
)

type APIConfig struct {
	DB *database.Queries
}

func CreateMux() *http.ServeMux {
	mux := http.NewServeMux()

	apiCfg := APIConfig{}

	queries := database.New(DB_CONN)
	apiCfg.DB = queries

	parseTemplates()

	fs := http.FileServer(http.Dir("static/assets"))

	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Templates.ExecuteTemplate(w, "index.html", nil)
	})
	mux.Handle("/tasks", apiCfg.AuthMiddleware(apiCfg.handlerTasks()))

	mux.Handle("POST /api/tasks/create", apiCfg.AuthMiddleware(apiCfg.CreateTask()))
	mux.Handle("POST /api/tasks/delete", apiCfg.AuthMiddleware(apiCfg.DeleteTask()))
	mux.Handle("POST /api/tasks/update", apiCfg.AuthMiddleware(apiCfg.UpdateTask()))

	mux.HandleFunc("POST /api/users/create", apiCfg.CreateUser)
	mux.HandleFunc("POST /api/users/login", apiCfg.LoginUser)

	return mux
}
