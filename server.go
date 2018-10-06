package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/GasOption/web-service/handler"
	"github.com/GasOption/web-service/txn"
	"github.com/gorilla/mux"
)

var (
	port    = flag.String("port", ":8080", "Port.")
	ethAddr = flag.String("eth_addr", "https://rinkeby.infura.io", "ETH blockchain address.")
)

func main() {
	// Transaction client.
	_, err := txn.NewClient(*ethAddr)
	if err != nil {
		log.Fatalf("txn.NewClient() = %v", err)
	}

	// REST server.
	router := mux.NewRouter()
	router.HandleFunc("/gasprice", handler.GetGasPrice).Methods("GET")
	router.HandleFunc("/transactionlist/{id}", handler.CreateTransactionList).Methods("POST")
	log.Printf("Running web service at port %v", *port)
	log.Fatal(http.ListenAndServe(*port, router))
}
