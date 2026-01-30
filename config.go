package main

import (
	"context"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/urfave/cli/v3"
)

type Config struct {
	UseOriginalDirectoryLayout bool     `yaml:"useOriginalDirectoryLayout"`
	FilenameTemplateParts      []string `json:"filenameTemplateParts"`
}

func getConfig(ctx context.Context, cmd *cli.Command) (*Config, error) {
	configFile := cmd.String(configFileFlag.Name)

	var config Config
	if configFile != "" {
		f, err := os.Open(configFile)
		if err != nil {
			return nil, err
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(f)

		yamlDecoder := yaml.NewDecoder(f, yaml.Strict())
		err = yamlDecoder.Decode(&config)
		if err != nil {
			return nil, err
		}
	}

	return &config, nil
}
