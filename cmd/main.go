package main

import (
	"context"
	"das_parser_tool/chain"
	"das_parser_tool/config"
	"das_parser_tool/transaction_parser"
	"flag"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/scorpiotzh/mylog"
	"sync"
)

var (
	log       = mylog.NewLogger("main", mylog.LevelDebug)
	ctxServer = context.Background()
	wgServer  = sync.WaitGroup{}
	dc        *core.DasCore
)

func main() {
	c := flag.String("c", "./config/config.yaml", "config file")
	h := flag.String("t", "", "transaction hash")
	j := flag.String("j", "", "transaction json")

	flag.Parse()

	log.Info(*c, *h, *j)
	log.Info("----- start tx parser -----")
	if *h != "" {
		hashParser(*c, *h)
	}
	if *j != "" {
		jsonParser(*c, *j)
	}
	log.Info("----- end tx parser -----")
}

func hashParser(c, h string) {
	// config
	if err := config.InitCfg(c); err != nil {
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
	bp, err := transaction_parser.NewTransactionParser(transaction_parser.ParamsTransactionParser{
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

func jsonParser(c, j string) {

}
