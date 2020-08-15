package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

const ConfigFileName = ".dorgconfig"

// Config stores all of dorg configuration.
type Config struct {
	Filepath       string
	DownloadsPath  string `json:"downloads_path"`
	TargetBasePath string `json:"target_path"`
}

// Default creates default dorg configuration.
func Default(downloadsPath string) Config {
	return Config{
		DownloadsPath:  downloadsPath,
		TargetBasePath: downloadsPath,
	}
}

// SaveToFile saves to file (created using CreateFileIfNotExist) configuration from Config.
func (c *Config) SaveToFile() error {
	err := c.CreateFileIfNotExist()
	if err != nil {
		return err
	}

	configJson, jsonErr := json.MarshalIndent(*c, "", "\t")
	if jsonErr != nil {
		msg := fmt.Sprintf("Cannot marshal [%v] into JSON", *c)
		return errors.Wrap(jsonErr, msg)
	}

	file, fErr := os.Create(c.Filepath)
	if fErr != nil {
		msg := fmt.Sprintf("Error while opening config file [%s]", c.Filepath)
		return errors.Wrap(fErr, msg)
	}

	_, wErr := file.Write(configJson)
	if wErr != nil {
		msg := fmt.Sprintf("Error while writing into config file [%s]", c.Filepath)
		return errors.Wrap(wErr, msg)
	}

	if closeErr := file.Close(); closeErr != nil {
		msg := fmt.Sprintf("Error while closing config file [%s]", c.Filepath)
		return errors.Wrap(closeErr, msg)
	}

	return nil
}

// TryParseConfig tries to load given file as Config object.
func TryParseConfig(configPath string) (Config, error) {
	configBytes, err := readBytesFromFile(configPath)
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

	config.Filepath = configPath

	return config, nil
}

// Checks if configuration files are indentical.
func (c Config) Equals(another Config) bool {
	pathEquals := c.Filepath == another.Filepath
	dpEquals := c.DownloadsPath == another.DownloadsPath
	tpEquals := c.TargetBasePath == another.TargetBasePath

	return pathEquals && dpEquals && tpEquals
}

// Read reads configuration files as bytes.
func readBytesFromFile(path string) ([]byte, error) {
	configBytes, err := ioutil.ReadFile(path)
	if err != nil {
		msg := fmt.Sprintf("Cannot read configuration file [%s]", path)
		return nil, errors.Wrap(err, msg)
	}
	return configBytes, nil
}
