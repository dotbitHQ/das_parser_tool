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

	if resp.ActionName != "" {
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
				resp.WitnessesMap = append(resp.WitnessesMap, map[string]interface{}{
					"name":         "ConfigCellMain",
					"witness":      common.Bytes2Hex(v),
					"witness_hash": common.Bytes2Hex(common.Blake2b(builder.ConfigCellMain.AsSlice())),
					"status":       common.Bytes2Hex(builder.ConfigCellMain.Status().RawData()),
					"type_id_table": map[string]interface{}{
						"account_cell":         common.Bytes2Hex(builder.ConfigCellMain.TypeIdTable().AccountCell().RawData()),
						"apply_register_cell":  common.Bytes2Hex(builder.ConfigCellMain.TypeIdTable().ApplyRegisterCell().RawData()),
						"balance_cell":         common.Bytes2Hex(builder.ConfigCellMain.TypeIdTable().BalanceCell().RawData()),
						"income_cell":          common.Bytes2Hex(builder.ConfigCellMain.TypeIdTable().IncomeCell().RawData()),
						"pre_account_cell":     common.Bytes2Hex(builder.ConfigCellMain.TypeIdTable().PreAccountCell().RawData()),
						"proposal_cell":        common.Bytes2Hex(builder.ConfigCellMain.TypeIdTable().ProposalCell().RawData()),
						"account_sale_cell":    common.Bytes2Hex(builder.ConfigCellMain.TypeIdTable().AccountSaleCell().RawData()),
						"account_auction_cell": common.Bytes2Hex(builder.ConfigCellMain.TypeIdTable().AccountAuctionCell().RawData()),
						"offer_cell":           common.Bytes2Hex(builder.ConfigCellMain.TypeIdTable().OfferCell().RawData()),
						"reverse_record_cell":  common.Bytes2Hex(builder.ConfigCellMain.TypeIdTable().ReverseRecordCell().RawData()),
					},
					"das_lock_out_point_table": map[string]interface{}{
						"ckb_signall": map[string]interface{}{
							"tx_hash": common.Bytes2Hex(builder.ConfigCellMain.DasLockOutPointTable().CkbSignall().TxHash().RawData()),
							"index":   common.Bytes2Hex(builder.ConfigCellMain.DasLockOutPointTable().CkbSignall().Index().RawData()),
						},
						"ckb_multisign": map[string]interface{}{
							"tx_hash": common.Bytes2Hex(builder.ConfigCellMain.DasLockOutPointTable().CkbMultisign().TxHash().RawData()),
							"index":   common.Bytes2Hex(builder.ConfigCellMain.DasLockOutPointTable().CkbMultisign().Index().RawData()),
						},
						"ckb_anyone_can_pay": map[string]interface{}{
							"tx_hash": common.Bytes2Hex(builder.ConfigCellMain.DasLockOutPointTable().CkbAnyoneCanPay().TxHash().RawData()),
							"index":   common.Bytes2Hex(builder.ConfigCellMain.DasLockOutPointTable().CkbAnyoneCanPay().Index().RawData()),
						},
						"eth": map[string]interface{}{
							"tx_hash": common.Bytes2Hex(builder.ConfigCellMain.DasLockOutPointTable().Eth().TxHash().RawData()),
							"index":   common.Bytes2Hex(builder.ConfigCellMain.DasLockOutPointTable().Eth().Index().RawData()),
						},
						"tron": map[string]interface{}{
							"tx_hash": common.Bytes2Hex(builder.ConfigCellMain.DasLockOutPointTable().Tron().TxHash().RawData()),
							"index":   common.Bytes2Hex(builder.ConfigCellMain.DasLockOutPointTable().Tron().Index().RawData()),
						},
					},
				})
				continue
			}
		}
	}

	return
}
