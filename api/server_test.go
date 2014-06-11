package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ralreegorganon/vaal/replay"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	goodReplayId = 1
	badReplayId  = 2
	goodEndpoint = "http://localhost:9999"
	badEndpoint  = "http://localhost:666"
	goodMatch    = "ABCDEF"
	badMatch     = "XYZ"
)

func TestHTTPServer(t *testing.T) {
	Convey("Subject: HttpServer responds to requests appropriately", t, func() {
		fixture := newServerFixture()
		Convey("When a replay is requested", func() {
			Convey("And it exists", func() {
				expectedReplay := &replay.Replay{Id: goodReplayId}
				status, replay := fixture.GetReplayById(expectedReplay.Id)

				Convey("The server returns it", func() {
					So(replay.Id, ShouldEqual, expectedReplay.Id)
				})

				Convey("The server returns HTTP 200 - OK", func() {
					So(status, ShouldEqual, http.StatusOK)
				})
			})
			Convey("And it doesn't exist", func() {
				expectedReplay := &replay.Replay{Id: badReplayId}
				status, replay := fixture.GetReplayById(expectedReplay.Id)

				Convey("The server returns nothing", func() {
					So(replay, ShouldEqual, nil)
				})

				Convey("The server returns HTTP 404 - Not Found", func() {
					So(status, ShouldEqual, http.StatusNotFound)
				})
			})

		})

		Convey("When a new match is requested", func() {
			status, match := fixture.CreateMatch()

			Convey("The server returns the match id", func() {
				So(match, ShouldEqual, goodMatch)
			})

			Convey("The server returns HTTP 200 - OK", func() {
				So(status, ShouldEqual, http.StatusOK)
			})
		})

		Convey("When a match start is requested", func() {
			status := fixture.StartMatch(goodMatch)

			Convey("The server returns HTTP 200 - OK", func() {
				So(status, ShouldEqual, http.StatusOK)
			})
		})

		Convey("When a client tries to join a match", func() {
			Convey("And provides a valid match id", func() {
				status := fixture.JoinMatch(goodEndpoint, goodMatch)

				Convey("The server returns HTTP 200 - OK", func() {
					So(status, ShouldEqual, http.StatusOK)
				})
			})
			Convey("And provides an invalid match id", func() {
				status := fixture.JoinMatch(goodEndpoint, badMatch)

				Convey("The server returns HTTP 400 - Bad Request", func() {
					So(status, ShouldEqual, http.StatusBadRequest)
				})
			})
			Convey("And provides a valid endpoint", func() {
				status := fixture.JoinMatch(goodEndpoint, goodMatch)

				Convey("The server returns HTTP 200 - OK", func() {
					So(status, ShouldEqual, http.StatusOK)
				})
			})
			Convey("And provides an invalid endpoint", func() {
				status := fixture.JoinMatch(badEndpoint, goodMatch)

				Convey("The server returns HTTP 400 - Bad Request", func() {
					So(status, ShouldEqual, http.StatusBadRequest)
				})
			})
		})
	})
}

type ServerFixture struct {
	server        *HTTPServer
	administrator *FakeAdministrator
	router        *mux.Router
}

func newServerFixture() *ServerFixture {
	self := &ServerFixture{}
	self.administrator = newFakeAdministrator()
	self.server = NewHTTPServer(self.administrator)
	self.router, _ = CreateRouter(self.server)
	return self
}

func (self *ServerFixture) GetReplayById(id int) (int, *replay.Replay) {
	request, _ := http.NewRequest("GET", fmt.Sprintf("/replays/%v", id), nil)
	response := httptest.NewRecorder()

	self.router.ServeHTTP(response, request)

	if response.Code == http.StatusOK {
		decoder := json.NewDecoder(strings.NewReader(response.Body.String()))
		replay := &replay.Replay{}
		decoder.Decode(replay)
		return response.Code, replay
	}

	return response.Code, nil
}

func (self *ServerFixture) CreateMatch() (int, string) {
	request, _ := http.NewRequest("POST", "/create", nil)
	response := httptest.NewRecorder()

	self.router.ServeHTTP(response, request)

	if response.Code == http.StatusOK {
		decoder := json.NewDecoder(strings.NewReader(response.Body.String()))
		match := &createMatchResponseMessage{}
		decoder.Decode(match)
		return response.Code, match.Match
	}

	return response.Code, ""
}

func (self *ServerFixture) JoinMatch(endpoint, match string) int {
	message := &joinMatchRequestMessage{
		Endpoint: endpoint,
		Match:    match,
	}

	json, err := json.Marshal(message)
	if err != nil {
		return http.StatusInternalServerError
	}

	request, _ := http.NewRequest("POST", "/join", bytes.NewBuffer(json))
	response := httptest.NewRecorder()

	self.router.ServeHTTP(response, request)

	return response.Code
}

func (self *ServerFixture) StartMatch(match string) int {
	message := &startMatchRequestMessage{
		Match: match,
	}

	json, err := json.Marshal(message)
	if err != nil {
		return http.StatusInternalServerError
	}

	request, _ := http.NewRequest("POST", "/start", bytes.NewBuffer(json))
	response := httptest.NewRecorder()

	self.router.ServeHTTP(response, request)

	return response.Code
}

type FakeAdministrator struct {
}

func newFakeAdministrator() *FakeAdministrator {
	return &FakeAdministrator{}
}

func (self *FakeAdministrator) GetReplayById(id int) (*replay.Replay, error) {
	switch id {
	case goodReplayId:
		return &replay.Replay{Id: id}, nil
	}
	return nil, errors.New("purposeful test failure")
}

func (self *FakeAdministrator) JoinMatch(endpoint, match string) error {
	switch endpoint {
	case goodEndpoint:
		switch match {
		case goodMatch:
			return nil
		case badMatch:
			return errors.New("bad match")
		}
		return nil
	case badEndpoint:
		return errors.New("bad endpoint")
	}
	return errors.New("all wrong")
}

func (a *FakeAdministrator) CreateMatch() (string, error) {
	return goodMatch, nil
}

func (a *FakeAdministrator) StartMatch(match string) error {
	switch match {
	case goodMatch:
		return nil
	default:
		return errors.New("bad match")
	}
}
