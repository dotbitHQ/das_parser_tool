package transaction_parser

import (
	"context"
	"das_parser_tool/chain"
	"das_parser_tool/config"
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
	ckbClient            *chain.Client
	ctx                  context.Context
	wg                   *sync.WaitGroup
}

type ParamsTransactionParser struct {
	DasCore   *core.DasCore
	CkbClient *chain.Client
	Ctx       context.Context
	Wg        *sync.WaitGroup
}

func NewTransactionParser(p ParamsTransactionParser) (*TransactionParser, error) {
	bp := TransactionParser{
		dasCore:   p.DasCore,
		ckbClient: p.CkbClient,
		ctx:       p.Ctx,
		wg:        p.Wg,
	}
	bp.registerTransactionHandle()
	return &bp, nil
}

func (b *TransactionParser) GetMapTransactionHandle(action common.DasAction) (FuncTransactionHandle, bool) {
	handler, ok := b.mapTransactionHandle[action]
	return handler, ok
}

func (b *TransactionParser) RunParser(h string) {
	tx, err := b.ckbClient.GetTransactionByHash(types.HexToHash(h))
	if err != nil {
		log.Fatal(err)
	}
	log.Info("RunParser txHash:", h)
	log.Info("RunParser status:", tx.TxStatus.Status)

	builder, err := witness.ActionDataBuilderFromTx(tx.Transaction)
	if err != nil {
		log.Fatal("ActionDataBuilderFromTx err:", err.Error())
	}
	handle, ok := b.mapTransactionHandle[builder.Action]
	if !ok {
		log.Fatal("mapTransactionHandle does not exist", builder.Action)
	}
	// transaction parse by action
	resp := handle(FuncTransactionHandleReq{
		Tx:     tx.Transaction,
		Hash:   h,
		Action: builder.Action,
	})
	if resp.Err != nil {
		log.Fatal("action handle resp:", builder.Action, resp.Err.Error())
	}

	if resp.ActionName != "" {
		log.Info("ActionName", resp.ActionName)
		b.parserTransaction(tx.Transaction)
	}
}

func (b *TransactionParser) parserTransaction(transaction *types.Transaction) {
	var cellDeps []string
	for _, v := range transaction.CellDeps {
		if cellDep, ok := config.Cfg.DasCore.CellDeps[v.OutPoint.TxHash.String()]; ok {
			cellDeps = append(cellDeps, cellDep)
		} else {
			cellDeps = append(cellDeps, "unknown")
		}
	}

	log.Info(cellDeps)

	// TODO auto changed cell
	quote, _ := b.dasCore.GetQuoteCell()
	log.Info(quote.Quote(), quote.ToCellDep().OutPoint.TxHash)

	height, _ := b.dasCore.GetHeightCell()
	log.Info(height.BlockNumber(), height.ToCellDep().OutPoint.TxHash)

	time, _ := b.dasCore.GetTimeCell()
	log.Info(time.Timestamp(), time.ToCellDep().OutPoint.TxHash)

}
