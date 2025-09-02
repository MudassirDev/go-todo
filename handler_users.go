package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MudassirDev/todo-go/db/database"
	"github.com/MudassirDev/todo-go/internal/auth"
	"github.com/google/uuid"
)

const (
	JSON_TYPE string = "application/json"
)

type UserRequest struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (apiCfg *APIConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	err := validateContentType(r.Header, JSON_TYPE)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid content type", err)
		return
	}

	var req UserRequest
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to decode request", err)
		return
	}

	if req.Username == "" {
		respondWithError(w, http.StatusBadRequest, "username cannot be empty", fmt.Errorf("empty username"))
		return
	}

	if req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "password cannot be empty", fmt.Errorf("empty password"))
		return
	}

	password, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to hash password", err)
		return
	}

	user, err := apiCfg.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      req.Name,
		Username:  req.Username,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		if !strings.Contains(err.Error(), "UNIQUE") {
			respondWithError(w, http.StatusInternalServerError, "failed to create user", err)
			return
		}
		if !strings.Contains(err.Error(), "users.username") {
			respondWithError(w, http.StatusBadRequest, "duplicate key", err)
			return
		}
		respondWithError(w, http.StatusBadRequest, "duplicate key: username", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}
