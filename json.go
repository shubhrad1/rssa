package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5XX error:", msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}
	respondJSON(w, code, errResponse{
		Error: msg,
	})
}

func respondJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)

	if err != nil {
		log.Printf("[ERROR] Failed to marshal JSON Response: %v", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
