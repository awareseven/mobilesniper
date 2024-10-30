package scan

import (
	"fmt"
	"log"
	"sync"

	"github.com/Cgboal/nessus"

	"github.com/awareseven/mobilesniper/pkg/models"
)

func RunNessusScan(targetNetOrIP string, conf models.NessusConf, wg *sync.WaitGroup, verbose bool) {
	defer wg.Done()

	nessus := nessus.NewNessus(conf.URL)
	nessus.Credentials(conf.Username, conf.Password)

	err := nessus.Authenticate()
	if err != nil {
		log.Fatal("Nessus Authentication failed.")
	}

	scanId, err := nessus.LaunchScan(fmt.Sprintf("[MobileSniper] - %s", targetNetOrIP), targetNetOrIP)
	nessus.Wait(scanId)

	report, _ := nessus.ExportAsNessus(scanId)
	log.Println(report)
}
