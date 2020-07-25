package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// CreateFileIfNotExist method creates configuration file (in case when it does not
// exist yet) in %HOME%/.dorgconfig location or ./.dorgconfig if %HOME% is not
// defined. It returns path to configuration file and error in case when config
// file cannot be created.
func (c Config) CreateFileIfNotExist() (string, error) {
	configPath := filepath.Join(homeDirPath(), ConfigFileName)
	_, configFileErr := os.Open(configPath)

	if os.IsNotExist(configFileErr) {
		_, newFileErr := os.Create(configPath)
		if newFileErr != nil {
			msg := fmt.Sprintf("Cannot create config file [%s]", configPath)
			return "", errors.Wrap(newFileErr, msg)
		}
	}

	return configPath, nil
}

// Returns path to dorg configuration file.
func homeDirPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}

	return homeDir
}
