package main

import (
	"context"
	"net/http"

	"github.com/MudassirDev/todo-go/internal/auth"
)

func (apiCfg *APIConfig) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(AUTH_KEY)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "user not logged in!", err)
			return
		}

		id, err := auth.VerifyJWT(JWT_SECRET, cookie.Value)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "user not logged in!", err)
			return
		}

		user, err := apiCfg.DB.GetUserWithUserID(context.Background(), id)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "user not logged in!", err)
			return
		}

		ctx := context.WithValue(r.Context(), AUTH_KEY, user)
		request := r.WithContext(ctx)
		next.ServeHTTP(w, request)
	})
}
