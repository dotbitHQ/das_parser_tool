package transaction_parser

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/DeAccountSystems/das-lib/witness"
)

func (t *TransactionParser) ActionTransferBalance(req FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	dasLock, err := core.GetDasContractInfo(common.DasContractNameDispatchCellType)
	if err != nil {
		resp.Err = fmt.Errorf("GetDasContractInfo err: %s", err.Error())
		return
	}

	dasBalance, err := core.GetDasContractInfo(common.DasContractNameBalanceCellType)
	if err != nil {
		resp.Err = fmt.Errorf("GetDasContractInfo err: %s", err.Error())
		return
	}

	log.Info("ActionTransferBalance:", req.Transaction.Hash.Hex())

	for _, v := range req.Transaction.Outputs {
		if v.Lock.CodeHash.Hex() != dasLock.ContractTypeId.Hex() {
			continue
		}
		if v.Type != nil && v.Type.CodeHash.Hex() != dasBalance.ContractTypeId.Hex() {
			continue
		}

		resp.ActionName = req.Builder.Action
	}
	if resp.ActionName == "" {
		return
	}

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
