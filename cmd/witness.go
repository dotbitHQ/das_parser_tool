package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/witness"
	"github.com/spf13/cobra"
)

var witnessCmd = &cobra.Command{
	Use:   "witness",
	Short: "Parser transaction by witness",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, v := range args {
			witnessParser(v)
		}
	},
}

func witnessParser(arg string) {
	witnessByte := common.Hex2Bytes(arg)
	b, err := json.Marshal(witness.ParserWitnessData(witnessByte))
	if err != nil {
		cobra.CheckErr(fmt.Errorf("Marshal err: %v ", err.Error()))
	}
	fmt.Println(string(b))
}
