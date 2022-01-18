package transaction_parser

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/DeAccountSystems/das-lib/witness"
)

func (t *TransactionParser) ActionConfigCell(req FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	configContract, err := core.GetDasContractInfo(common.DasContractNameConfigCellType)
	if err != nil {
		resp.Err = fmt.Errorf("GetDasContractInfo err: %s", err.Error())
		return
	} else if configContract.ContractTypeId != req.Transaction.Outputs[0].Type.CodeHash {
		log.Warn("not current version config cell")
		return
	}

	log.Info("ActionConfigCell:", req.Transaction.Hash.Hex())
	// config cell update，rsync config cell out point
	if err = t.dasCore.AsyncDasConfigCell(); err != nil {
		resp.Err = fmt.Errorf("AsyncDasConfigCell err: %s", err.Error())
		return
	}

	resp.ActionName = req.Builder.Action

	inputsSize := len(req.Transaction.Inputs)
	for k, v := range req.Transaction.Witnesses {
		if k < inputsSize {
			resp.WitnessesMap = append(resp.WitnessesMap, map[string]interface{}{
				"name":    "unknown",
				"witness": common.Bytes2Hex(v),
			})
			continue
		}

		if k == inputsSize {
			resp.WitnessesMap = append(resp.WitnessesMap, map[string]interface{}{
				"name":         "ActionData",
				"witness":      common.Bytes2Hex(v),
				"witness_hash": common.Bytes2Hex(common.Blake2b(req.Builder.ActionData.AsSlice())),
				"action":       req.Builder.Action,
				"params":       req.Builder.ParamsStr,
			})
			continue
		}

		builder, _ := witness.ConfigCellDataBuilderByTypeArgs(req.Transaction, common.ConfigCellTypeArgsMain)
		if builder.ConfigCellMain != nil {
			resp.WitnessesMap = append(resp.WitnessesMap, t.parserConfigCellMain(v, builder.ConfigCellMain))
			continue
		}
	}

	return
}
