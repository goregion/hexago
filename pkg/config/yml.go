package config

import (
	"os"

	yaml "gopkg.in/yaml.v3"
)

func ParseYmlConfig[ConfigType any](data []byte) (*ConfigType, error) {
	var appConfig = new(ConfigType)
	if err := yaml.Unmarshal(data, appConfig); err != nil {
		return nil, err
	}
	return appConfig, nil
}

func ParseYmlFileConfig[ConfigType any](filePath string) (*ConfigType, error) {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return ParseYmlConfig[ConfigType](yamlFile)
}
