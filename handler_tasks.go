package main

import "net/http"

func (apiCfg *APIConfig) CreateTask() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}
