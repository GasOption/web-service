package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/GasOption/web-service/handler"
	"github.com/GasOption/web-service/store"
	"github.com/GasOption/web-service/txn"
	"github.com/gorilla/mux"
)

var (
	port    = flag.String("port", ":8080", "Port.")
	ethAddr = flag.String("eth_addr", "https://rinkeby.infura.io", "ETH blockchain address.")
)

func main() {
	// ctx := context.Background()

	// Transaction client.
	_, err := txn.NewClient(*ethAddr)
	if err != nil {
		log.Fatalf("txn.NewClient() = %v", err)
	}

	// Test.
	/*
		hexTxn := "f86d8202b28477359400825208944592d8f8d7b001e72cb26a73e4fa1806a51ac79d880de0b6b3a7640000802ca05924bde7ef10aa88db9c66dd4f5fb16b46dff2319b9968be983118b57bb50562a001b24b31010004f13d9a26b320845257a6cfc2bf819a3d55e3fc86263c5f0772"
		if err := txnClient.SendTransaction(ctx, hexTxn); err != nil {
			log.Fatalf("txnClient.SendTransaction() = %v", err)
		}
	*/

	// Store client.
	storeClient := store.New()

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
