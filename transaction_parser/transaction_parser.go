package transaction_parser

import (
	"context"
	"das_parser_tool/chain"
	"das_parser_tool/config"
	"encoding/json"
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

func (t *TransactionParser) GetMapTransactionHandle(action common.DasAction) (FuncTransactionHandle, bool) {
	handler, ok := t.mapTransactionHandle[action]
	return handler, ok
}

func (t *TransactionParser) RunParser(h string) {
	tx, err := t.ckbClient.GetTransactionByHash(types.HexToHash(h))
	if err != nil {
		log.Fatal(err)
	}
	log.Info("RunParser txHash:", tx.Transaction.Hash)
	log.Info("RunParser version:", tx.Transaction.Version)
	log.Info("RunParser status:", tx.TxStatus.Status)

	builder, err := witness.ActionDataBuilderFromTx(tx.Transaction)
	if err != nil {
		log.Fatal("ActionDataBuilderFromTx err:", err.Error())
	}
	handle, ok := t.mapTransactionHandle[builder.Action]
	if !ok {
		log.Fatal("action doesn't exist", builder.Action)
	}
	// transaction parse by action
	resp := handle(FuncTransactionHandleReq{
		Tx:     tx.Transaction,
		Hash:   h,
		Action: builder.Action,
	})
	if resp.Err != nil {
		log.Fatal("action handle err:", builder.Action, resp.Err.Error())
	}

	if resp.ActionName != "" {
		log.Info("ActionName", resp.ActionName)
		t.parserTransaction(tx.Transaction)
	}
}

func (t *TransactionParser) parserTransaction(transaction *types.Transaction) {
	// Warn: if you need order json, use ordered map
	out := map[string]interface{}{}
	out["cell_deps"] = t.parserCellDeps(transaction.CellDeps)
	out["inputs"] = t.parserInputs(transaction.Inputs)
	out["outputs"] = t.parserOutputs(transaction.Outputs, transaction.OutputsData)

	b, _ := json.Marshal(out)
	log.Info(string(b))
}

func (t *TransactionParser) parserCellDeps(cellDeps []*types.CellDep) (cellDepsMap []interface{}) {
	for _, v := range cellDeps {
		if cellDep, ok := config.Cfg.DasCore.CellDeps[v.OutPoint.TxHash.String()]; ok {
			cellDepsMap = append(cellDepsMap, map[string]interface{}{
				"name": cellDep,
			})
			continue
		}
		res, err := t.ckbClient.GetTransactionByHash(v.OutPoint.TxHash)
		if err != nil {
			log.Fatal("GetTransactionByHash err:", err.Error())
		}

		output := res.Transaction.Outputs[v.OutPoint.Index]
		if output.Type.CodeHash.Hex() == config.Cfg.DasCore.THQCodeHash {
			switch common.Bytes2Hex(output.Type.Args) {
			case common.ArgsQuoteCell:
				cell, _ := t.dasCore.GetQuoteCell()
				cellDepsMap = append(cellDepsMap, map[string]interface{}{
					"name":  "QuoteCell",
					"quote": cell.Quote(),
				})
			case common.ArgsTimeCell:
				cell, _ := t.dasCore.GetTimeCell()
				cellDepsMap = append(cellDepsMap, map[string]interface{}{
					"name":      "TimeCell",
					"timestamp": cell.Timestamp(),
				})
			case common.ArgsHeightCell:
				cell, _ := t.dasCore.GetHeightCell()
				cellDepsMap = append(cellDepsMap, map[string]interface{}{
					"name":         "HeightCell",
					"block_number": cell.BlockNumber(),
				})
			}
			continue
		}

		if output.Type.CodeHash.Hex() == config.Cfg.DasCore.DasContractCodeHash {
			script := common.GetScript(config.Cfg.DasCore.DasContractCodeHash, common.Bytes2Hex(output.Type.Args))
			if contractName, ok := core.DasContractByTypeIdMap[common.ScriptToTypeId(script).String()]; ok {
				cellDepsMap = append(cellDepsMap, map[string]interface{}{
					"name": string(contractName),
					"type": t.convertOutputTypeScript(output),
				})
				continue
			}
		}

		if output.Type.CodeHash.Hex() == config.Cfg.DasCore.DasConfigCodeHash {
			if value, ok := core.DasConfigCellMap.Load(common.Bytes2Hex(output.Type.Args)); ok {
				cellDepsMap = append(cellDepsMap, map[string]interface{}{
					"name": value.(*core.DasConfigCellInfo).Name,
					"type": t.convertOutputTypeScript(output),
					"data": common.Bytes2Hex(res.Transaction.OutputsData[v.OutPoint.Index]),
				})
				continue
			}
		}

		cellDepsMap = append(cellDepsMap, map[string]interface{}{
			"name": "unknown",
		})
	}
	return
}

func (t *TransactionParser) parserInputs(inputs []*types.CellInput) (inputsMap []interface{}) {
	for _, v := range inputs {
		res, err := t.ckbClient.GetTransactionByHash(v.PreviousOutput.TxHash)
		if err != nil {
			log.Fatal("GetTransactionByHash err:", err.Error())
		}
		inputsMap = append(inputsMap, t.parserOutput(res.Transaction.Outputs[v.PreviousOutput.Index], res.Transaction.OutputsData[v.PreviousOutput.Index]))
	}
	return
}

func (t *TransactionParser) parserOutputs(outputs []*types.CellOutput, outputsData [][]byte) (outputsMap []interface{}) {
	for k, v := range outputs {
		outputsMap = append(outputsMap, t.parserOutput(v, outputsData[k]))
	}
	return
}

func (t *TransactionParser) parserOutput(output *types.CellOutput, outputData []byte) (outputMap interface{}) {
	if output.Type == nil {
		return map[string]interface{}{
			"name":     "NormalCell",
			"capacity": output.Capacity,
			"lock":     t.convertOutputLockScript(output),
			"data":     common.Bytes2Hex(outputData),
		}
	}

	if contractName, ok := core.DasContractByTypeIdMap[output.Type.CodeHash.Hex()]; ok {
		if string(contractName) == "account-cell-type" {
			accountId, _ := common.OutputDataToAccountId(outputData)
			return map[string]interface{}{
				"name":       string(contractName),
				"capacity":   output.Capacity,
				"account_id": common.Bytes2Hex(accountId),
				"account":    string(outputData[80:]), // Warn: can't convert empty account
				"type":       t.convertOutputTypeScript(output),
				"data":       common.Bytes2Hex(outputData),
			}
		}
		return map[string]interface{}{
			"name":     string(contractName),
			"capacity": output.Capacity,
			"type":     t.convertOutputTypeScript(output),
			"data":     common.Bytes2Hex(outputData),
		}
	}

	return map[string]interface{}{
		"name": "unknown",
	}
}

func (t *TransactionParser) convertOutputLockScript(output *types.CellOutput) map[string]interface{} {
	return map[string]interface{}{
		"code_hash": output.Lock.CodeHash,
		"hash_type": output.Lock.HashType,
		"args":      common.Bytes2Hex(output.Lock.Args),
	}
}

func (t *TransactionParser) convertOutputTypeScript(output *types.CellOutput) map[string]interface{} {
	return map[string]interface{}{
		"code_hash": output.Type.CodeHash,
		"hash_type": output.Type.HashType,
		"args":      common.Bytes2Hex(output.Type.Args),
	}
}
