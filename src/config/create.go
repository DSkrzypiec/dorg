package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// CreateFileIfNotExist method creates configuration file (in case when it does not
// exist yet) in %HOME%/.dorgconfig location or ./.dorgconfig if %HOME% is not
// defined. Config Filepath is updated after creating new config file.
func (c *Config) CreateFileIfNotExist() error {
	configPath := filepath.Join(homeDirPath(), ConfigFileName)
	file, configFileErr := os.Open(configPath)
	defer file.Close()

	if os.IsNotExist(configFileErr) {
		_, newFileErr := os.Create(configPath)
		if newFileErr != nil {
			msg := fmt.Sprintf("Cannot create config file [%s]", configPath)
			return errors.Wrap(newFileErr, msg)
		}
	}

	c.Filepath = configPath
	return nil
}

// Returns path to dorg configuration file.
func homeDirPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}

	return homeDir
}
