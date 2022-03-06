package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/AliRostami1/baagh/internal/server"
	"github.com/jedib0t/go-pretty/table"
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

	Run: chipGetCmdRun,
}

type chipGetSuccessResponse struct {
	Success bool              `json:"success"`
	Data    []server.ChipInfo `json:"data"`
}

func chipGetCmdRun(cmd *cobra.Command, args []string) {
	stdOut, stdErr := rootCmd.OutOrStdout(), rootCmd.ErrOrStderr()
	if len(args) == 0 {
		daemonUrl, err := url.Parse(daemonAddr)
		if err != nil {
			fmt.Fprintf(stdErr, "daemon url parsing failed: %v", err)
			return
		}
		daemonUrl.Path = "/api/chips/"
		res, err := http.Get(daemonUrl.String())
		if err != nil {
			fmt.Fprintf(stdErr, "error sending http request to daemon: %v", err)
			return
		}
		defer res.Body.Close()

		decoder := json.NewDecoder(res.Body)

		if res.StatusCode != 200 {
			var errResp server.ErrorResponse
			err = decoder.Decode(&errResp)

			if err != nil {
				fmt.Fprintf(stdErr, "error unmarshaling json: %v", err)
				return
			}
			fmt.Fprintf(stdErr, "daemon responded with the following error: %s", errResp.Message)
			return
		}

		var sucResp chipGetSuccessResponse
		err = decoder.Decode(&sucResp)
		if err != nil {
			fmt.Fprintf(stdErr, "error unmarshaling json: %v", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(stdOut)

		t.AppendHeader(table.Row{"#", "Name", "Label", "Lines", "UapiAbiVersion"})
		for index, row := range sucResp.Data {
			t.AppendRow(table.Row{index, row.Name, row.Label, row.Lines, row.UapiAbiVersion})
		}
		t.Render()

	} else {

	}
}
