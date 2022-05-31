package chain

import (
	"context"
	"das_parser_tool/config"
	"encoding/json"
	"fmt"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"testing"
)

func getCkbClient() *Client {
	config.InitCfg("")
	return NewClient(context.Background(), config.Cfg.Chain.CkbUrl, config.Cfg.Chain.IndexUrl)
}

func TestGetTransaction(t *testing.T) {
	hash := "0x5e594f15662fc75fe01fd67c76cee02c79dea2e6573509e5408af5114afd459e"
	client := getCkbClient()
	tx := client.GetTransactionByHash(types.HexToHash(hash))
	s, _ := json.Marshal(tx)
	fmt.Println(string(s))
}
