package enum

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/tomsteele/go-nmap"

	"github.com/awareseven/mobilesniper/pkg/models"
	utils "github.com/awareseven/mobilesniper/pkg/utils"
)

func DiscoverOpenPorts(targetNetOrIP string, targetChan chan<- models.Target, wg *sync.WaitGroup, maxConcurrency int, hostTimeout string, verbose bool, nmapArgs ...string) error {
	defer wg.Done()

	ips, err := utils.GetIPsInCIDR(targetNetOrIP)
	if err != nil {
		return fmt.Errorf("%s is not a valid target. Expect network in CIDR notation or IP address", targetNetOrIP)
	}

	semaphore := make(chan struct{}, maxConcurrency)

	for _, ip := range ips {

		if verbose {
			log.Printf("Enumerating Ports on %s\n", ip)
		}

		wg.Add(1)
		semaphore <- struct{}{} // add to channel

		go func() {
			defer func() {
				wg.Done()
				<-semaphore

				if verbose {
					log.Printf("Quit port scan goroutine for %s", ip)
				}
			}()

			// default nmap options to optimize parsing, performance & accurancy
			var nmapCmd = []string{"-Pn", "-oX", "-", "-p-", "--host-timeout", hostTimeout, "--max-retries", "3", "-T5"}

			if len(nmapArgs) != 0 {
				nmapCmd = append(nmapCmd, nmapArgs...)
			}

			// add target ip
			nmapCmd = append(nmapCmd, ip)

			// execute nmap
			cmd := exec.Command("nmap", nmapCmd...)

			// read nmap xml output
			var out bytes.Buffer
			cmd.Stdout = &out

			err := cmd.Run()
			if err != nil {
				log.Printf("Error running nmap for %s: %v", ip, err)
				return
			}

			var nmapRun nmap.NmapRun
			err = xml.Unmarshal(out.Bytes(), &nmapRun)
			if err != nil {
				log.Printf("Error parsing nmap XML output for %s: %v", ip, err)
				return
			}

			for _, host := range nmapRun.Hosts {
				// only hosts with open ports are interesting
				if len(host.Ports) > 0 {

					log.Printf("Finish scan for %s in %v and found %d services.",
						ip,
						time.Time(host.EndTime).Sub(time.Time(host.StartTime)),
						len(host.Ports),
					)

					for _, port := range host.Ports {
						if port.State.State == "open" {

							targetModel := models.Target{
								IP: ip, Port: int(port.PortId), Protocol: port.Protocol,
								Service: models.Service{
									Name: port.Service.Name, Version: port.Service.Version,
									Product: port.Service.Product,
								},
							}

							if verbose {
								log.Printf("Found Port %d on %s", targetModel.Port, targetModel.IP)
							}

							targetChan <- targetModel
						}
					}
				}
			}

			if verbose {
				log.Printf("Finished Ports on %s", ip)
			}
		}()
	}

	return err
}
