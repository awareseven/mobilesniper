package enum

import (
	"github.com/spf13/cobra"
)

var EnumCmd = &cobra.Command{
	Use:   "enum",
	Short: "Enumeration commands",
	Long:  `This command group contains commands related to enumeration.`,
}

func init() {
	nfsCmd.Flags().String(
		"openapi", "assets/5GC-APIs", "Path to 3GPP OpenAPI definitions of 5G network functions",
	)
	nfsCmd.Flags().Float64P(
		"threshold", "t", 70.0, "The threshold of accurancy a NF should be considered as detected.",
	)

	sdcCmd.PersistentFlags().String(
		"tlsCrt", "", "The TLS certificate file to use within the SDC communication.",
	)
	sdcCmd.PersistentFlags().String(
		"tlsKey", "", "The TLS private key file to use within the SDC communication.",
	)
	sdcCmd.PersistentFlags().String(
		"tlsCA", "", "The TLS CA certificate file to use within the SDC communication.",
	)

	EnumCmd.PersistentFlags().String(
		"host-timeout", "4m", "The time before give up a scan on a single host.",
	)

	sdcCmd.AddCommand(providerCmd)
	sdcCmd.AddCommand(consumerCmd)

	EnumCmd.AddCommand(servicesCmd)
	EnumCmd.AddCommand(nfsCmd)
	EnumCmd.AddCommand(sdcCmd)
}
