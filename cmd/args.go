package cmd

import (
	"context"
	"das_parser_tool/chain"
	"das_parser_tool/config"
	"das_parser_tool/dascore"
	"encoding/json"
	"fmt"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/DeAccountSystems/das-lib/witness"
	"github.com/spf13/cobra"
)

var argsCmd = &cobra.Command{
	Use:   "args",
	Short: "Parser config cell by args",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, v := range args {
			argsParser(v)
		}
	},
}

func argsParser(args string) {
	// ckb node
	ckbClient := chain.NewClient(context.Background(), config.Cfg.Chain.CkbUrl, config.Cfg.Chain.IndexUrl)
	// contract init
	_ = dascore.NewDasCore(ckbClient.Client())
	// config cell parser
	configCell, err := core.GetDasConfigCellInfo(args)
	if err != nil {
		cobra.CheckErr(fmt.Errorf("GetDasConfigCellInfo err: %s", err.Error()))
	}

	res := ckbClient.GetTransactionByHash(configCell.OutPoint.TxHash)
	var witnessByte []byte
	for _, v := range res.Transaction.Witnesses {
		actionDataType := witness.ParserWitnessAction(v)
		if actionDataType == args {
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
