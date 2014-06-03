package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ralreegorganon/vaal/models"
)

func GetReplay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 0, 0)
	if err != nil {
		log.Println("Trash")
	}
	thing := models.Replay{Id: id}

	writeJSON(w, thing)
}
