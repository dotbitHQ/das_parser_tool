package chain

import (
	"fmt"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/transaction"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/spf13/cobra"
)

func (c *Client) GetTransactionByHash(hash types.Hash) *types.TransactionWithStatus {
	res, err := c.client.GetTransaction(c.ctx, hash)
	if err != nil {
		cobra.CheckErr(fmt.Errorf("GetTransaction err: %v", err.Error()))
	}
	return res
}

func (c *Client) GetBalance(addr string) uint64 {
	parseAddr, err := address.Parse(addr)
	if err != nil {
		cobra.CheckErr(fmt.Errorf("Parse err: %v ", err.Error()))
	}
	searchKey := &indexer.SearchKey{
		Script: &types.Script{
			CodeHash: types.HexToHash(transaction.SECP256K1_BLAKE160_SIGHASH_ALL_TYPE_HASH),
			HashType: types.HashTypeType,
			Args:     parseAddr.Script.Args,
		},
		ScriptType: indexer.ScriptTypeLock,
		Filter: &indexer.CellsFilter{
			OutputDataLenRange: &[2]uint64{0, 1},
		},
	}
	res, err := c.client.GetCellsCapacity(c.ctx, searchKey)
	if err != nil {
		cobra.CheckErr(fmt.Errorf("GetCellsCapacity err: %v", err.Error()))
	}

	return res.Capacity
}
