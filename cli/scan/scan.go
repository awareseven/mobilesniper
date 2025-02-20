package scan

import (
	"github.com/spf13/cobra"
)

var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Perfrom a vulnerability scan.",
	Long:  "This command hroup contains commands to discovery network vulnerabilities.",
}

func init() {
	nessusCmd.PersistentFlags().String(
		"username", "", "Nessus username to authenticate",
	)
	nessusCmd.PersistentFlags().String(
		"password", "", "Nessus password for given user to authenticate.",
	)
	nessusCmd.PersistentFlags().String(
		"url", "https://127.0.0.1:8834", "URL to Nessus Web-UI",
	)

	ScanCmd.AddCommand(nucleiCmd)
	ScanCmd.AddCommand(nessusCmd)
}
