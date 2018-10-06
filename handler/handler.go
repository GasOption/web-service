// Package hander provides HTTP handling methods for the REST web service.
package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/GasOption/web-service/store"
	"github.com/GasOption/web-service/txn"
	"github.com/ethereum/go-ethereum/core/types"
)

type Handler struct {
	store *store.Store
}

func New(store *store.Store) *Handler {
	return &Handler{store: store}
}

// GetGasPrice gets latest gas price.
func (h *Handler) GetGasPrice(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)
	json.NewEncoder(w).Encode("hello")
}

// UpdateTxnList adds a list of transactions with different gas prices to transactions
// pool.
func (h *Handler) UpdateTxnList(w http.ResponseWriter, r *http.Request) {
	// Parse POST body.
	var hexTxnList txn.HexTxnList
	_ = json.NewDecoder(r.Body).Decode(&hexTxnList)
	log.Print(hexTxnList)

	// Decode hex transactions into raw transactions.
	var rawTxns []*types.Transaction
	for _, hexTxn := range hexTxnList.List {
		log.Printf("Decoding hex transaction: %v", hexTxn)
		rawTxn, err := txn.Decode(hexTxn)
		if err != nil {
			json.NewEncoder(w).Encode(fmt.Errorf("txn.Decode() = %v", err))
			return
		}
		rawTxns = append(rawTxns, rawTxn)
	}

	// Add raw transactions to storage pool,
	if err := h.store.AddToPool(rawTxns); err != nil {
		json.NewEncoder(w).Encode(fmt.Errorf("store.AddToPool() = %v", err))
	}

	json.NewEncoder(w).Encode("200")
}

func (h *Handler) GetPool(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(h.store.Pool())
}
