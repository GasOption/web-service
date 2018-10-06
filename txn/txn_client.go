package txn

import (
	"context"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

type TxnClient struct {
	client *ethclient.Client
}

func NewClient(addr string) (*TxnClient, error) {
	client, err := ethclient.Dial(addr)
	if err != nil {
		return nil, err
	}
	return &TxnClient{client: client}, nil
}

func (t *TxnClient) SendTransaction(ctx context.Context, hexTxn string) error {
	var txn *types.Transaction

	rawTxBytes, err := hex.DecodeString(hexTxn)
	if err != nil {
		return err
	}

	rlp.DecodeBytes(rawTxBytes, &txn)

	return t.client.SendTransaction(context.Background(), txn)
}
