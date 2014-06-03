package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ralreegorganon/vaal/contract"
)

type HTTPServer struct {
	administrator contract.Administrator
}

func NewHTTPServer(administrator contract.Administrator) *HTTPServer {
	self := new(HTTPServer)
	self.administrator = administrator
	return self
}

func CreateRoutes(server contract.Server) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/replays/{id:[0-9]+}", addDefaultHeaders(server.GetReplay)).Methods("GET")
	return r
}

func addDefaultHeaders(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		fn(w, r)
	}
}

func (self *HTTPServer) GetReplay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
	}

	thing := self.administrator.GetReplayById(id)

	writeJSON(w, thing)
}

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
