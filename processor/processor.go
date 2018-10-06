package processor

import (
	"context"
	"log"
	"time"

	"github.com/GasOption/web-service/store"
	"github.com/GasOption/web-service/txn"
)

func Process(ctx context.Context, storeClient *store.Store, txnClient *txn.TxnClient) {
	for {
		time.Sleep(10 * time.Second)
		log.Printf("Processor waking up to work ...")

		for _, txnBundle := range storeClient.Pool() {
			if txnBundle.Processed || len(txnBundle.Bundle) == 0 {
				continue
			}

			for index, txn := range txnBundle.Bundle {
				// TODO: Implement an algorithm here to determine:
				// 1. When is the best time to submit the transaction.
				// 2. Which transaction in the bundle has the best gas price to be used.
				// Currently, we just naively pick the first one.
				if index > 0 {
					continue
				}

				log.Printf("Submitting the transaction with hash value: %v", txn.HexHash)
				if err := txnClient.SendTransaction(ctx, txn.RawTxn); err != nil {
					log.Fatalf("txnClient.SendTransaction() = %v", err)
				}

				txn.Submitted = true
			}
			txnBundle.Processed = true
		}
	}
}
