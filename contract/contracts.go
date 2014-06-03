package contract

import (
	"net/http"

	"github.com/ralreegorganon/vaal/models"
)

type Server interface {
	GetReplay(writer http.ResponseWriter, request *http.Request)
}

type Administrator interface {
	GetReplayById(id int) *models.Replay
}
