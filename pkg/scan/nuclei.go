package scan

import (
	"context"
	"log"
	"sync"

	nuclei "github.com/projectdiscovery/nuclei/v3/lib"
	"github.com/projectdiscovery/nuclei/v3/pkg/output"

	utils "github.com/awareseven/mobilesniper/pkg/utils"
)

func RunNucleiScan(targetNetOrIP string, resultsChan chan<- output.ResultEvent, wg *sync.WaitGroup, verbose bool) {
	defer wg.Done()

	// create nuclei engine with options
	ne, err := nuclei.NewNucleiEngineCtx(
		context.Background(),
		nuclei.WithTemplateFilters(nuclei.TemplateFilters{ProtocolTypes: "http"}),
	)
	defer ne.Close()

	if err != nil && verbose {
		log.Fatal(err)
	}

	// load list of targets
	targets, _ := utils.GetIPsInCIDR(targetNetOrIP)

	// load targets
	ne.LoadTargets(targets, true)
	err = ne.ExecuteWithCallback(func(event *output.ResultEvent) {
		resultsChan <- *event
	})

	if err != nil && verbose {
		log.Fatal(err)
	}
}
