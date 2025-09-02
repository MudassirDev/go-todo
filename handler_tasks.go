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
			respondWithError(w, http.StatusBadRequest, "task cannot be empty!", fmt.Errorf("no task found"))
			return
		}

		task, err := apiCfg.DB.CreateTask(context.Background(), database.CreateTaskParams{
			ID:        uuid.New(),
			UserID:    user.ID,
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

func (apiCfg *APIConfig) DeleteTask() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type Request struct {
			TaskID string `json:"task_id"`
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

		if req.TaskID == "" {
			respondWithError(w, http.StatusBadRequest, "task id cannot be empty!", fmt.Errorf("no task id"))
			return
		}

		err = apiCfg.DB.DeleteTaskWithID(context.Background(), database.DeleteTaskWithIDParams{
			ID:     req.TaskID,
			UserID: user.ID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "failed to delete task!", err)
			return
		}

		respondWithJSON(w, http.StatusOK, struct {
			Message string `json:"msg"`
		}{
			Message: "Task deleted successfully!",
		})
	})
}

func (apiCfg *APIConfig) UpdateTask() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type Request struct {
			TaskID      string `json:"task_id"`
			IsCompleted bool   `json:"is_completed"`
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

		if req.TaskID == "" {
			respondWithError(w, http.StatusBadRequest, "task id cannot be empty!", fmt.Errorf("no task id"))
			return
		}

		task, err := apiCfg.DB.UpdateTaskWithID(context.Background(), database.UpdateTaskWithIDParams{
			IsCompleted: req.IsCompleted,
			ID:          req.TaskID,
			UserID:      user.ID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "failed to update task", err)
			return
		}

		respondWithJSON(w, http.StatusOK, task)
	})
}
