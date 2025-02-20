package scan

import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"

	"github.com/awareseven/mobilesniper/cli/core"
	"github.com/awareseven/mobilesniper/pkg/models"
	"github.com/awareseven/mobilesniper/pkg/scan"
	utils "github.com/awareseven/mobilesniper/pkg/utils"
)

var nessusCmd = &cobra.Command{
	Use:   "nessus <network range or single IP>",
	Short: "Run a nessus scan on the given target.",
	Long:  "This command runs nessus to perform a vulnerability scan.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		cidrOrIP := args[0]
		_, err := utils.GetIPsInCIDR(cidrOrIP)
		if err != nil {
			panic(err)
		}

		nessus_username, errName := cmd.Flags().GetString("username")
		nessus_password, errPwd := cmd.Flags().GetString("password")
		nessus_url, errUrl := cmd.Flags().GetString("url")

		if errName != nil || errPwd != nil || errUrl != nil {
			panic(err)
		}

		config := models.NessusConf{
			Username: nessus_username,
			Password: nessus_password,
			URL:      nessus_url,
		}

		bar, _ := core.NewProgressBar(1, fmt.Sprintf("Scanning: %s", cidrOrIP))
		defer bar.Finish()

		var wg sync.WaitGroup

		wg.Add(1)
		go scan.RunNessusScan(cidrOrIP, config, &wg, core.Verbose)
		wg.Wait()
	},
}
