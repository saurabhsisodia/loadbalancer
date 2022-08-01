package routes

import (
	"github.com/gorilla/mux"

	"github.com/saurabhsisodia/loadbalancer/handlers"
)

func Handlers() *mux.Router {

	r := mux.NewRouter()

	r.HandleFunc("/urls/register", handlers.Register).Methods("POST")
	r.HandleFunc("/proxy", handlers.Proxy)

	r.HandleFunc("/urls/get", handlers.Get)
	r.HandleFunc("/urls/delete", handlers.Delete).Methods("DELETE")

	return r

}
