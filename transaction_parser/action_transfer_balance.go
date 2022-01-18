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

	log.Info("ActionTransferBalance:", req.Transaction.Hash.Hex())

	var currentVersion bool
	for _, v := range req.Transaction.Outputs {
		if v.Lock.CodeHash.Hex() != dasLock.ContractTypeId.Hex() {
			continue
		}
		if v.Type != nil && v.Type.CodeHash.Hex() != dasBalance.ContractTypeId.Hex() {
			continue
		}

		currentVersion = true
	}
	if !currentVersion {
		return
	}

	for k, witnessByte := range req.Transaction.Witnesses {
		if k < len(req.Transaction.Inputs) {
			resp.WitnessesMap = append(resp.WitnessesMap, t.parserNormalWitness(witnessByte))
			continue
		}

		if k == len(req.Transaction.Inputs) {
			resp.WitnessesMap = append(resp.WitnessesMap, t.parserActionDataWitness(witnessByte, req.Builder))
			continue
		}

		configCellMain := t.parserConfigCellMainWitnesses(witnessByte, req.Transaction)
		if configCellMain != nil {
			resp.WitnessesMap = append(resp.WitnessesMap, configCellMain)
		}
	}

	return
}
