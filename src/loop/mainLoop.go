// Package loop contains implementation of the main program loop which scans
// and moves files from "downloads" catalog.
package loop

import (
	"dorg/config"
	"time"

	"github.com/rs/zerolog/log"
)

type MainInputs struct {
	InitialConfig config.Config
	Config        chan config.Config
	ConfigErr     chan error
}

// TODO
func Main(inputs MainInputs) {
	cnf := inputs.InitialConfig

	for {
		select {
		case newCnf := <-inputs.Config:
			cnf = newCnf
			log.Info().Bool("Configuration changed", true).Fields(map[string]interface{}{
				"New Config": newCnf,
			}).Send()
		case err := <-inputs.ConfigErr:
			log.Fatal().Err(err).Send()
		default:
			time.Sleep(1 * time.Second)
			log.Info().Msgf("Waiting for new files to organize...[%s]", cnf.DownloadsPath)
		}
	}
}
