package api

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ralreegorganon/vaal/endpoint"
	"github.com/ralreegorganon/vaal/engine"
	"github.com/ralreegorganon/vaal/replay"
	"github.com/satori/go.uuid"
)

type Administrator struct {
	db                *sqlx.DB
	activeMatchesLock sync.RWMutex
	activeMatches     map[string]*engine.Match
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

	m := make(map[string]*engine.Match)

	return &Administrator{db: db, activeMatches: m}
}

func (a *Administrator) GetReplayById(id int) (*replay.Replay, error) {
	jsonContainer := &jsonContainer{}

	err := a.db.Get(jsonContainer, "select data from replays where replay_id = $1", id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	replay := &replay.Replay{}
	json.Unmarshal([]byte(jsonContainer.Data), replay)

	return replay, nil
}

func (a *Administrator) JoinMatch(root, match string) error {
	e := &endpoint.Endpoint{
		Root: root,
	}

	err := e.Validate()
	if err != nil {
		return err
	}

	a.activeMatchesLock.Lock()
	defer a.activeMatchesLock.Unlock()

	if _, ok := a.activeMatches[match]; !ok {
		log.Println("match doesn't exist in memory, getting it")
		m := &engine.Match{}
		err := a.db.Get(m, "select match, open, complete from matches where match = $1", match)
		if err != nil {
			return err
		}

		if !m.Open {
			return errors.New("attempted to join closed match")
		}

		m.Endpoints = make([]*endpoint.Endpoint, 0)
		a.activeMatches[match] = m
	}
	m := a.activeMatches[match]

	m.Endpoints = append(m.Endpoints, e)

	return nil
}

func (a *Administrator) CreateMatch() (string, error) {
	u := uuid.NewV4().String()

	_, err := a.db.Exec("insert into matches (match) values ($1)", u)

	if err != nil {
		return "", err
	}

	m := &engine.Match{}

	err = a.db.Get(m, "select match from matches where match = $1", u)
	if err != nil {
		return "", err
	}

	return m.Match, nil
}

func (a *Administrator) StartMatch(match string) error {
	a.activeMatchesLock.Lock()
	m, ok := a.activeMatches[match]
	a.activeMatchesLock.Unlock()

	if ok {
		go func() {
			m.Open = false
			a.db.Exec("update matches set open = false where match = $1", m.Match)
			m.Start()

			text, err := json.Marshal(m.Replay)
			if err != nil {
				log.Println(err)
			}

			m.Complete = true
			result, err := a.db.NamedQuery("insert into replays (data) values (:text) returning replay_id", map[string]interface{}{"text": text})
			if err != nil {
				log.Println(err)
			}
			var id int
			if result.Next() {
				result.Scan(&id)
			}
			a.db.Exec("update matches set complete = true, replay_id = $2 where match = $1", m.Match, id)

			a.activeMatchesLock.Lock()
			delete(a.activeMatches, m.Match)
			a.activeMatchesLock.Unlock()
		}()
	}

	return nil
}
