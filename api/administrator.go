package api

import (
	"encoding/json"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ralreegorganon/vaal/models"
)

type Administrator struct {
	db *sqlx.DB
}

type jsonContainer struct {
	Data string
}

func NewAdministrator() *Administrator {
	connectionString := os.Getenv("VAAL_CONNECTION_STRING")

	if connectionString == "" {
		connectionString = "user=vaal password=ourrobotoverlords dbname=vaal sslmode=disable"
	}

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatalln(err)
	}

	return &Administrator{db: db}
}

func (self *Administrator) GetReplayById(id int) (*models.Replay, error) {
	jsonContainer := &jsonContainer{}

	err := self.db.Get(jsonContainer, "select data from replays where replay_id = $1", id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	replay := &models.Replay{}
	json.Unmarshal([]byte(jsonContainer.Data), replay)

	return replay, nil
}

func (self *Administrator) JoinMatch(endpoint string) error {
	return nil
}
