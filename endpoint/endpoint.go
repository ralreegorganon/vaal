package endpoint

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type Endpoint struct {
	Root string
}

func (e *Endpoint) Validate() error {
	err := e.Status()
	return err
}

func (e *Endpoint) Status() error {
	u, err := url.Parse(e.Root)
	if err != nil {
		return err
	}

	u.Path = "status"

	log.Println("Making request to ", u.String())
	r, err := http.Get(u.String())
	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("expected HTTP 200 - OK, got %v", r.StatusCode))
	}

	return nil
}
