package main

import (
	"log"
	"net/http"

	"github.com/saurabhsisodia/loadbalancer/routes"
)

func main() {

	http.Handle("/", routes.Handlers())

	log.Fatal(http.ListenAndServe(":8080", nil))
}
