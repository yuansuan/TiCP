package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Resources struct {
	Resources map[string]*Resource `yaml:"resources"`
}

type Resource struct {
	Type      string      `yaml:"type"`
	Resolver  string      `yaml:"resolver"`
	KeySuffix string      `yaml:"key_suffix"`
	Data      []*KeyValue `yaml:"data"`
}

type KeyValue struct {
	Key   string   `yaml:"key"`
	Value []string `yaml:"value"`
}

func GetResources(path string) (map[string]*Resource, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	resource := &Resources{}
	err = yaml.Unmarshal(content, &resource)
	if err != nil {
		return nil, err
	}

	return resource.Resources, nil
}
