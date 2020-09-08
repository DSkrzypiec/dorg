package config

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Reloader struct {
	ConfigFilePath string
	ReloadInterval time.Duration
	CurrentConfig  Config
}

// ListenAndReload checks dorg configuration files once every ReloadInterval
// for any changes to be updated.
func (r *Reloader) ListenAndReload(configChan chan<- Config, errChan chan<- error) {
	for {
		time.Sleep(r.ReloadInterval)
		log.Info().Msg("Checking for config file changes...")
		configFile, err := r.CurrentConfig.OpenFile()

		if err != nil {
			msg := fmt.Sprintf("Cannot open config file [%s]", r.CurrentConfig.Filepath)
			errChan <- errors.Wrap(err, msg)
		}

		newConfig, err := TryParseConfig(configFile)
		if err != nil {
			msg := fmt.Sprintf("Cannot parse configuration file [%s]", r.ConfigFilePath)
			errChan <- errors.Wrap(err, msg)
			// TODO what now?
			continue
		}

		if !r.CurrentConfig.Equals(newConfig) {
			r.CurrentConfig = newConfig
			configChan <- newConfig
		}
	}
}
