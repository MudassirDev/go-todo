package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MudassirDev/todo-go/db/database"
	"github.com/google/uuid"
)

func (apiCfg *APIConfig) CreateTask() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type Request struct {
			Task string `json:"task"`
		}

		rawUser := r.Context().Value(AUTH_KEY)
		user, ok := rawUser.(database.GetUserWithUserIDRow)
		if !ok {
			respondWithError(w, http.StatusUnauthorized, "not logged in!", fmt.Errorf("user type not matched!"))
			return
		}

		err := validateContentType(r.Header, JSON_TYPE)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid content type", err)
			return
		}

		var req Request
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		if err := decoder.Decode(&req); err != nil {
			respondWithError(w, http.StatusInternalServerError, "failed to decode request!", err)
			return
		}

		if req.Task == "" {
			respondWithError(w, http.StatusBadRequest, "task cannot be empty!", err)
			return
		}

		task, err := apiCfg.DB.CreateTask(context.Background(), database.CreateTaskParams{
			ID:        uuid.New(),
			Userid:    user.ID,
			Task:      req.Task,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "failed to create task!", err)
			return
		}

		respondWithJSON(w, http.StatusCreated, task)
	})
}
