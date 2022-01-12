package chain

import (
	"fmt"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/nervosnetwork/ckb-sdk-go/collector"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/transaction"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

func (c *Client) GetTipBlockNumber() (uint64, error) {
	if blockNumber, err := c.client.GetTipBlockNumber(c.ctx); err != nil {
		return 0, fmt.Errorf("GetTipBlockNumber err:%s", err.Error())
	} else {
		return blockNumber, nil
	}
}

func (c *Client) GetBlockByNumber(number uint64) (*types.Block, error) {
	return c.client.GetBlockByNumber(c.ctx, number)
}

func (c *Client) GetTransactionByHash(hash types.Hash) (*types.TransactionWithStatus, error) {
	if res, err := c.client.GetTransaction(c.ctx, hash); err != nil {
		return nil, fmt.Errorf("GetTransaction err:%s", err.Error())
	} else {
		return res, nil
	}
}

func (c *Client) GetHeaderByHash(hash types.Hash) (*types.Header, error) {
	if res, err := c.client.GetHeader(c.ctx, hash); err != nil {
		return nil, fmt.Errorf("GetHeader err:%s", err.Error())
	} else {
		return res, nil
	}
}

func (c *Client) GetNormalLiveCell(addr string, limit uint64) ([]*indexer.LiveCell, uint64, error) {
	parseAddr, err := address.Parse(addr)
	if err != nil {
		return nil, 0, fmt.Errorf("address.Parse err: %s", err.Error())
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
	co := collector.NewLiveCellCollector(c.client, searchKey, indexer.SearchOrderAsc, indexer.SearchLimit, "")
	iterator, err := co.Iterator()
	if err != nil {
		return nil, 0, fmt.Errorf("iterator err:%s", err.Error())
	}
	var cells []*indexer.LiveCell
	total := uint64(0)
	for iterator.HasNext() {
		liveCell, err := iterator.CurrentItem()
		if err != nil {
			return nil, 0, fmt.Errorf("CurrentItem err:%s", err.Error())
		}
		//fmt.Println(liveCell)
		cells = append(cells, liveCell)
		total += liveCell.Output.Capacity
		if limit > 0 && (total == limit || total-limit > uint64(6100000000)) { // limit 为转账金额+手续费
			break
		}
		if err = iterator.Next(); err != nil {
			return nil, 0, fmt.Errorf("next err:%s", err.Error())
		}
	}
	return cells, total, nil
}

func (c *Client) GetBalance(addr string) (uint64, error) {
	parseAddr, err := address.Parse(addr)
	if err != nil {
		return 0, fmt.Errorf("address.Parse err: %s", err.Error())
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
		return 0, fmt.Errorf("GetCellsCapacity err: %s", err.Error())
	}
	return res.Capacity, nil
}
