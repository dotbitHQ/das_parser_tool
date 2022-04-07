package dascore

import (
	"context"
	"das_parser_tool/config"
	"fmt"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/scorpiotzh/mylog"
	"github.com/spf13/cobra"
	"sync"
)

var (
	ctxServer = context.Background()
	wgServer  = sync.WaitGroup{}
)

func NewDasCore(client rpc.Client) *core.DasCore {
	core.SetLogLevel(mylog.LevelError)
	config.Env = core.InitEnv(config.Cfg.Chain.Net)
	opts := []core.DasCoreOption{
		core.WithClient(client),
		core.WithDasContractArgs(config.Env.ContractArgs),
		core.WithDasContractCodeHash(config.Env.ContractCodeHash),
		core.WithDasNetType(config.Cfg.Chain.Net),
		core.WithTHQCodeHash(config.Env.THQCodeHash),
	}
	dc := core.NewDasCore(ctxServer, &wgServer, opts...)
	dc.InitDasContract(config.Env.MapContract)
	if err := dc.InitDasConfigCell(); err != nil {
		cobra.CheckErr(fmt.Errorf("InitDasConfigCell err: %v", err))
	}
	if err := dc.InitDasSoScript(); err != nil {
		cobra.CheckErr(fmt.Errorf("InitDasSoScript err: %v", err))
	}

	return dc
}
