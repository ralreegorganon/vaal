package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ralreegorganon/vaal/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHTTPServer(t *testing.T) {
	Convey("Subject: HttpServer responds to requests appropriately", t, func() {
		fixture := newServerFixture()
		Convey("When a replay is requested", func() {
			id := 1
			expectedReplay := &models.Replay{Id: id}
			status, replay := fixture.GetReplayById(id)

			Convey("The server returns it", func() {
				So(replay.Id, ShouldEqual, expectedReplay.Id)
			})

			Convey("The server returns HTTP 200 - OK", func() {
				So(status, ShouldEqual, http.StatusOK)
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
	self := new(ServerFixture)
	self.administrator = newFakeAdministrator()
	self.server = NewHTTPServer(self.administrator)
	self.router = CreateRoutes(self.server)
	return self
}

func (self *ServerFixture) GetReplayById(id int) (int, *models.Replay) {
	request, _ := http.NewRequest("GET", fmt.Sprintf("/replays/%v", id), nil)
	response := httptest.NewRecorder()

	self.router.ServeHTTP(response, request)

	decoder := json.NewDecoder(strings.NewReader(response.Body.String()))
	replay := new(models.Replay)
	decoder.Decode(replay)

	return response.Code, replay
}

type FakeAdministrator struct {
}

func newFakeAdministrator() *FakeAdministrator {
	return new(FakeAdministrator)
}

func (self *FakeAdministrator) GetReplayById(id int) *models.Replay {
	replay := &models.Replay{Id: id}
	return replay
}
