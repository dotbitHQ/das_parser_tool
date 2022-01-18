package transaction_parser

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
)

func (t *TransactionParser) ActionTransferPayment(req FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	dasLock, err := core.GetDasContractInfo(common.DasContractNameDispatchCellType)
	if err != nil {
		resp.Err = fmt.Errorf("GetDasContractInfo err: %s", err.Error())
		return
	}

	balanceType, err := core.GetDasContractInfo(common.DasContractNameBalanceCellType)
	if err != nil {
		resp.Err = fmt.Errorf("GetDasContractInfo err: %s", err.Error())
		return
	}

	log.Info("ActionTransferPayment:", req.Transaction.Hash.Hex())

	res, err := t.ckbClient.GetTransactionByHash(req.Transaction.Inputs[0].PreviousOutput.TxHash)
	if err != nil {
		resp.Err = fmt.Errorf("GetTransactionByHash err: %s", err.Error())
		return
	}
	cellOutput := res.Transaction.Outputs[req.Transaction.Inputs[0].PreviousOutput.Index]
	if !dasLock.IsSameTypeId(cellOutput.Lock.CodeHash) {
		log.Warn("ActionTransferPayment: das lock not match", req.Transaction.Hash.Hex())
		return
	}
	if cellOutput.Type != nil && !balanceType.IsSameTypeId(cellOutput.Type.CodeHash) {
		log.Warn("ActionTransferPayment: balance type not match", req.Transaction.Hash.Hex())
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
