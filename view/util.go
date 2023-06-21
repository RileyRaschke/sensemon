package view

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func AsJson(w http.ResponseWriter, s interface{}) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")
	if err := enc.Encode(s); err != nil {
		log.Errorf("Encoding error: %w", err)
		JsonErrorMsg(w, http.StatusInternalServerError, "Encoding error")
		return
	}
}

func JsonErrorMsg(w http.ResponseWriter, stat int, msg string) {
	w.WriteHeader(stat)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = msg
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}
