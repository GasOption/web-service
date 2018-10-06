package txn

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type HexTxnList struct {
	List []string `json:"list"`
}

func Decode(hexTxn string) (*types.Transaction, error) {
	var rawTxn *types.Transaction
	rawTxBytes, err := hex.DecodeString(hexTxn)
	if err != nil {
		return nil, err
	}
	rlp.DecodeBytes(rawTxBytes, &rawTxn)
	return rawTxn, nil
}
