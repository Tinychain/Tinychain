package executor

import (
	"github.com/tinychain/tinychain/common"
	"github.com/tinychain/tinychain/core"
	"github.com/tinychain/tinychain/core/types"
	"github.com/tinychain/tinychain/db"
)

func (ex *Executor) commit(block *types.Block) error {
	if err := ex.persistTxs(block); err != nil {
		log.Errorf("failed to persist tx metas, err: %s", err)
		return err
	}

	if receipts, exist := ex.receiptsCache.Load(block.Height()); exist {
		err := ex.persistReceipts(block.Transactions, receipts.(types.Receipts), block.Height())
		if err != nil {
			log.Errorf("failed to persist receipts, err: %s", err)
			return err
		}
		ex.receiptsCache.Delete(block.Height())
	}

	if _, err := ex.stateCommit(block.Height()); err != nil {
		log.Errorf("failed to put state in batch, err: %s", err)
		return err
	}

	if err := ex.commitBlock(block); err != nil {
		log.Errorf("failed to append block to chain, err: %s", err)
		return err
	}

	// Commit data in batch
	if err := db.CommitBatch(ex.db.LDB(), block.Height()); err != nil {
		log.Errorf("failed to commit db.Batch, err: %s", err)
		return err
	}

	log.Infof("New block height = #%d commits. Hash = %s", block.Height(), block.Hash().Hex())
	go ex.event.Post(&core.CommitCompleteEvent{
		Block: block,
	})
	ex.resetVersion()
	return nil
}

// stateCommit commits the state transition at the given block height
func (ex *Executor) stateCommit(height uint64) (common.Hash, error) {
	return ex.state.Commit(db.GetBatch(ex.db.LDB(), height))
}

func (ex *Executor) persistTxs(block *types.Block) error {
	return ex.db.PutTxMetas(db.GetBatch(ex.db.LDB(), block.Height()), block.Transactions, block.Hash(), block.Height(), false, false)
}

func (ex *Executor) persistReceipts(txs types.Transactions, receipts types.Receipts, height uint64) error {
	if len(txs) != len(receipts) {
		return errReceiptNumInvalid
	}
	for idx, rp := range receipts {
		tx := txs[idx]
		if err := ex.db.PutReceipt(db.GetBatch(ex.db.LDB(), height), tx.Hash(), rp, false, false); err != nil {
			return err
		}
	}
	return nil
}

func (ex *Executor) commitBlock(block *types.Block) error {
	return ex.chain.CommitBlock(db.GetBatch(ex.db.LDB(), block.Height()), block)
}
