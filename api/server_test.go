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
	"github.com/ralreegorganon/vaal/models"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	goodReplayId = 1
	badReplayId  = 2
	goodEndpoint = "http://localhost:9999"
	badEndpoint  = "http://localhost:666"
)

func TestHTTPServer(t *testing.T) {
	Convey("Subject: HttpServer responds to requests appropriately", t, func() {
		fixture := newServerFixture()
		Convey("When a replay is requested", func() {
			Convey("And it exists", func() {
				expectedReplay := &models.Replay{Id: goodReplayId}
				status, replay := fixture.GetReplayById(expectedReplay.Id)

				Convey("The server returns it", func() {
					So(replay.Id, ShouldEqual, expectedReplay.Id)
				})

				Convey("The server returns HTTP 200 - OK", func() {
					So(status, ShouldEqual, http.StatusOK)
				})
			})
			Convey("And it doesn't exist", func() {
				expectedReplay := &models.Replay{Id: badReplayId}
				status, replay := fixture.GetReplayById(expectedReplay.Id)

				Convey("The server returns nothing", func() {
					So(replay, ShouldEqual, nil)
				})

				Convey("The server returns HTTP 404 - Not Found", func() {
					So(status, ShouldEqual, http.StatusNotFound)
				})
			})

		})
		Convey("When a client tries to join a match", func() {
			Convey("And provides a valid endpoint", func() {
				status := fixture.JoinMatch(goodEndpoint)

				Convey("The server returns HTTP 200 - OK", func() {
					So(status, ShouldEqual, http.StatusOK)
				})
			})
			Convey("And provides an invalid endpoint", func() {
				status := fixture.JoinMatch(badEndpoint)

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

func (self *ServerFixture) GetReplayById(id int) (int, *models.Replay) {
	request, _ := http.NewRequest("GET", fmt.Sprintf("/replays/%v", id), nil)
	response := httptest.NewRecorder()

	self.router.ServeHTTP(response, request)

	if response.Code == http.StatusOK {
		decoder := json.NewDecoder(strings.NewReader(response.Body.String()))
		replay := &models.Replay{}
		decoder.Decode(replay)
		return response.Code, replay
	}

	return response.Code, nil
}

func (self *ServerFixture) JoinMatch(endpoint string) int {
	message := &joinMatchMessage{
		Endpoint: endpoint,
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

type FakeAdministrator struct {
}

func newFakeAdministrator() *FakeAdministrator {
	return &FakeAdministrator{}
}

func (self *FakeAdministrator) GetReplayById(id int) (*models.Replay, error) {
	switch id {
	case goodReplayId:
		return &models.Replay{Id: id}, nil
	}
	return nil, errors.New("purposeful test failure")
}

func (self *FakeAdministrator) JoinMatch(endpoint string) error {
	switch endpoint {
	case goodEndpoint:
		return nil
	case badEndpoint:
		return errors.New("purposeful test failure")
	}
	return errors.New("all wrong")
}
