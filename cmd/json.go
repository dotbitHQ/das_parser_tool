package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
)

var (
	jsonCmd = &cobra.Command{
		Use:   "json",
		Short: "Parser transaction by transaction json",
	}
	jsonFileCmd = &cobra.Command{
		Use:   "file",
		Short: "Parser transaction by transaction json file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			file, err := ioutil.ReadFile(args[0])
			if err != nil {
				cobra.CheckErr(fmt.Errorf("ReadFile err: %v ", err.Error()))
			}

			jsonParser(string(file))
		},
	}
	jsonDataCmd = &cobra.Command{
		Use:   "data",
		Short: "Parser transaction by transaction json data",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for _, v := range args {
				jsonParser(v)
			}
		},
	}
)

func init() {
	jsonCmd.AddCommand(jsonFileCmd)
	jsonCmd.AddCommand(jsonDataCmd)
}

func jsonParser(arg string) {
	// TODO parser by transaction json
	fmt.Println(arg)
}
