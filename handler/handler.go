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
	// Header.
	enableCors(&w)

	//params := mux.Vars(r)
	json.NewEncoder(w).Encode("hello")
}

// UpdateTxnList adds a list of transactions with different gas prices to transactions
// pool.
func (h *Handler) UpdateTxnList(w http.ResponseWriter, r *http.Request) {
	// Header.
	enableCors(&w)
	if (*r).Method == "OPTIONS" {
		return
	}

	// Parse POST body.
	var hexTxnList txn.HexTxnList
	if err := json.NewDecoder(r.Body).Decode(&hexTxnList); err != nil {
		json.NewEncoder(w).Encode(fmt.Errorf("json.NewDecoder().Decode() = %v", err))
	}
	log.Print(hexTxnList)

	// Decode hex transaction into raw transaction.
	hexTxn := hexTxnList.List[0]
	log.Printf("Decoding hex transaction: %v", hexTxn)
	rawTxn, err := txn.Decode(hexTxn)
	if err != nil {
		json.NewEncoder(w).Encode(fmt.Errorf("txn.Decode() = %v", err))
		return
	}

	// Add raw transactions to storage pool,
	if err := h.store.AddToPool(rawTxn); err != nil {
		json.NewEncoder(w).Encode(fmt.Errorf("store.AddToPool() = %v", err))
	}

	json.NewEncoder(w).Encode("200")
}

func (h *Handler) GetPool(w http.ResponseWriter, r *http.Request) {
	// Header.
	enableCors(&w)

	json.NewEncoder(w).Encode(h.store.Pool())
}

func (h *Handler) CreateJsonRpc(w http.ResponseWriter, r *http.Request) {
	// Header.
	enableCors(&w)
	if (*r).Method == "OPTIONS" {
		return
	}

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
		hexTxn := jsonRpcRequest.Params[0]
		log.Printf("Decoding hex transaction: %v", hexTxn)
		rawTxn, err := txn.Decode(hexTxn.(string))
		if err != nil {
			log.Printf("txn.Decode() = %v", err)
			//json.NewEncoder(w).Encode(fmt.Errorf("txn.Decode() = %v", err))
		}

		// Add raw transactions to storage pool,
		if err := h.store.AddToPool(rawTxn); err != nil {
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

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, OPTIONS, HEAD, CONNECT, TRACE, PATCH")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
}
