package main

import (
	"log"
	"net/http"
	"os"

	"github.com/saurabhsisodia/loadbalancer/routes"
)

func main() {

	http.Handle("/", routes.Handlers())

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
