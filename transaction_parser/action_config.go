package transaction_parser

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
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
