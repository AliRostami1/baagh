package cli

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(chipCmd)
}

var chipCmd = &cobra.Command{
	Use:     "chip",
	Aliases: []string{"c", "ch", "chi", "chip", "chips"},
	Short:   "chip short",
	Long:    "chip long",
	Example: "baagh-cli chip",

	Run: chipCmdRun,
}

func chipCmdRun(cmd *cobra.Command, args []string) {
	rootCmd.OutOrStdout()
	fmt.Fprintf(rootCmd.OutOrStdout(), `hello from chip`)
}

func init() {
	chipCmd.AddCommand(chipGetCmd)
}

var chipGetCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"g", "ge", "get"},
	Short:   "chip short",
	Long:    "chip long",
	Example: "baagh-cli chip get",
	Args:    cobra.MaximumNArgs(1),

	RunE: chipGetCmdRun,
}

func chipGetCmdRun(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		dUrl, err := url.Parse(daemonAddr)
		if err != nil {
			return fmt.Errorf("daemon url parsing failed: %v", err)
		}
		dUrl.Path = "/api/chips"
		res, err := http.Get(dUrl.String())
		if err != nil {
			return fmt.Errorf("error sending http request to daemon: %v", err)
		}

		//We Read the response body on the line below.
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error reading the response body: %v", err)
		}

		fmt.Fprintf(rootCmd.OutOrStdout(), "response: %s", string(body))
	} else {

	}
	return nil
}
