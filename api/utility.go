package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func writeJSON(w http.ResponseWriter, c interface{}) {
	cj, err := json.MarshalIndent(c, "", "  ")
	if checkError(err, w, "Error creating JSON response", http.StatusInternalServerError) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", cj)
}

func checkError(e error, w http.ResponseWriter, m string, c int) bool {
	if e != nil {
		w.WriteHeader(c)
		return true
	}
	return false
}
