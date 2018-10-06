package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/GasOption/web-service/handler"
	"github.com/GasOption/web-service/processor"
	"github.com/GasOption/web-service/store"
	"github.com/GasOption/web-service/txn"
	"github.com/gorilla/mux"
)

var (
	port    = flag.String("port", ":8080", "Port.")
	ethAddr = flag.String("eth_addr", "https://rinkeby.infura.io", "ETH blockchain address.")
)

func main() {
	ctx := context.Background()

	// Transaction client.
	txnClient, err := txn.NewClient(*ethAddr)
	if err != nil {
		log.Fatalf("txn.NewClient() = %v", err)
	}

	// Store client.
	storeClient := store.New()

	// Processor.
	go processor.Process(ctx, storeClient, txnClient)

	// Router and handler.
	router := mux.NewRouter()
	handlerClient := handler.New(storeClient)
	router.HandleFunc("/gasprice", handlerClient.GetGasPrice).Methods("GET")
	router.HandleFunc("/transactionlist", handlerClient.UpdateTxnList).Methods("POST")
	router.HandleFunc("/pool", handlerClient.GetPool).Methods("GET")

	// HTTP server.
	log.Printf("Running web service at port %v", *port)
	log.Fatal(http.ListenAndServe(*port, router))
}
