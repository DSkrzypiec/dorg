package config

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
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
		log.Println("Checking for config file changes...")
		newConfig, err := TryParseConfig(r.ConfigFilePath)
		log.Println("Config: ", newConfig)
		if err != nil {
			msg := fmt.Sprintf("Cannot parse configuration file [%s]", r.ConfigFilePath)
			log.Println(msg)
			errChan <- errors.Wrap(err, msg)
			// TODO what now?
			continue
		}

		if !r.CurrentConfig.Equals(newConfig) {
			log.Println("Configuration has changed!")
			r.CurrentConfig = newConfig
			configChan <- newConfig
		}
	}
}
