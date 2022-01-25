package main

import (
	"context"
	"das_parser_tool/chain"
	"das_parser_tool/config"
	"das_parser_tool/dascore"
	"das_parser_tool/parser"
	"encoding/json"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/witness"
	"github.com/spf13/cobra"
)

var (
	cfgFile  string
	jsonFile string

	rootCmd = &cobra.Command{
		Use:   "das_parser_tool",
		Short: "A tool for das parser Transaction",
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Das parser tool",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Das parser tool v1.0.0 -- HEAD")
		},
	}
	hashCmd = &cobra.Command{
		Use:   "hash",
		Short: "Parser transaction by hash",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for _, v := range args {
				hashParser(v)
			}
		},
	}
	witnessCmd = &cobra.Command{
		Use:   "witness",
		Short: "Parser transaction by witness",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for _, v := range args {
				witnessParser(v)
			}
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "parser config file (default is ./config/config.yaml)")
	// TODO parser by transaction json and transaction json file
	rootCmd.PersistentFlags().StringVar(&jsonFile, "json", "", "Parser transaction by transaction json")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(hashCmd)
	rootCmd.AddCommand(witnessCmd)
}

func initConfig() {
	config.InitCfg(cfgFile)
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

func witnessParser(arg string) {
	witnessByte := common.Hex2Bytes(arg)
	b, err := json.Marshal(witness.ParserWitnessData(witnessByte))
	if err != nil {
		cobra.CheckErr(fmt.Errorf("Marshal err: %v ", err.Error()))
	}
	fmt.Println(string(b))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		cobra.CheckErr(err)
	}
}
