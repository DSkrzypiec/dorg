// Package loop contains implementation of the main program loop which scans
// and moves files from "downloads" catalog.
package loop

import (
	"dorg/config"
	"dorg/dir"
	"time"

	"github.com/rs/zerolog/log"
)

type MainInputs struct {
	InitialConfig config.Config
	DirDiff       chan dir.Dir
	DirDiffErr    chan error
	Config        chan config.Config
	ConfigErr     chan error
}

// TODO
func Main(inputs MainInputs) {
	cnf := inputs.InitialConfig

	for {
		select {
		case newDirDiff := <-inputs.DirDiff:
			log.Info().Msg("New file!:")
			log.Info().Msg(newDirDiff.String())

		case dirErr := <-inputs.DirDiffErr:
			log.Error().Msgf("Error from DirDiff: %s", dirErr.Error())

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
