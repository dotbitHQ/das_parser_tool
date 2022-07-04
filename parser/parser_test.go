package parser

import (
	"context"
	"das_parser_tool/chain"
	"das_parser_tool/config"
	"das_parser_tool/dascore"
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/spf13/cobra"
	"strings"
	"testing"
)

func parserHash(h string) {
	// config
	config.InitCfg("")

	// ckb node
	ckbClient := chain.NewClient(context.Background(), config.Cfg.Chain.CkbUrl, config.Cfg.Chain.IndexUrl)

	// das contract init
	dc := dascore.NewDasCore(ckbClient.Client())

	// transaction parser
	bp := Parser{
		CkbClient: ckbClient,
		DasCore:   dc,
	}
	out := bp.HashParser(h)

	b, err := json.Marshal(out)
	if err != nil {
		cobra.CheckErr(fmt.Errorf("Marshal err: %v ", err.Error()))
	}
	fmt.Println(string(b))
}

func TestTransferBalance(t *testing.T) {
	parserHash("0xeb4a0ac7582277bef6df42e06a7d0efd6d24a11cb362b9dc2fee1da0fafa8228")
}

func TestWithdrawFromWallet(t *testing.T) {
	parserHash("0xeb65b00ed081682219fa0a3d45f566441ab19db0372533b531de9bc3c596e210")
}

func TestTransfer(t *testing.T) {
	parserHash("0x8ef4ce4a33e97319b303e028b3bc3e9ce6f34fbb4b63ccf6fbe0844a90ba9fb2")
}

func TestCreateIncome(t *testing.T) {
	parserHash("0x4e79fc8a02249555bb57546c8374a9f747bab1f13df6ce8aa610d72d4214a5c8")
}

func TestConsolidateIncome(t *testing.T) {
	parserHash("0x74890a9a3800c583ee7bf704d5fe84a8ffca0a3d7877a488493a333e7c1af08f")
}

func TestApplyRegister(t *testing.T) {
	parserHash("0xef817a24ee7d1ae82e78c20a38d18fb9a53803d50d875168b47cdf20d5b53392")
}

func TestPreRegister(t *testing.T) {
	parserHash("0xd7cb4234ae72340ef7f42c0d5fbc8eb422ba25741ddfe24d72078ce2be08a020")
}

func TestPropose(t *testing.T) {
	parserHash("0xf45e2fb69a56441701a1549c853c1c03324544147da0e0023e98341077b3ff57")
}

func TestConfirmProposal(t *testing.T) {
	parserHash("0xe03ad8605220dba3a65f778ef86c9c07207410abc28462bee16b7977df1f174a")
}

func TestEditRecords(t *testing.T) {
	parserHash("0xf077361d3078a6cab299e3d813c9fe5cb9916d92aa3052e45d2e79bda89db7ab")
}

func TestEditManager(t *testing.T) {
	parserHash("0x023fed77f3a72c34dc9642aa7d542128a2668020193923852009be364981171d")
}

func TestRenewAccount(t *testing.T) {
	parserHash("0xf633b25d3272cbb8aea0af58ad0e5430183e682bc3dea83f18c0bb38df908296")
}

func TestTransferAccount(t *testing.T) {
	parserHash("0x477fdf091f553d1811599a3292198ee547a5b3cfb9d51f125c163c1c27dbc932")
}

func TestStartAccountSale(t *testing.T) {
	parserHash("0x0b762bcd7679365be06696f7a8ff59472bc28b1294ee55374e840ee500f72436")
}

func TestEditAccountSale(t *testing.T) {
	parserHash("0x208cf033969fec9ba2a8b889bf884804a33f7655db7e8c4f223c533614cdd33c")
}

func TestCancelAccountSale(t *testing.T) {
	parserHash("0x89c60ffe04afea217f4cdf524805b4b98e0a42608ac7a4bcc4d0c4a0e4986382")
}

func TestBuyAccount(t *testing.T) {
	parserHash("0xe388608052e2d4008336d6e6ab3f7ca457397df84fa043199191a3a7350f5b0d")
}

func TestMakeOffer(t *testing.T) {
	parserHash("0x28d70b6fcc59290c3a73fcf5fc0e006b80cb461ac17e74223d75cf81d32706d1")
}

func TestEditOffer(t *testing.T) {
	parserHash("0xed3e4a0c665ba970013ffde09cf1aecdab9dd03af103ffd88dfb35d1ddc6cfec")
}

func TestCancelOffer(t *testing.T) {
	parserHash("0xfcef771cc7c7199d8ffac419b8b77fcd5581ff29a0cf1785ca59f38eade75587")
}

func TestAcceptOffer(t *testing.T) {
	parserHash("0xff14b214049c9ff61660e58d11e1e06402aa2f45291218f7f0e2ca49c1c67684")
}

func TestDeclareReverseRecord(t *testing.T) {
	parserHash("0x141bd83f467c631dde1384572871a24a2aff247dec8d07f0fce0a5a5d15180a7")
}

func TestRedeclareReverseRecord(t *testing.T) {
	parserHash("0x1fc8520f73c2c5383ecd0c824a03ea7a61c4b71579704535bab3bda948cca296")
}

func TestRetractReverseRecord(t *testing.T) {
	parserHash("0x39ff46cea5dfaeca4dc591d69913a9484507db30030de2e62024f7662214ee51")
}

func TestConfigCellRelease(t *testing.T) {
	parserHash("0x53d077fe2f29027f29985a54c2514f1978b5a37167113f9908289cbf3d2761ac")
}

func parserArgs(a string) {
	// config
	config.InitCfg("")

	// ckb node
	ckbClient := chain.NewClient(context.Background(), config.Cfg.Chain.CkbUrl, config.Cfg.Chain.IndexUrl)

	// das contract init
	_ = dascore.NewDasCore(ckbClient.Client())

	// config cell parser
	configCell, err := core.GetDasConfigCellInfo(a)
	if err != nil {
		cobra.CheckErr(fmt.Errorf("GetDasConfigCellInfo err: %s", err.Error()))
	}

	res := ckbClient.GetTransactionByHash(configCell.OutPoint.TxHash)
	var witnessByte []byte
	fmt.Println(strings.Repeat("-", 100))
	for _, v := range res.Transaction.Witnesses {
		actionDataType := witness.ParserWitnessAction(v)
		if actionDataType == a {
			witnessByte = v
			break
		}
	}

	out := witness.ParserWitnessData(witnessByte)
	b, err := json.Marshal(out)
	if err != nil {
		cobra.CheckErr(fmt.Errorf("Marshal err: %v ", err.Error()))
	}
	fmt.Println(string(b))
}

func TestParserArgs(t *testing.T) {
	parserArgs(common.ConfigCellTypeArgsAccount)
	parserArgs(common.ConfigCellTypeArgsMain)
	parserArgs(common.ConfigCellTypeArgsSubAccount)

	parserArgs(common.ConfigCellTypeArgsPreservedAccount00)
	parserArgs(common.ConfigCellTypeArgsPreservedAccount19)

	parserArgs(common.ConfigCellTypeArgsCharSetEmoji)
	parserArgs(common.ConfigCellTypeArgsCharSetHanT)
}

func TestEnableSubAccount(t *testing.T) {
	parserHash("0xd72ab7cd90fc30337663a60f1809888dc0277f7409b28721cb33079225fada53")
}

func TestCreateSubAccount(t *testing.T) {
	parserHash("0xf8d4a5466e7c9942699664fc133ff63ff79dbcc5819c2c81cd14db83102ca23b")
}

func TestEditSubAccount(t *testing.T) {
	parserHash("0x95730a147d07031ce8240cb8ff760172f280604598bdfcb05c09fc1840de39bc")
}

func TestLockAccountForCrossChain(t *testing.T) {
	parserHash("0x70064f9f4d9bd3850dc34d99abaca1b659e7055939eb53c9c5ed51ece74b25e0")
}

func TestUnlockAccountForCrossChain(t *testing.T) {
	parserHash("0x42d8f70429b46a742c3511f9914cd698e3c963ba456a9a7a2bea56dd76f9a818")
}

func TestConfigSubAccountCreatingScript(t *testing.T) {
	parserHash("0xc5e95701d4f2e20a0f023b138762c9f4faf59b1b43673297c9a116338b836ef4")
}

func TestRecycleExpiredAccount(t *testing.T) {
	parserHash("0xdaa0e66a9ff150986e0c433ca315b37c20e40926c73aeecf3eb2c9a434f781f3")
}

func TestCompatibleWitnessParser(t *testing.T) {
	witnessByte := common.Hex2Bytes("0x646173010000006e0100001000000010000000100000005e010000100000001400000018000000000000000500000042010000420100003400000048000000c9000000d1000000d9000000e1000000e9000000ea0000002f01000030010000380100003d010000a35ea5d5ef43a74e95351254802c334237bdde2b81000000180000002d00000042000000570000006c000000150000000c00000010000000020000000100000074150000000c0000001000000002000000010000007a150000000c00000010000000020000000100000068150000000c00000010000000010000000100000030150000000c0000001000000001000000010000003180cbde6000000000000000000000000000000000000000005474b562000000000045000000080000003d00000018000000230000002e00000032000000390000000700000070726f66696c65070000007477697474657200000000030000003131312c01000000000000000000000001000000de01000000de")
	b, err := json.Marshal(witness.ParserWitnessData(witnessByte))
	if err != nil {
		cobra.CheckErr(fmt.Errorf("Marshal err: %v ", err.Error()))
	}
	fmt.Println(string(b))
}
