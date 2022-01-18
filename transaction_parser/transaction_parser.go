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

	// Warn: if you need order json, use ordered map
	out := map[string]interface{}{}
	out["cell_deps"] = t.parserCellDeps(tx.Transaction.CellDeps)
	out["inputs"] = t.parserInputs(tx.Transaction.Inputs)
	out["outputs"] = t.parserOutputs(tx.Transaction.Outputs, tx.Transaction.OutputsData)
	out["witness"] = t.parserWitnesses(tx.Transaction)

	b, _ := json.Marshal(out)
	log.Info(string(b))
}

func (t *TransactionParser) parserCellDeps(cellDeps []*types.CellDep) (cellDepsMap []interface{}) {
	for _, v := range cellDeps {
		if v.DepType == types.DepTypeDepGroup {
			cellDepsMap = append(cellDepsMap, map[string]interface{}{
				"name": "secp256k1_blake160",
			})
			continue
		}
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
					"name":         value.(*core.DasConfigCellInfo).Name,
					"type":         t.convertOutputTypeScript(output),
					"witness_hash": common.Bytes2Hex(res.Transaction.OutputsData[v.OutPoint.Index]),
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
			"name":     "normal-cell",
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

func (t *TransactionParser) parserWitnesses(transaction *types.Transaction) (witnessesMap []interface{}) {
	builder, err := witness.ActionDataBuilderFromTx(transaction)
	if err != nil {
		log.Fatal("ActionDataBuilderFromTx err:", err.Error())
	}
	handle, ok := t.mapTransactionHandle[builder.Action]
	if !ok {
		log.Fatal("action doesn't exist", builder.Action)
	}
	// transaction parse by action
	resp := handle(FuncTransactionHandleReq{
		Transaction: transaction,
		Builder:     builder,
	})
	if resp.Err != nil {
		log.Fatal("action handle err:", builder.Action, resp.Err.Error())
	}

	witnessesMap = resp.WitnessesMap
	return
}

func (t *TransactionParser) parserNormalWitnesses(req FuncTransactionHandleReq, inputsSize int) (witnessesMap []interface{}) {
	for k, v := range req.Transaction.Witnesses {
		if k < inputsSize {
			witnessesMap = append(witnessesMap, map[string]interface{}{
				"name":    "unknown",
				"witness": common.Bytes2Hex(v),
			})
			continue
		}

		if k == inputsSize {
			witnessesMap = append(witnessesMap, map[string]interface{}{
				"name":         "ActionData",
				"witness":      common.Bytes2Hex(v),
				"witness_hash": common.Bytes2Hex(common.Blake2b(req.Builder.ActionData.AsSlice())),
				"action":       req.Builder.Action,
				"action_hash":  common.Bytes2Hex(req.Builder.ActionData.Action().RawData()),
				"params":       req.Builder.ParamsStr,
			})
			continue
		}

		builder, _ := witness.ConfigCellDataBuilderByTypeArgs(req.Transaction, common.ConfigCellTypeArgsMain)
		if builder != nil && builder.ConfigCellMain != nil {
			witnessesMap = append(witnessesMap, map[string]interface{}{
				"name":         "ConfigCellMain",
				"witness":      common.Bytes2Hex(v),
				"witness_hash": common.Bytes2Hex(common.Blake2b(builder.ConfigCellMain.AsSlice())),
				"status":       common.Bytes2Hex(builder.ConfigCellMain.Status().RawData()),
				"type_id_table": map[string]interface{}{
					"account_cell":         common.Bytes2Hex(builder.ConfigCellMain.TypeIdTable().AccountCell().RawData()),
					"apply_register_cell":  common.Bytes2Hex(builder.ConfigCellMain.TypeIdTable().ApplyRegisterCell().RawData()),
					"balance_cell":         common.Bytes2Hex(builder.ConfigCellMain.TypeIdTable().BalanceCell().RawData()),
					"income_cell":          common.Bytes2Hex(builder.ConfigCellMain.TypeIdTable().IncomeCell().RawData()),
					"pre_account_cell":     common.Bytes2Hex(builder.ConfigCellMain.TypeIdTable().PreAccountCell().RawData()),
					"proposal_cell":        common.Bytes2Hex(builder.ConfigCellMain.TypeIdTable().ProposalCell().RawData()),
					"account_sale_cell":    common.Bytes2Hex(builder.ConfigCellMain.TypeIdTable().AccountSaleCell().RawData()),
					"account_auction_cell": common.Bytes2Hex(builder.ConfigCellMain.TypeIdTable().AccountAuctionCell().RawData()),
					"offer_cell":           common.Bytes2Hex(builder.ConfigCellMain.TypeIdTable().OfferCell().RawData()),
					"reverse_record_cell":  common.Bytes2Hex(builder.ConfigCellMain.TypeIdTable().ReverseRecordCell().RawData()),
				},
				"das_lock_out_point_table": map[string]interface{}{
					"ckb_signall": map[string]interface{}{
						"tx_hash": common.Bytes2Hex(builder.ConfigCellMain.DasLockOutPointTable().CkbSignall().TxHash().RawData()),
						"index":   common.Bytes2Hex(builder.ConfigCellMain.DasLockOutPointTable().CkbSignall().Index().RawData()),
					},
					"ckb_multisign": map[string]interface{}{
						"tx_hash": common.Bytes2Hex(builder.ConfigCellMain.DasLockOutPointTable().CkbMultisign().TxHash().RawData()),
						"index":   common.Bytes2Hex(builder.ConfigCellMain.DasLockOutPointTable().CkbMultisign().Index().RawData()),
					},
					"ckb_anyone_can_pay": map[string]interface{}{
						"tx_hash": common.Bytes2Hex(builder.ConfigCellMain.DasLockOutPointTable().CkbAnyoneCanPay().TxHash().RawData()),
						"index":   common.Bytes2Hex(builder.ConfigCellMain.DasLockOutPointTable().CkbAnyoneCanPay().Index().RawData()),
					},
					"eth": map[string]interface{}{
						"tx_hash": common.Bytes2Hex(builder.ConfigCellMain.DasLockOutPointTable().Eth().TxHash().RawData()),
						"index":   common.Bytes2Hex(builder.ConfigCellMain.DasLockOutPointTable().Eth().Index().RawData()),
					},
					"tron": map[string]interface{}{
						"tx_hash": common.Bytes2Hex(builder.ConfigCellMain.DasLockOutPointTable().Tron().TxHash().RawData()),
						"index":   common.Bytes2Hex(builder.ConfigCellMain.DasLockOutPointTable().Tron().Index().RawData()),
					},
				},
			})
			continue
		}
	}
	return
}
