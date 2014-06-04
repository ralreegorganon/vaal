package contracts

import (
	"net/http"

	"github.com/ralreegorganon/vaal/models"
)

type Server interface {
	GetReplay(writer http.ResponseWriter, request *http.Request, vars map[string]string) error
}

type Administrator interface {
	GetReplayById(id int) *models.Replay
}
