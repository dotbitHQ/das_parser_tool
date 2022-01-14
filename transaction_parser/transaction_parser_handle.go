package transaction_parser

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

const DasActionTransferBalance = "transfer_balance"

func (t *TransactionParser) registerTransactionHandle() {
	t.mapTransactionHandle = make(map[string]FuncTransactionHandle)
	t.mapTransactionHandle[DasActionTransferBalance] = t.ActionTransferBalance
	t.mapTransactionHandle[common.DasActionConfig] = t.ActionConfigCell
	t.mapTransactionHandle[common.DasActionWithdrawFromWallet] = t.ActionWithdrawFromWallet
	t.mapTransactionHandle[common.DasActionTransfer] = t.ActionTransferPayment
	t.mapTransactionHandle[common.DasActionCreateIncome] = t.ActionCreateIncome
	t.mapTransactionHandle[common.DasActionConsolidateIncome] = t.ActionConsolidateIncome

	t.mapTransactionHandle[common.DasActionApplyRegister] = t.ActionApplyRegister
	t.mapTransactionHandle[common.DasActionPreRegister] = t.ActionPreRegister
	t.mapTransactionHandle[common.DasActionPropose] = t.ActionPropose
	t.mapTransactionHandle[common.DasActionExtendPropose] = t.ActionPropose
	t.mapTransactionHandle[common.DasActionConfirmProposal] = t.ActionConfirmProposal
	t.mapTransactionHandle[common.DasActionEditRecords] = t.ActionEditRecords
	t.mapTransactionHandle[common.DasActionEditManager] = t.ActionEditManager
	t.mapTransactionHandle[common.DasActionRenewAccount] = t.ActionRenewAccount
	t.mapTransactionHandle[common.DasActionTransferAccount] = t.ActionTransferAccount

	t.mapTransactionHandle[common.DasActionStartAccountSale] = t.ActionStartAccountSale
	t.mapTransactionHandle[common.DasActionEditAccountSale] = t.ActionEditAccountSale
	t.mapTransactionHandle[common.DasActionCancelAccountSale] = t.ActionCancelAccountSale
	t.mapTransactionHandle[common.DasActionBuyAccount] = t.ActionBuyAccount

	t.mapTransactionHandle[common.DasActionMakeOffer] = t.ActionMakeOffer
	t.mapTransactionHandle[common.DasActionEditOffer] = t.ActionEditOffer
	t.mapTransactionHandle[common.DasActionCancelOffer] = t.ActionCancelOffer
	t.mapTransactionHandle[common.DasActionAcceptOffer] = t.ActionAcceptOffer

	t.mapTransactionHandle[common.DasActionDeclareReverseRecord] = t.ActionDeclareReverseRecord
	t.mapTransactionHandle[common.DasActionRedeclareReverseRecord] = t.ActionRedeclareReverseRecord
	t.mapTransactionHandle[common.DasActionRetractReverseRecord] = t.ActionRetractReverseRecord
}

func isCurrentVersionTx(tx *types.Transaction, name common.DasContractName) (bool, error) {
	contract, err := core.GetDasContractInfo(name)
	if err != nil {
		return false, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	isCV := false
	for _, v := range tx.Outputs {
		if v.Type == nil {
			continue
		}
		if contract.IsSameTypeId(v.Type.CodeHash) {
			isCV = true
			break
		}
	}
	return isCV, nil
}

type FuncTransactionHandleReq struct {
	Tx     *types.Transaction
	Hash   string
	Action common.DasAction
}

type FuncTransactionHandleResp struct {
	ActionName string
	Err        error
}

type FuncTransactionHandle func(FuncTransactionHandleReq) FuncTransactionHandleResp
