package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	URL string `json:"url"`
}

var Conf Config

func init() {
	configFile, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
	}
	json.Unmarshal(configFile, &Conf)
}
