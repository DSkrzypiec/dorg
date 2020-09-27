package main

import (
	"dorg/config"
	"dorg/dir"
	"dorg/loop"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Starts dorg main services.
// TODO: extend after all services will be added.
func runMainServices(cnf config.Config) {
	newConfigChan := make(chan config.Config)
	newConfigErrChan := make(chan error)
	newDirDiffChan := make(chan dir.Dir)
	newDirErrChan := make(chan error)

	cr := config.Reloader{
		ConfigFilePath: cnf.Filepath,
		ReloadInterval: 5 * time.Second,
		CurrentConfig:  cnf,
	}

	dirListener, dlErr := dir.NewDirListener(cnf)
	if dlErr != nil {
		log.Fatal().Err(dlErr).Send()
	}

	mainLoopInputs := loop.MainInputs{
		InitialConfig: cnf,
		Config:        newConfigChan,
		ConfigErr:     newConfigErrChan,
		DirDiff:       newDirDiffChan,
		DirDiffErr:    newDirErrChan,
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go cr.ListenAndReload(newConfigChan, newConfigErrChan)
	go dirListener.Listen(newDirDiffChan, newDirErrChan)
	go loop.Main(mainLoopInputs)

	wg.Wait()
}

// Setup config. Create config.Config object and save it into ~/.dorgconfig
// file.
func setupConfig(downloadsPath, targetPath string) config.Config {
	cnf := config.Config{
		DownloadsPath:  downloadsPath,
		TargetBasePath: targetPath,
	}
	configFile, err := cnf.CreateFile()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	cfErr := cnf.SaveToFile(configFile)
	if cfErr != nil {
		log.Fatal().Err(cfErr).Send()
	}

	return cnf
}

// Setup log based on program flags.
func setupLog(logIntoFile bool, logFilePath string) {
	if !logIntoFile {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		return
	}

	logFile, err := os.Create(logFilePath)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	// Just for initial info
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msgf("Logs will be logged into [%s] file.", logFilePath)

	// Setup log file
	log.Logger = zerolog.New(logFile).With().Timestamp().Logger()
}
