package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var Config *CfgDict

type CfgDict struct {
	Storage struct {
		File struct {
			Path string `yaml:"path"`
		} `yaml:"file"`
	} `yaml:"storage"`
}

func ReadCfgFromFile(filepath string) error {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("read config file failed: %w", err)
	}

	var cfg CfgDict
	if err := yaml.Unmarshal(content, &cfg); err != nil {
		return fmt.Errorf("parse config file failed: %w", err)
	}

	Config = &cfg
	return nil
}
