package analyze

import (
	"github.com/spf13/cobra"
)

var AnalyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyzing commands",
	Long:  `This command group contains commands related to offline analyzing.`,
}

func init() {
	AnalyzeCmd.PersistentFlags().String(
		"openapi", "assets/5GC-APIs", "Path to 3GPP OpenAPI definitions of 5G network functions",
	)

	AnalyzeCmd.AddCommand(pcapCmd)
}
