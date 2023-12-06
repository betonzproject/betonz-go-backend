package jsonutils

import (
	"encoding/json"
	"log"
	"net/http"
)

func Write(w http.ResponseWriter, v any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Panicln("Can't encode json: " + err.Error())
	}
}
