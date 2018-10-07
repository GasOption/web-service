// Package hander provides HTTP handling methods for the REST web service.
package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/GasOption/web-service/store"
	"github.com/GasOption/web-service/txn"
	"github.com/ethereum/go-ethereum/core/types"
)

type Handler struct {
	store   *store.Store
	ethAddr string
}

func New(store *store.Store, ethAddr string) *Handler {
	return &Handler{store: store, ethAddr: ethAddr}
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
	if err := json.NewDecoder(r.Body).Decode(&hexTxnList); err != nil {
		json.NewEncoder(w).Encode(fmt.Errorf("json.NewDecoder().Decode() = %v", err))
	}
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

func (h *Handler) CreateJsonRpc(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	tee := io.TeeReader(r.Body, &buf)

	type JsonRpcRequest struct {
		JsonRpc string        `json:"jsonrpc"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
		Id      int           `json:"id"`
	}

	var jsonRpcRequest JsonRpcRequest
	if err := json.NewDecoder(tee).Decode(&jsonRpcRequest); err != nil {
		log.Printf("json.NewDecoder().Decode() = %v", err)
		//json.NewEncoder(w).Encode(fmt.Errorf("json.NewDecoder().Decode() = %v", err))
	}

	if jsonRpcRequest.Method == "eth_sendRawTransaction" {
		log.Printf("params for eth_sendRawTransaction: %v", jsonRpcRequest.Params)

		// Decode hex transactions into raw transactions.
		var rawTxns []*types.Transaction
		for _, hexTxn := range jsonRpcRequest.Params {
			log.Printf("Decoding hex transaction: %v", hexTxn)
			rawTxn, err := txn.Decode(hexTxn.(string))
			if err != nil {
				log.Printf("txn.Decode() = %v", err)
				//json.NewEncoder(w).Encode(fmt.Errorf("txn.Decode() = %v", err))
				break
			}
			rawTxns = append(rawTxns, rawTxn)
		}

		// Add raw transactions to storage pool,
		if err := h.store.AddToPool(rawTxns); err != nil {
			log.Printf("store.AddToPool() = %v", err)
			//json.NewEncoder(w).Encode(fmt.Errorf("store.AddToPool() = %v", err))
		}

		type SendRawMessageResponse struct {
			Id      int    `json:"id"`
			JsonRpc string `json:"jsonrpc"`
			Result  string `json:"result"`
		}
		sendRawMessageResponse := SendRawMessageResponse{
			Id:      jsonRpcRequest.Id,
			JsonRpc: jsonRpcRequest.JsonRpc,
			Result:  "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331",
		}
		json.NewEncoder(w).Encode(sendRawMessageResponse)
	} else {
		//log.Printf("json rpc incoming: %v\n", r)
		resp, err := http.Post(h.ethAddr, "application/json", &buf)
		if err != nil {
			json.NewEncoder(w).Encode(fmt.Errorf("http.Post() = %v", err))
		}

		//log.Printf("json rpc response: %v", resp)
		defer resp.Body.Close()
		buf2, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			json.NewEncoder(w).Encode(fmt.Errorf("ioutil.ReadAll() = %v", err))
		}
		w.Write(buf2)
	}
}
