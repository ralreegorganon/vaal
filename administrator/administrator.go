package administrator

import "github.com/ralreegorganon/vaal/models"

type Administrator struct {
}

func (self *Administrator) GetReplayById(id int) *models.Replay {
	replay := &models.Replay{Id: id}
	return replay
}

func NewAdministrator() *Administrator {
	return &Administrator{}
}
