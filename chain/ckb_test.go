package chain

import (
	"context"
	"das_parser_tool/config"
	"encoding/hex"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/scorpiotzh/toolib"
	"testing"
	"time"
)

func getCkbClient() (*Client, error) {
	if err := config.InitCfg("../config/config.yaml"); err != nil {
		panic(fmt.Errorf("InitCfg err: %s", err))
	}
	return NewClient(context.Background(), config.Cfg.Chain.CkbUrl, config.Cfg.Chain.IndexUrl)
}

func TestGetBalance(t *testing.T) {
	client, err := getCkbClient()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("111")
	fmt.Println(client.GetBalance("ckt1qyq27z6pccncqlaamnh8ttapwn260egnt67ss2cwvz"))
	fmt.Println(client.GetNormalLiveCell("ckt1qyq27z6pccncqlaamnh8ttapwn260egnt67ss2cwvz", 50999994254))
}

func TestGetCells(t *testing.T) {
	codeHash := "0x96248cdefb09eed910018a847cfb51ad044c2d7db650112931760e3ef34a7e9a"
	client, err := getCkbClient()
	if err != nil {
		t.Fatal(err)
	}
	searchKey := &indexer.SearchKey{
		Script:     common.GetScript(codeHash, "0x01"),
		ScriptType: indexer.ScriptTypeType,
	}
	res, err := client.Client().GetCells(context.Background(), searchKey, indexer.SearchOrderDesc, 100, "")
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range res.Objects {
		fmt.Println(v.OutPoint.TxHash.Hex(), v.OutPoint.Index)
	}
}

func TestGetBlockByNumber(t *testing.T) {
	client, err := getCkbClient()
	if err != nil {
		t.Fatal(err)
	}
	block, err := client.GetBlockByNumber(1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(block.Header.Timestamp)
	fmt.Println(time.Now().UnixNano() / 1e6)
	fmt.Println(time.Unix(int64(block.Header.Timestamp/1e3), 0).String())
}

func TestGetTransaction(t *testing.T) {
	hash := "0x5e594f15662fc75fe01fd67c76cee02c79dea2e6573509e5408af5114afd459e"
	client, err := getCkbClient()
	if err != nil {
		t.Fatal(err)
	}
	tx, err := client.Client().GetTransaction(context.Background(), types.HexToHash(hash))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(toolib.JsonString(tx))
}

func TestGetTransactions(t *testing.T) {
	client, err := getCkbClient()
	if err != nil {
		t.Fatal(err)
	}
	key := indexer.SearchKey{
		Script:     common.GetScript("0xbf43c3602455798c1a61a596e0d95278864c552fafe231c063b3fabf97a8febc", "0x26e9aa4899003b28de08c00a6e946d422c18bbba"),
		ScriptType: indexer.ScriptTypeLock,
		Filter: &indexer.CellsFilter{
			BlockRange: &[2]uint64{5377315, 5377442},
		},
	}
	list, err := client.Client().GetTransactions(context.Background(), &key, indexer.SearchOrderDesc, 1000, "")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(list.Objects))
	var mapTx = make(map[string]types.Hash)
	for _, v := range list.Objects {
		mapTx[v.TxHash.Hex()] = v.TxHash
	}
	fmt.Println(len(mapTx))
	count := 0
	for k, v := range mapTx {
		tx, err := client.Client().GetTransaction(context.Background(), v)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			if len(tx.Transaction.Outputs) > 0 {
				if "6d91285768e7c96f1cea0173c8167ada2cfeabe8" == hex.EncodeToString(tx.Transaction.Outputs[0].Lock.Args) {
					count++
					fmt.Println(count, k, string(tx.Transaction.OutputsData[0]))
				}
			}
		}
	}

}
