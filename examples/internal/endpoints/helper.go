package endpoints

import (
	"encoding/json"
	"log"
	"net/http"
)

func handler500(w http.ResponseWriter) {
	respondWithError(w,
		http.StatusInternalServerError,
		http.StatusText(http.StatusInternalServerError))
}

func handler404(w http.ResponseWriter) {
	respondWithError(w,
		http.StatusNotFound,
		http.StatusText(http.StatusNotFound))
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Println("failed to marshall", err)
		handler500(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		log.Println("failed to write", err)
	}
}
