// Package store provides in-memory storage for the web service.
package store

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
)

type Txn struct {
	Submitted bool
	GasPrice  uint64
	HexHash   string
	RawTxn    *types.Transaction
}

type TxnBundle struct {
	Processed bool
	Bundle    []*Txn
}

type Store struct {
	pool map[uint64]*TxnBundle
}

func New() *Store {
	pool := make(map[uint64]*TxnBundle)
	return &Store{pool: pool}
}

func (s *Store) AddToPool(rawTxns []*types.Transaction) error {
	if len(rawTxns) == 0 {
		return fmt.Errorf("empty rawTxns")
	}

	txnBundle := &TxnBundle{Processed: false}
	for _, rawTxn := range rawTxns {
		txnBundle.Bundle = append(txnBundle.Bundle, &Txn{
			Submitted: false,
			GasPrice:  rawTxn.Gas(),
			HexHash:   rawTxn.Hash().Hex(),
			RawTxn:    rawTxn,
		})
	}

	s.pool[rawTxns[0].Nonce()] = txnBundle

	return nil
}

func (s *Store) Pool() map[uint64]*TxnBundle {
	return s.pool
}
