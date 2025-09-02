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

func (apiCfg *APIConfig) LoginUser(w http.ResponseWriter, r *http.Request) {
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

	user, err := apiCfg.DB.GetUserWithUsername(context.Background(), req.Username)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "no user found!", err)
		return
	}

	err = auth.VerifyPassword(req.Password, user.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "wrong password!", err)
		return
	}

	userId, ok := user.ID.(string)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "failed to create access token!", fmt.Errorf("user id not string!"))
		return
	}

	token, err := auth.CreateJWT(JWT_SECRET, JWT_EXPIRES_IN, userId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create access token!", err)
		return
	}

	cookie := http.Cookie{
		Name:     AUTH_KEY,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(JWT_EXPIRES_IN),
		MaxAge:   int(JWT_EXPIRES_IN),
		Secure:   false,
		HttpOnly: false,
	}

	http.SetCookie(w, &cookie)
	respondWithJSON(w, http.StatusOK, struct {
		ID        any       `json:"id"`
		Name      string    `json:"name"`
		Username  string    `json:"username"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}{
		ID:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}
