package main

import (
	"fmt"
	"net/http"

	"github.com/ralreegorganon/vaal/api"
	"github.com/ralreegorganon/vaal/config"
)

func main() {
	api.CreateRoutes()
	fmt.Printf("Vaal server started at http://%s\n", config.Conf.URL)
	http.ListenAndServe(config.Conf.URL, nil)
}
