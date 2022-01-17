package transaction_parser

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
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

	log.Info("ActionTransferBalance:", req.Hash)

	for _, v := range req.Tx.Outputs {
		if v.Lock.CodeHash.Hex() != dasLock.ContractTypeId.Hex() {
			continue
		}
		if v.Type != nil && v.Type.CodeHash.Hex() != dasBalance.ContractTypeId.Hex() {
			continue
		}

		resp.ActionName = req.Action
	}

	if resp.ActionName != "" {
		inputsSize := len(req.Tx.Inputs)
		for k, v := range req.Tx.Witnesses {
			if k < inputsSize {
				resp.WitnessesMap = append(resp.WitnessesMap, map[string]interface{}{
					"name":    "unknown",
					"witness": common.Bytes2Hex(v),
				})
				continue
			}

			if k == inputsSize {
				// TODO action and params should in witness map
				// log.Info("parserWitnesses action", builder.Action)
				// log.Info("parserWitnesses params", builder.ParamsStr)

				resp.WitnessesMap = append(resp.WitnessesMap, map[string]interface{}{
					"name":    "action_data",
					"witness": common.Bytes2Hex(v),
				})
				continue
			}

			resp.WitnessesMap = append(resp.WitnessesMap, map[string]interface{}{
				"name":    req.Action,
				"witness": common.Bytes2Hex(v),
			})
			continue
		}
	}

	return
}
