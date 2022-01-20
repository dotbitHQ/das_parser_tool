package parser

import (
	"context"
	"das_parser_tool/chain"
	"das_parser_tool/config"
	"github.com/DeAccountSystems/das-lib/core"
	"sync"
	"testing"
)

var (
	ctxServer = context.Background()
	wgServer  = sync.WaitGroup{}
	dc        *core.DasCore
)

func parserHash(h string) {
	if err := config.InitCfg("../config/config.yaml"); err != nil {
		log.Fatal(err)
	}

	// ckb node
	ckbClient, err := chain.NewClient(context.Background(), config.Cfg.Chain.CkbUrl, config.Cfg.Chain.IndexUrl)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("node ok")

	// das contract init
	opts := []core.DasCoreOption{
		core.WithClient(ckbClient.Client()),
		core.WithDasContractArgs(config.Cfg.DasCore.DasContractArgs),
		core.WithDasContractCodeHash(config.Cfg.DasCore.DasContractCodeHash),
		core.WithDasNetType(config.Cfg.Chain.Net),
		core.WithTHQCodeHash(config.Cfg.DasCore.THQCodeHash),
	}
	dc = core.NewDasCore(ctxServer, &wgServer, opts...)
	dc.InitDasContract(config.Cfg.DasCore.MapDasContract)
	if err = dc.InitDasConfigCell(); err != nil {
		log.Fatal(err)
	}
	if err = dc.InitDasSoScript(); err != nil {
		log.Fatal(err)
	}
	log.Info("contract ok")

	// transaction parser
	bp, err := NewParser(ParamsParser{
		DasCore:   dc,
		CkbClient: ckbClient,
		Ctx:       ctxServer,
		Wg:        &wgServer,
	})
	if err != nil {
		log.Fatal(err)
	}
	bp.RunParser(h)
}

func TestTransferBalance(t *testing.T) {
	parserHash("0x072a302f9563e557bad9969514df671fe9d0f7253c6e471da3f3d99bb56779f6")
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
