package transaction_parser

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
)

func (t *TransactionParser) ActionWithdrawFromWallet(req FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
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

	log.Info("ActionWithdrawFromWallet:", req.Transaction.Hash.Hex())

	res, err := t.ckbClient.GetTransactionByHash(req.Transaction.Inputs[0].PreviousOutput.TxHash)
	if err != nil {
		resp.Err = fmt.Errorf("GetTransactionByHash err: %s", err.Error())
		return
	}
	cellOutput := res.Transaction.Outputs[req.Transaction.Inputs[0].PreviousOutput.Index]
	if !dasLock.IsSameTypeId(cellOutput.Lock.CodeHash) {
		log.Warn("ActionWithdrawFromWallet: das lock not match", req.Transaction.Hash.Hex())
		return
	}
	if cellOutput.Type != nil && !balanceType.IsSameTypeId(cellOutput.Type.CodeHash) {
		log.Warn("ActionWithdrawFromWallet: balance type not match", req.Transaction.Hash.Hex())
		return
	}

	resp.ActionName = req.Builder.Action
	resp.WitnessesMap = t.parserNormalWitnesses(req, len(req.Transaction.Inputs))

	return
}
