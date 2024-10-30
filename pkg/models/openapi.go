package models

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Server struct {
	URL       string               `yaml:"url"`
	Variables map[string]ServerVar `yaml:"variables,omitempty"`
}

type ServerVar struct {
	Default     string `yaml:"default"`
	Description string `yaml:"description,omitempty"`
}

type Parameter struct {
	Name     string      `yaml:"name"`
	In       string      `yaml:"in"`
	Required bool        `yaml:"required"`
	Type     string      `yaml:"type"`
	Example  interface{} `yaml:"example"`
	Default  interface{} `yaml:"default"`
}

type Operation struct {
	Responses  map[string]interface{} `yaml:"responses"`
	Parameters []Parameter            `yaml:"parameters"`
}

type Info struct {
	Title       string `yaml:"title"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
}

type OpenAPI struct {
	Info    Info                            `yaml:"info"`
	Servers []Server                        `yaml:"servers"`
	Paths   map[string]map[string]Operation `yaml:"paths"`
}

func ParseOpenAPIFile(filePath string) (*OpenAPI, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read OpenAPI file: %v", err)
	}

	var openapi OpenAPI
	err = yaml.Unmarshal(data, &openapi)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI file: %v", err)
	}

	return &openapi, nil
}

func ValidateOpenAPIFile(filePath string) (*OpenAPI, error) {
	// Check if the file has a YAML extension
	if !(strings.HasSuffix(filePath, ".yaml") || strings.HasSuffix(filePath, ".yml")) {
		return nil, fmt.Errorf("file %s is not a YAML file", filePath)
	}

	swagger, err := ParseOpenAPIFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("invalid OpenAPI file: %s - %v", filePath, err)
	}

	return swagger, nil
}
