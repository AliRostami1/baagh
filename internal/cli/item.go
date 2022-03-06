package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(itemCmd)
}

var itemCmd = &cobra.Command{
	Use:     "item",
	Aliases: []string{"i", "it", "ite", "item", "items"},
	Short:   "item short",
	Long:    "item long",
	Example: "baagh-cli item gpiochip0 1",

	Run: itemCmdRun,
}

func itemCmdRun(cmd *cobra.Command, args []string) {
	rootCmd.OutOrStdout()
	fmt.Fprintf(rootCmd.OutOrStdout(), `hello from item`)
}
