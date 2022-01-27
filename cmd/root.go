package cmd

import (
	"das_parser_tool/config"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"os"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "tx_parser",
		Short: "A tool for das parser Transaction",
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "parser config file (default is ./config/config.yaml)")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(hashCmd)
	rootCmd.AddCommand(witnessCmd)
	rootCmd.AddCommand(jsonCmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func initConfig() {
	config.InitCfg(cfgFile)
}

func Execute() {
	// Gen markdown files
	err := doc.GenMarkdownTree(rootCmd, "./doc")
	if err != nil {
		cobra.CheckErr(fmt.Errorf("GenMarkdown err: %v ", err.Error()))
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
