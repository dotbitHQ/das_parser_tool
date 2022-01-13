package transaction_parser

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
)

func (b *TransactionParser) ActionTransferBalance(req FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
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

	return
}
