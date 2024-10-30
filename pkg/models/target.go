package models

import (
	"fmt"
	"strings"
)

type Target struct {
	IP       string
	Port     int
	Service  Service
	Protocol string
}

func (t Target) String() string {
	if t.Service.String() == "unknown" {
		return fmt.Sprintf("%s:%d/%s", t.IP, t.Port, t.Protocol)
	}
	return fmt.Sprintf("%s:%d/%s: %s", t.IP, t.Port, t.Protocol, t.Service)
}

func PrintTargets(t *[]Target) string {
	targets := make([]string, len(*t))
	for i, target := range *t {
		targets[i] = fmt.Sprint(target)
	}

	fmt.Printf("%s", strings.Join(targets, "\n"))
	return strings.Join(targets, "\n")
}
