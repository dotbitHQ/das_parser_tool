package parser

import (
	"das_parser_tool/chain"
	"das_parser_tool/config"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/DeAccountSystems/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type Parser struct {
	dasCore   *core.DasCore
	ckbClient *chain.Client
}

type ParamsParser struct {
	DasCore   *core.DasCore
	CkbClient *chain.Client
}

func NewParser(p ParamsParser) *Parser {
	return &Parser{
		dasCore:   p.DasCore,
		ckbClient: p.CkbClient,
	}
}

func (t *Parser) HashParser(h string) map[string]interface{} {
	tx := t.ckbClient.GetTransactionByHash(types.HexToHash(h))
	// Warn: if you need order json, use ordered map
	return map[string]interface{}{
		"hash":      tx.Transaction.Hash.Hex(),
		"version":   tx.Transaction.Version,
		"cell_deps": t.parserCellDeps(tx.Transaction.CellDeps),
		"inputs":    t.parserInputs(tx.Transaction.Inputs),
		"outputs":   t.parserOutputs(tx.Transaction.Outputs, tx.Transaction.OutputsData),
		"witnesses": t.parserWitnesses(tx.Transaction),
		"status":    tx.TxStatus.Status,
	}
}

func (t *Parser) parserCellDeps(cellDeps []*types.CellDep) (cellDepsMap []interface{}) {
	for _, v := range cellDeps {
		if v.DepType == types.DepTypeDepGroup {
			cellDepsMap = append(cellDepsMap, map[string]interface{}{
				"secp256k1_blake160": v.OutPoint,
			})
			continue
		}
		if cellDep, ok := config.Cfg.DasCore.CellDeps[v.OutPoint.TxHash.String()]; ok {
			cellDepsMap = append(cellDepsMap, map[string]interface{}{
				cellDep: v.OutPoint,
			})
			continue
		}

		res := t.ckbClient.GetTransactionByHash(v.OutPoint.TxHash)
		output := res.Transaction.Outputs[v.OutPoint.Index]
		if output.Type != nil && output.Type.CodeHash.Hex() == config.Cfg.DasCore.THQCodeHash {
			switch common.Bytes2Hex(output.Type.Args) {
			case common.ArgsQuoteCell:
				cell, _ := t.dasCore.GetQuoteCell()
				cellDepsMap = append(cellDepsMap, map[string]interface{}{
					"QuoteCell": v.OutPoint,
					"quote":     cell.Quote(),
				})
			case common.ArgsTimeCell:
				cell, _ := t.dasCore.GetTimeCell()
				cellDepsMap = append(cellDepsMap, map[string]interface{}{
					"TimeCell":  v.OutPoint,
					"timestamp": cell.Timestamp(),
				})
			case common.ArgsHeightCell:
				cell, _ := t.dasCore.GetHeightCell()
				cellDepsMap = append(cellDepsMap, map[string]interface{}{
					"HeightCell":   v.OutPoint,
					"block_number": cell.BlockNumber(),
				})
			}
			continue
		}

		if output.Type != nil && output.Type.CodeHash.Hex() == config.Cfg.DasCore.DasContractCodeHash {
			script := common.GetScript(config.Cfg.DasCore.DasContractCodeHash, common.Bytes2Hex(output.Type.Args))
			if contractName, ok := core.DasContractByTypeIdMap[common.ScriptToTypeId(script).String()]; ok {
				cellDepsMap = append(cellDepsMap, map[string]interface{}{
					string(contractName): v.OutPoint,
					"output":             t.convertOutputTypeScript(output),
				})
				continue
			}
		}

		if output.Type != nil && output.Type.CodeHash.Hex() == config.Cfg.DasCore.DasConfigCodeHash {
			if value, ok := core.DasConfigCellMap.Load(common.Bytes2Hex(output.Type.Args)); ok {
				cellDepsMap = append(cellDepsMap, map[string]interface{}{
					value.(*core.DasConfigCellInfo).Name: v.OutPoint,
					"output":                             t.convertOutputTypeScript(output),
					"witness_hash":                       common.Bytes2Hex(res.Transaction.OutputsData[v.OutPoint.Index]),
				})
				continue
			}
		}

		cellDepsMap = append(cellDepsMap, map[string]interface{}{
			"unknown": v.OutPoint,
		})
	}
	return
}

func (t *Parser) parserInputs(inputs []*types.CellInput) (inputsMap []interface{}) {
	for _, v := range inputs {
		res := t.ckbClient.GetTransactionByHash(v.PreviousOutput.TxHash)
		inputsMap = append(inputsMap, t.parserOutput(res.Transaction.Outputs[v.PreviousOutput.Index], res.Transaction.OutputsData[v.PreviousOutput.Index]))
	}
	return
}

func (t *Parser) parserOutputs(outputs []*types.CellOutput, outputsData [][]byte) (outputsMap []interface{}) {
	for k, v := range outputs {
		outputsMap = append(outputsMap, t.parserOutput(v, outputsData[k]))
	}
	return
}

func (t *Parser) parserOutput(output *types.CellOutput, outputData []byte) (outputMap interface{}) {
	if output.Type == nil {
		return map[string]interface{}{
			"normal-cell": map[string]interface{}{
				"output":      t.convertOutputLockScript(output),
				"output_data": common.Bytes2Hex(outputData),
			},
		}
	}

	if contractName, ok := core.DasContractByTypeIdMap[output.Type.CodeHash.Hex()]; ok {
		switch string(contractName) {
		case "account-cell-type":
			id, _ := common.OutputDataToAccountId(outputData)
			next, _ := common.GetAccountCellNextAccountIdFromOutputData(outputData)
			expiredAt, _ := common.GetAccountCellExpiredAtFromOutputData(outputData)
			return map[string]interface{}{
				string(contractName): map[string]interface{}{
					"output_account": map[string]interface{}{
						"account":    string(outputData[80:]), // Warn: can't convert empty account
						"id":         common.Bytes2Hex(id),
						"next":       common.Bytes2Hex(next),
						"expired_at": expiredAt,
					},
					"output":      t.convertOutputTypeScript(output),
					"output_data": common.Bytes2Hex(outputData),
				},
			}
		default:
			return map[string]interface{}{
				string(contractName): map[string]interface{}{
					"output":      t.convertOutputTypeScript(output),
					"output_data": common.Bytes2Hex(outputData),
				},
			}
		}
	}

	return map[string]interface{}{
		"unknown": t.convertOutputTypeScript(output),
	}
}

func (t *Parser) convertOutputLockScript(output *types.CellOutput) map[string]interface{} {
	return map[string]interface{}{
		"capacity": output.Capacity,
		"lock": map[string]interface{}{
			"code_hash": output.Lock.CodeHash.Hex(),
			"hash_type": output.Lock.HashType,
			"args":      common.Bytes2Hex(output.Lock.Args),
		},
	}
}

func (t *Parser) convertOutputTypeScript(output *types.CellOutput) map[string]interface{} {
	return map[string]interface{}{
		"capacity": output.Capacity,
		"lock": map[string]interface{}{
			"code_hash": output.Lock.CodeHash.Hex(),
			"hash_type": output.Lock.HashType,
			"args":      common.Bytes2Hex(output.Lock.Args),
		},
		"type": map[string]interface{}{
			"code_hash": output.Type.CodeHash.Hex(),
			"hash_type": output.Type.HashType,
			"args":      common.Bytes2Hex(output.Type.Args),
		},
	}
}

func (t *Parser) parserWitnesses(transaction *types.Transaction) (witnessesMap []interface{}) {
	for _, witnessByte := range transaction.Witnesses {
		witnessesMap = append(witnessesMap, witness.ParserWitnessData(witnessByte))
	}

	return
}
