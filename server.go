package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/GasOption/web-service/handler"
	"github.com/gorilla/mux"
)

var (
	port = flag.String("port", ":8080", "Port.")
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/gasprice", handler.GetGasPrice).Methods("GET")
	router.HandleFunc("/transactionlist/{id}", handler.CreateTransactionList).Methods("POST")

	log.Printf("Running web service at port %v", *port)
	log.Fatal(http.ListenAndServe(*port, router))
}
