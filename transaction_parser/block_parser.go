package transaction_parser

import (
	"context"
	"das_parser_tool/chain/chain_ckb"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/DeAccountSystems/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/scorpiotzh/mylog"
	"sync"
)

var log = mylog.NewLogger("transaction_parser", mylog.LevelDebug)

type TransactionParser struct {
	dasCore              *core.DasCore
	mapTransactionHandle map[common.DasAction]FuncTransactionHandle
	ckbClient            *chain_ckb.Client
	ctx                  context.Context
	wg                   *sync.WaitGroup
}

type ParamsTransactionParser struct {
	DasCore            *core.DasCore
	CkbClient          *chain_ckb.Client
	Ctx                context.Context
	Wg                 *sync.WaitGroup
}

func NewTransactionParser(p ParamsTransactionParser) (*TransactionParser, error) {
	bp := TransactionParser{
		dasCore:            p.DasCore,
		ckbClient:          p.CkbClient,
		ctx:                p.Ctx,
		wg:                 p.Wg,
	}
	bp.registerTransactionHandle()
	return &bp, nil
}

func (b *TransactionParser) GetMapTransactionHandle(action common.DasAction) (FuncTransactionHandle, bool) {
	handler, ok := b.mapTransactionHandle[action]
	return handler, ok
}

func (b *TransactionParser) RunParser(t string) {
	tx, err := b.ckbClient.GetTransactionByHash(types.HexToHash(t))
	if err != nil {
		log.Fatal(err)
	}
	log.Info("parsingTransactionData txHash:", t)

	if builder, err := witness.ActionDataBuilderFromTx(tx.Transaction); err != nil {
		log.Warn("ActionDataBuilderFromTx err:", err.Error())
	} else {
		if handle, ok := b.mapTransactionHandle[builder.Action]; ok {
			// transaction parse by action
			resp := handle(FuncTransactionHandleReq{
				Tx:             tx.Transaction,
				Action:         builder.Action,
			})
			if resp.Err != nil {
				log.Error("action handle resp:", builder.Action, resp.Err.Error())
			}
		}

	}
}
