package config

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func (c Config) FileExist() bool {

	configPath := filepath.Join(homeDir, ConfigFileName)
	_, configFileErr := os.Open(configPath)

	if configFileErr != nil {
		return false, 
	}

	return true, nil
}

// Returns path to dorg configuration file.
func homeDirPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ConfigFileName)
	}

	return filepath.Join(homeDir, ConfigFileName)
}
