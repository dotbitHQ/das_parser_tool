package transaction_parser

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
)

func (t *TransactionParser) ActionCreateIncome(req FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Transaction, common.DasContractNameIncomeCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
		return
	} else if !isCV {
		log.Warn("not current version create income tx")
		return
	}

	log.Info("ActionCreateIncome:", req.Transaction.Hash.Hex())

	for k, witnessByte := range req.Transaction.Witnesses {
		if k < len(req.Transaction.Inputs) {
			resp.WitnessesMap = append(resp.WitnessesMap, t.parserNormalWitness(witnessByte))
			continue
		}

		if k == len(req.Transaction.Inputs) {
			resp.WitnessesMap = append(resp.WitnessesMap, t.parserActionDataWitness(witnessByte, req.Builder))
			continue
		}

		// TODO add das lib parser witness direct
		configCellMain := t.parserConfigCellMainWitnesses(witnessByte, req.Transaction)
		if configCellMain != nil {
			resp.WitnessesMap = append(resp.WitnessesMap, configCellMain)
			continue
		}

		configCellIncome := t.parserConfigCellIncomeWitnesses(witnessByte, req.Transaction)
		if configCellIncome != nil {
			resp.WitnessesMap = append(resp.WitnessesMap, configCellIncome)
			continue
		}

		// IncomeCellData
		incomeCellData := t.parserIncomeCellDataWitnesses(witnessByte, req.Transaction)
		if incomeCellData != nil {
			resp.WitnessesMap = append(resp.WitnessesMap, incomeCellData)
		}
	}

	return
}
