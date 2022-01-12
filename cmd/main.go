package main

import (
	"context"
	"das_parser_tool/chain/chain_ckb"
	"das_parser_tool/config"
	"das_parser_tool/transaction_parser"
	"flag"
	"fmt"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/scorpiotzh/mylog"
	"sync"
	"time"
)

var (
	log       = mylog.NewLogger("main", mylog.LevelDebug)
	ctxServer = context.Background()
	wgServer  = sync.WaitGroup{}
	dc        *core.DasCore
)

func main() {
	c := flag.String("c", "./config/config.yaml", "config file")
	t := flag.String("t", "", "transaction hash")

	flag.Parse()
	if *t == "" {
		log.Fatal("transaction hash is empty")
	}

	fmt.Println(*c, *t)
	fmt.Println("----- start tx parser -----")
	txParser(*c, *t)
	fmt.Println("----- end tx parser -----")
}

func txParser(c, t string) {
	// config
	if err := config.InitCfg(c); err != nil {
		log.Fatal(err)
	}

	// ckb node
	ckbClient, err := chain_ckb.NewClient(context.Background(), config.Cfg.Chain.CkbUrl, config.Cfg.Chain.IndexUrl)
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
	dc.RunAsyncDasContract(time.Minute * 5)   // contract outpoint
	dc.RunAsyncDasConfigCell(time.Minute * 3) // config cell outpoint
	dc.RunAsyncDasSoScript(time.Minute * 7)   // so
	log.Info("contract ok")

	// transaction parser
	bp, err := transaction_parser.NewTransactionParser(transaction_parser.ParamsTransactionParser{
		DasCore:            dc,
		CkbClient:          ckbClient,
		Ctx:                ctxServer,
		Wg:                 &wgServer,
	})
	if err != nil {
		log.Fatal(err)
	}
	bp.RunParser(t)
}

