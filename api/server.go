package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ralreegorganon/vaal/contracts"
)

func CreateRouter(server contracts.Server) (*mux.Router, error) {
	r := mux.NewRouter()
	m := map[string]map[string]HttpApiFunc{
		"GET": {
			"/replays/{id:[0-9]+}": server.GetReplay,
		},
		"POST": {
			"/join": server.JoinMatch,
		},
	}

	for method, routes := range m {
		for route, handler := range routes {
			localRoute := route
			localHandler := handler
			localMethod := method
			f := makeHttpHandler(localMethod, localRoute, localHandler)

			r.Path(localRoute).Methods(localMethod).HandlerFunc(f)
		}
	}

	return r, nil
}

func makeHttpHandler(localMethod string, localRoute string, handlerFunc HttpApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeCorsHeaders(w, r)
		if err := handlerFunc(w, r, mux.Vars(r)); err != nil {
			httpError(w, err)
		}
	}
}

func writeCorsHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
}

type HttpApiFunc func(w http.ResponseWriter, r *http.Request, vars map[string]string) error

type HTTPServer struct {
	administrator contracts.Administrator
}

func NewHTTPServer(administrator contracts.Administrator) *HTTPServer {
	self := &HTTPServer{
		administrator: administrator,
	}

	return self
}

func writeJSON(w http.ResponseWriter, code int, thing interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	val, err := json.Marshal(thing)
	w.Write(val)
	return err
}

func httpError(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError

	if err != nil {
		http.Error(w, err.Error(), statusCode)
	}
}

func (self *HTTPServer) GetReplay(w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}

	thing, err := self.administrator.GetReplayById(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	writeJSON(w, http.StatusOK, thing)

	return nil
}

func (self *HTTPServer) JoinMatch(w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	decoder := json.NewDecoder(r.Body)
	message := &joinMatchMessage{}

	err := decoder.Decode(message)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	err = self.administrator.JoinMatch(message.Endpoint)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
