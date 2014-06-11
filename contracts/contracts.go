package contracts

import (
	"net/http"

	"github.com/ralreegorganon/vaal/replay"
)

type Server interface {
	GetReplay(writer http.ResponseWriter, request *http.Request, vars map[string]string) error
	JoinMatch(writer http.ResponseWriter, request *http.Request, vars map[string]string) error
	CreateMatch(writer http.ResponseWriter, request *http.Request, vars map[string]string) error
	StartMatch(writer http.ResponseWriter, request *http.Request, vars map[string]string) error
}

type Administrator interface {
	GetReplayById(id int) (*replay.Replay, error)
	JoinMatch(endpoint, match string) error
	CreateMatch() (string, error)
	StartMatch(match string) error
}
