package cmd

import (
	"context"
	"das_parser_tool/chain"
	"das_parser_tool/config"
	"das_parser_tool/dascore"
	"das_parser_tool/parser"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
)

var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "Parser transaction by hash",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, v := range args {
			hashParser(v)
		}
	},
}

func hashParser(arg string) {
	// ckb node
	ckbClient := chain.NewClient(context.Background(), config.Cfg.Chain.CkbUrl, config.Cfg.Chain.IndexUrl)
	// contract init
	dasCore := dascore.NewDasCore(ckbClient.Client())
	// transaction parser
	bp := parser.NewParser(parser.ParamsParser{
		DasCore:   dasCore,
		CkbClient: ckbClient,
	})
	out := bp.HashParser(arg)

	b, err := json.Marshal(out)
	if err != nil {
		cobra.CheckErr(fmt.Errorf("Marshal err: %v ", err.Error()))
	}
	fmt.Println(string(b))
}
