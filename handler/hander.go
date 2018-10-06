package handler

import (
	"encoding/json"
	"net/http"
)

// GetGasPrice gets latest gas price.
func GetGasPrice(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)
	json.NewEncoder(w).Encode("hello")
}

// CreateTransactionList adds a list of transactions with different gas prices to transactions
// pool.
func CreateTransactionList(w http.ResponseWriter, r *http.Request) {
	//  params := mux.Vars(r)
	json.NewEncoder(w).Encode("ok")
}
