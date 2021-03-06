package endpoint

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/ralreegorganon/vaal/common"
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
	} else {
		log.Println("OK")
	}

	return nil
}

func (e *Endpoint) Start(a *common.Arena) error {
	u, err := url.Parse(e.Root)
	if err != nil {
		return err
	}

	u.Path = "start"

	log.Println("Making request to ", u.String())
	js, _ := json.Marshal(a)
	r, err := http.Post(u.String(), "application/json", bytes.NewBuffer(js))
	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("expected HTTP 200 - OK, got %v", r.StatusCode))
	} else {
		log.Println("OK")
	}

	return nil
}

func (e *Endpoint) End(won bool) error {
	u, err := url.Parse(e.Root)
	if err != nil {
		return err
	}

	u.Path = "end"

	log.Println("Making request to ", u.String())
	m := &endMatchRequestMessage{
		Won: won,
	}
	js, _ := json.Marshal(m)
	r, err := http.Post(u.String(), "application/json", bytes.NewBuffer(js))
	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("expected HTTP 200 - OK, got %v", r.StatusCode))
	} else {
		log.Println("OK")
	}

	return nil
}

func (e *Endpoint) Think(state *common.RobotState) (*common.RobotCommands, error) {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	u, err := url.Parse(e.Root)
	if err != nil {
		return nil, err
	}
	u.Path = "think"

	log.Println("Making request to ", u.String())
	js, _ := json.Marshal(state)
	r, err := client.Post(u.String(), "application/json", bytes.NewBuffer(js))
	if err != nil {
		return nil, err
	}
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("expected HTTP 200 - OK, got %v", r.StatusCode))
	} else {
		log.Println("OK")
	}

	decoder := json.NewDecoder(r.Body)
	message := &common.RobotCommands{}
	err = decoder.Decode(message)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return message, nil
}
