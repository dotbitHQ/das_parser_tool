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

	resp.ActionName = req.Builder.Action
	resp.WitnessesMap = t.parserNormalWitnesses(req, len(req.Transaction.Inputs))

	return
}
