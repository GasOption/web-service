package main

import (
	"log"
	"net/http"

	"github.com/GasOption/web-service/handler"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/gasprice", handler.GetGasPrice).Methods("GET")
	router.HandleFunc("/transactionlist/{id}", handler.CreateTransactionList).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
