package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ralreegorganon/vaal/administrator"
	"github.com/ralreegorganon/vaal/api"
	"github.com/ralreegorganon/vaal/config"
)

func main() {
	administrator := administrator.NewAdministrator()
	server := api.NewHTTPServer(administrator)
	router := api.CreateRoutes(server)
	http.Handle("/", router)

	log.Printf("Vaal server started at http://%s\n", config.Conf.URL)
	err := http.ListenAndServe(config.Conf.URL, nil)
	if err != nil {
		fmt.Println(err)
	}
}
