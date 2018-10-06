// Package hander provides HTTP handling methods for the REST web service.
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GasOption/web-service/store"
)

// GetGasPrice gets latest gas price.
func GetGasPrice(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)
	json.NewEncoder(w).Encode("hello")
}

// UpdateTxnList adds a list of transactions with different gas prices to transactions
// pool.
func UpdateTxnList(w http.ResponseWriter, r *http.Request) {
	var txnList store.TxnList
	_ = json.NewDecoder(r.Body).Decode(&txnList)
	fmt.Print(txnList)
	json.NewEncoder(w).Encode("ok")
}
