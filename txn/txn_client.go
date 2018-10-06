package txn

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
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

func (t *TxnClient) SendTransaction(ctx context.Context, rawTxn *types.Transaction) error {
	return t.client.SendTransaction(ctx, rawTxn)
}
