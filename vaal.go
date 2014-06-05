package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ralreegorganon/vaal/api"
)

func main() {
	conf := loadConfigFromFile("./config.json")

	administrator := api.NewAdministrator()
	server := api.NewHTTPServer(administrator)
	router, err := api.CreateRouter(server)
	if err != nil {
		fmt.Println(err)
		return
	}
	http.Handle("/", router)

	log.Printf("Vaal server started at http://%s\n", conf.URL)
	err = http.ListenAndServe(conf.URL, nil)
	if err != nil {
		fmt.Println(err)
	}
}

type config struct {
	URL string `json:"url"`
}

func loadConfigFromFile(path string) *config {
	conf := &config{}
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
	}
	json.Unmarshal(configFile, conf)
	return conf
}
