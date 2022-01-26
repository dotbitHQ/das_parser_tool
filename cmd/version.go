package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Das parser tool",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Das parser tool v1.0.0 -- HEAD")
	},
}
