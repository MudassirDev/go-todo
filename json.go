package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, statusCode int, message string, err error) {
	log.Printf("error while processing request: %v", err)

	type Response struct {
		Message string `json:"msg"`
	}

	respondWithJSON(w, statusCode, Response{
		Message: message,
	})
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to stringify request"))
		return
	}

	w.WriteHeader(statusCode)
	w.Write(data)
}
