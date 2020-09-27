package main

import (
	"flag"
	"os"
	"sync"
	"time"

	"dorg/config"
	"dorg/dir"
	"dorg/loop"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	downloadsPath := flag.String("d", "~/Downloads", "A path do downloads directory")
	logIntoFile := flag.Bool("lf", false, "Logs will be saved into file. Default is .dorglog")
	logFile := flag.String("logFile", ".dorglog", "Path to log file")
	flag.Parse()

	setupLog(*logIntoFile, *logFile)

	cnf := config.Config{
		Filepath:       "",
		DownloadsPath:  *downloadsPath,
		TargetBasePath: *downloadsPath,
	}
	configFile, err := cnf.CreateFile()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	cfErr := cnf.SaveToFile(configFile)
	if cfErr != nil {
		log.Fatal().Err(cfErr).Send()
	}

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
