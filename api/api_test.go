package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetReplayReturnsWithStatusOK(t *testing.T) {
	request, _ := http.NewRequest("GET", "/replays/1", nil)
	response := httptest.NewRecorder()

	GetReplay(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("Non-expected status code%v:\n\tbody: %v", "200", response.Code)
	}
}
