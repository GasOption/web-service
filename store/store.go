// Package store provides in-memory storage for the web service.
package store

import (
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

func (s *Store) AddToPool(rawTxn *types.Transaction) error {
	if _, ok := s.pool[rawTxn.Nonce()]; !ok {
		s.pool[rawTxn.Nonce()] = &TxnBundle{Processed: false}
	}
	s.pool[rawTxn.Nonce()].Bundle = append(s.pool[rawTxn.Nonce()].Bundle, &Txn{
		Submitted: false,
		GasPrice:  rawTxn.Gas(),
		HexHash:   rawTxn.Hash().Hex(),
		RawTxn:    rawTxn,
	})
	return nil
}

func (s *Store) Pool() map[uint64]*TxnBundle {
	return s.pool
}
