package txn

import (
	"encoding/hex"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type HexTxnList struct {
	List []string `json:"list"`
}

func Decode(hexTxn string) (*types.Transaction, error) {
	if strings.HasPrefix(hexTxn, "0x") {
		hexTxn = hexTxn[2:]
	}

	var rawTxn *types.Transaction
	rawTxBytes, err := hex.DecodeString(hexTxn)
	if err != nil {
		return nil, err
	}
	rlp.DecodeBytes(rawTxBytes, &rawTxn)
	return rawTxn, nil
}
