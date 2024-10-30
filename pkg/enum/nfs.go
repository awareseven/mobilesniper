package enum

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/awareseven/mobilesniper/pkg/models"
)

func DiscoverNetworkFunctions(target models.Target, openapiPath string, nfrChan chan<- models.NetworkFunctionResult, wg *sync.WaitGroup, maxConcurrency int, verbose bool) error {
	defer wg.Done()

	semaphore := make(chan struct{}, maxConcurrency)

	if verbose {
		log.Printf("Enumerating NF definitions on %s:%d", target.IP, target.Port)
	}

	err := filepath.Walk(openapiPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		wg.Add(1)
		semaphore <- struct{}{} // add to channel

		go func() {
			defer func() {
				wg.Done()
				<-semaphore // remove from channel

				if verbose {
					log.Printf("Quit %s goroutine for %s:%d", path, target.IP, target.Port)
				}
			}()

			openapi, err := models.ValidateOpenAPIFile(path)
			if err != nil {
				if verbose {
					log.Printf("Invalid OpenAPI file: %s - %v\n", path, err)
				}
				return // continue processing other files
			}

			var reachableCount int
			var totalExpectedResponses int
			client := &http.Client{
				Timeout: 30 * time.Second,
			}

			// Iterate over all paths defined in the OpenAPI specification
			for path, methods := range openapi.Paths {
				// Iterate over all HTTP methods (GET, POST, etc.) for each path
				for method, operation := range methods {

					// Replace path variables with example or default values
					for _, param := range operation.Parameters {
						if param.In == "path" {
							placeholder := fmt.Sprintf("{%s}", param.Name)
							var value string
							if param.Example != nil {
								value = fmt.Sprintf("%v", param.Example)
							} else if param.Default != nil {
								value = fmt.Sprintf("%v", param.Default)
							} else {
								// Provide a generic placeholder if no example or default is provided
								value = placeholder
							}
							path = strings.Replace(path, placeholder, value, -1)
						}
					}

					// Construct the request URL for the given path and method
					url := fmt.Sprintf("http://%s:%d%s", target.IP, target.Port, path)
					req, err := http.NewRequest(strings.ToUpper(method), url, nil)
					if err != nil {
						if verbose {
							log.Printf("Error creating request for %s %s: %v", method, path, err)
						}
						continue // Skip to the next method if there's an error creating the request
					}

					// Execute the request
					resp, err := client.Do(req)
					if err != nil {
						if verbose {
							log.Printf("Error executing request for %s %s: %v", method, path, err)
						}
						continue // Skip to the next method if there's an error executing the request
					}

					// Check if the returned status code is one of the expected codes
					if _, ok := operation.Responses[fmt.Sprintf("%d", resp.StatusCode)]; ok {
						// Increment the reachable count if the response matches one of the expected status codes
						reachableCount++
						if verbose {
							log.Printf("Matched expected status code %d for %s %s", resp.StatusCode, method, path)
						}
					} else {
						if verbose {
							log.Printf("Unexpected status code %d for %s %s", resp.StatusCode, method, path)
						}
					}

					// Explicitly close the response body here
					resp.Body.Close()

					// Increment the total expected responses count
					totalExpectedResponses += 1
				}
			}

			// Calculate accuracy based on the number of expected responses
			accuracy := (float64(reachableCount) / float64(totalExpectedResponses)) * 100

			apiName := openapi.Info.Title

			result := models.NetworkFunctionResult{
				Target:          target,
				NetworkFunction: apiName,
				Accuracy:        accuracy,
			}

			if verbose {
				log.Printf("%f%% NF %s on %s:%d",
					result.Accuracy, result.NetworkFunction,
					result.Target.IP, result.Target.Port,
				)
			}

			nfrChan <- result
		}()

		return nil // Return nil to continue processing other files
	})

	return err
}
