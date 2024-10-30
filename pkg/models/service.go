package models

import "fmt"

type Service struct {
	Name    string
	Product string
	Version string
}

func (s Service) String() string {
	if s.Product == "" || s.Name == "unknown" {
		if s.Name == "" {
			// product & name empty
			return "unknown"
		}
		return s.Name
	}

	if s.Version == "" {
		return fmt.Sprintf("%s (%s)", s.Product, s.Name)
	}

	return fmt.Sprintf("%s (%s) v%s", s.Product, s.Name, s.Version)
}
