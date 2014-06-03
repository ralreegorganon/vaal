package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func CreateRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/replays/{id:[0-9]+}", GetReplay).Methods("GET")
	http.Handle("/", r)
}
