package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const ConfigFileName = ".dorgconfig"

// Config stores all of dorg configuration.
type Config struct {
	Filepath       string
	DownloadsPath  string `json:"downloads_path"`
	TargetBasePath string `json:"target_path"`
}

// CreateFile creates or alter configuration file and returns it as io.ReadWriteCloser.
func (c *Config) CreateFile() (io.ReadWriteCloser, error) {
	err := c.createDorgConfigFile()
	if err != nil {
		return nil, err
	}

	configFile, err := os.Create(c.Filepath)
	if err != nil {
		msg := fmt.Sprintf("Cannot create config file [%s]", c.Filepath)
		return nil, errors.Wrap(err, msg)
	}
	return configFile, nil
}

// OpenFile opens configuration file as io.ReadWriteCloser.
func (c *Config) OpenFile() (io.ReadWriteCloser, error) {
	configFile, err := os.Open(c.Filepath)
	if err != nil {
		msg := fmt.Sprintf("Cannot open config file [%s]", c.Filepath)
		return nil, errors.Wrap(err, msg)
	}
	return configFile, nil
}

// SaveToFile saves to file (created using CreateFileIfNotExist) configuration from Config.
func (c *Config) SaveToFile(configFile io.WriteCloser) error {
	err := c.createDorgConfigFile()
	if err != nil {
		return err
	}

	configJson, jsonErr := json.MarshalIndent(*c, "", "\t")
	if jsonErr != nil {
		msg := fmt.Sprintf("Cannot marshal [%v] into JSON", *c)
		return errors.Wrap(jsonErr, msg)
	}

	_, wErr := configFile.Write(configJson)
	if wErr != nil {
		msg := fmt.Sprintf("Error while writing into config file [%s]", c.Filepath)
		return errors.Wrap(wErr, msg)
	}

	if closeErr := configFile.Close(); closeErr != nil {
		msg := fmt.Sprintf("Error while closing config file [%s]", c.Filepath)
		return errors.Wrap(closeErr, msg)
	}

	return nil
}

// TryParseConfig tries to load given file as Config object.
func TryParseConfig(configFile io.ReadCloser) (Config, error) {
	configBytes, err := readBytesFromFile(configFile)
	if err != nil {
		return Config{}, err
	}

	var config Config
	jsonErr := json.Unmarshal(configBytes, &config)
	if jsonErr != nil {
		msg := fmt.Sprintf("Cannot unmarshal this file into JSON: \n [%s]",
			string(configBytes))
		return Config{}, errors.Wrap(jsonErr, msg)
	}

	if closeErr := configFile.Close(); closeErr != nil {
		msg := fmt.Sprintf("Error while closing config file")
		return Config{}, errors.Wrap(closeErr, msg)
	}

	return config, nil
}

// Checks if configuration files are indentical.
func (c Config) Equals(another Config) bool {
	pathEquals := c.Filepath == another.Filepath
	dpEquals := c.DownloadsPath == another.DownloadsPath
	tpEquals := c.TargetBasePath == another.TargetBasePath

	return pathEquals && dpEquals && tpEquals
}

// createDorgConfigFile method creates configuration file (in case when it does not
// exist yet) in %HOME%/.dorgconfig location or ./.dorgconfig if %HOME% is not
// defined. Config Filepath is updated after creating new config file.
func (c *Config) createDorgConfigFile() error {
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

// Read reads configuration files as bytes.
func readBytesFromFile(reader io.Reader) ([]byte, error) {
	configBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return configBytes, nil
}
