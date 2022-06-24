package parser

import (
	"das_parser_tool/chain"
	"das_parser_tool/config"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"strings"
)

type Parser struct {
	CkbClient *chain.Client
	DasCore   *core.DasCore
}

func (t *Parser) HashParser(hash string) map[string]interface{} {
	tx := t.CkbClient.GetTransactionByHash(types.HexToHash(hash))
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

func (t *Parser) JsonParser(transaction *types.Transaction) map[string]interface{} {
	return map[string]interface{}{
		"hash":      transaction.Hash.Hex(),
		"version":   transaction.Version,
		"cell_deps": t.parserCellDeps(transaction.CellDeps),
		"inputs":    t.parserInputs(transaction.Inputs),
		"outputs":   t.parserOutputs(transaction.Outputs, transaction.OutputsData),
		"witnesses": t.parserWitnesses(transaction),
	}
}

var soScriptType = []common.SoScriptType{common.SoScriptTypeEth, common.SoScriptTypeTron, common.SoScriptTypeCkbSingle, common.SoScriptTypeCkbMulti, common.SoScriptTypeEd25519}

func (t *Parser) parserCellDeps(cellDeps []*types.CellDep) (cellDepsMap []interface{}) {
	var soScriptMap = map[string]string{}
	for _, soScriptName := range soScriptType {
		soScript, _ := core.GetDasSoScript(soScriptName)
		if soScript != nil {
			soScriptMap[soScript.OutPoint.TxHash.String()] = fmt.Sprintf("%s-lib", strings.ToLower(string(soScriptName)))
		}
	}
	for _, v := range cellDeps {
		if v.DepType == types.DepTypeDepGroup {
			cellDepsMap = append(cellDepsMap, map[string]interface{}{
				"name":     "secp256k1_blake160",
				"cell_dep": v.OutPoint,
			})
			continue
		}
		if cellDep, ok := soScriptMap[v.OutPoint.TxHash.String()]; ok {
			cellDepsMap = append(cellDepsMap, map[string]interface{}{
				"name":     cellDep,
				"cell_dep": v.OutPoint,
			})
			continue
		}

		cellDepsMap = append(cellDepsMap, t.parserCellDep(v.OutPoint))
	}
	return
}

func (t *Parser) parserCellDep(outpoint *types.OutPoint) interface{} {
	res := t.CkbClient.GetTransactionByHash(outpoint.TxHash)
	output := res.Transaction.Outputs[outpoint.Index]
	if output.Type == nil {
		return map[string]interface{}{
			"name":     "unknown",
			"cell_dep": outpoint,
		}
	}
	configCodeHash := common.ScriptToTypeId(&types.Script{
		CodeHash: types.HexToHash(config.Env.ContractCodeHash),
		HashType: types.HashTypeType,
		Args:     common.Hex2Bytes(config.Env.MapContract[common.DasContractNameConfigCellType]),
	})

	switch output.Type.CodeHash.Hex() {
	case config.Env.THQCodeHash:
		switch common.Bytes2Hex(output.Type.Args) {
		case common.ArgsQuoteCell:
			cell, _ := t.DasCore.GetQuoteCell()
			return map[string]interface{}{
				"name":     "QuoteCell",
				"cell_dep": outpoint,
				"quote":    cell.Quote(),
			}
		case common.ArgsTimeCell:
			cell, _ := t.DasCore.GetTimeCell()
			return map[string]interface{}{
				"name":      "TimeCell",
				"cell_dep":  outpoint,
				"timestamp": witness.ConvertTimestamp(cell.Timestamp()),
			}
		case common.ArgsHeightCell:
			cell, _ := t.DasCore.GetHeightCell()
			return map[string]interface{}{
				"name":         "HeightCell",
				"cell_dep":     outpoint,
				"block_number": cell.BlockNumber(),
			}
		}
	case config.Env.ContractCodeHash:
		script := common.GetScript(config.Env.ContractCodeHash, common.Bytes2Hex(output.Type.Args))
		if contractName, ok := core.DasContractByTypeIdMap[common.ScriptToTypeId(script).String()]; ok {
			return map[string]interface{}{
				"name":     string(contractName),
				"cell_dep": outpoint,
				"output":   t.convertOutputTypeScript(output),
			}
		}
	case configCodeHash.String():
		if value, ok := core.DasConfigCellMap.Load(common.Bytes2Hex(output.Type.Args)); ok {
			return map[string]interface{}{
				"name":         value.(*core.DasConfigCellInfo).Name,
				"cell_dep":     outpoint,
				"output":       t.convertOutputTypeScript(output),
				"witness_hash": common.Bytes2Hex(res.Transaction.OutputsData[outpoint.Index]),
			}
		}
	default:
		if contractName, ok := core.DasContractByTypeIdMap[output.Type.CodeHash.Hex()]; ok {
			return map[string]interface{}{
				"name":     string(contractName),
				"cell_dep": outpoint,
				"output":   t.convertOutputTypeScript(output),
			}
		}
	}

	return map[string]interface{}{
		"name":     "unknown",
		"cell_dep": outpoint,
	}
}

func (t *Parser) parserInputs(inputs []*types.CellInput) (inputsMap []interface{}) {
	for _, v := range inputs {
		res := t.CkbClient.GetTransactionByHash(v.PreviousOutput.TxHash)
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
			"name":        "normal-cell",
			"output":      t.convertOutputLockScript(output),
			"output_data": common.Bytes2Hex(outputData),
		}
	}

	if contractName, ok := core.DasContractByTypeIdMap[output.Type.CodeHash.Hex()]; ok {
		switch string(contractName) {
		case "account-cell-type":
			id, _ := common.OutputDataToAccountId(outputData)
			next, _ := common.GetAccountCellNextAccountIdFromOutputData(outputData)
			expiredAt, _ := common.GetAccountCellExpiredAtFromOutputData(outputData)
			return map[string]interface{}{
				"name": string(contractName),
				"output_account": map[string]interface{}{
					"account":    string(outputData[80:]), // Warn: can't convert empty account
					"id":         common.Bytes2Hex(id),
					"next":       common.Bytes2Hex(next),
					"expired_at": witness.ConvertTimestamp(int64(expiredAt)),
				},
				"output":      t.convertOutputTypeScript(output),
				"output_data": common.Bytes2Hex(outputData),
			}
		default:
			return map[string]interface{}{
				"name":   string(contractName),
				"output": t.convertOutputTypeScript(output),
			}
		}
	}

	return map[string]interface{}{
		"name":   "unknown",
		"output": t.convertOutputTypeScript(output),
	}
}

func (t *Parser) convertOutputLockScript(output *types.CellOutput) map[string]interface{} {
	return map[string]interface{}{
		"capacity": witness.ConvertCapacity(output.Capacity),
		"lock": map[string]interface{}{
			"code_hash": output.Lock.CodeHash.Hex(),
			"hash_type": output.Lock.HashType,
			"args":      common.Bytes2Hex(output.Lock.Args),
		},
	}
}

func (t *Parser) convertOutputTypeScript(output *types.CellOutput) map[string]interface{} {
	return map[string]interface{}{
		"capacity": witness.ConvertCapacity(output.Capacity),
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
