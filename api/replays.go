package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func GetReplay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	writeJSON(w, id)
}
