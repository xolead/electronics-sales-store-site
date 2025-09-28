package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"electronic/pkg/handlers"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/product", handlers.CreateProduct).Methods("POST")
	router.HandleFunc("/product", handlers.ReadAllProduct).Methods("GET")
	router.HandleFunc("/product/{id}", handlers.ReadProduct).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
