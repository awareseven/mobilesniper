package models

import "fmt"

type NetworkFunctionResult struct {
	Target          Target
	NetworkFunction string
	Accuracy        float64
}

func (nfr *NetworkFunctionResult) String() string {
	return fmt.Sprintf("- %v - %s (%f%%)", nfr.Target, nfr.NetworkFunction, nfr.Accuracy)
}

// ContainsNFResult checks if a given NetworkFunctionResult is already present in the list of results.
func ContainsNFResult(results []NetworkFunctionResult, result NetworkFunctionResult) bool {
	for _, r := range results {
		if r.Target.IP == result.Target.IP &&
			r.Target.Port == result.Target.Port &&
			r.NetworkFunction == result.NetworkFunction {
			return true
		}
	}
	return false
}
