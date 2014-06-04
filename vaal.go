package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ralreegorganon/vaal/api"
)

func main() {
	administrator := administrator.NewAdministrator()
	server := api.NewHTTPServer(administrator)
	router := api.CreateRouter(server)
	http.Handle("/", router)

	log.Printf("Vaal server started at http://%s\n", config.conf.URL)
	err := http.ListenAndServe(config.conf.URL, nil)
	if err != nil {
		fmt.Println(err)
	}
}
