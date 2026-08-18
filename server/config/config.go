package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var Config *CfgDict

func init() {
	err := ReadCfgFromFile("config.yaml")
	if err != nil {
		panic(fmt.Sprintf("failed to read config file: %v", err))
	}
}

type CfgDict struct {
	Storage struct {
		File struct {
			Path string `yaml:"path"`
		} `yaml:"file"`
	} `yaml:"storage"`
	Cache struct {
		Type     string `yaml:"type"`
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DBIndex  int    `yaml:"db_index"`
	} `yaml:"cache"`
	Model struct {
		Embedding struct {
			APIURL    string `yaml:"api_url"`
			ModelName string `yaml:"model_name"`
			APIKey    string `yaml:"api_key"`
		} `yaml:"embedding"`
		Reranking struct {
			APIURL    string `yaml:"api_url"`
			ModelName string `yaml:"model_name"`
			APIKey    string `yaml:"api_key"`
		} `yaml:"reranking"`
		LLM struct {
			APIURL    string `yaml:"api_url"`
			ModelName string `yaml:"model_name"`
			APIKey    string `yaml:"api_key"`
		} `yaml:"llm"`
	} `yaml:"model"`
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
