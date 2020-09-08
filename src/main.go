package main

import (
	"flag"
	"sync"
	"time"

	"dorg/config"
	"dorg/loop"

	"github.com/rs/zerolog/log"
)

func main() {
	downloadsPath := flag.String("d", "~/Downloads", "A path do downloads directory")
	flag.Parse()

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
		log.Fatal().Err(err).Send()
	}

	newConfigChan := make(chan config.Config)
	newConfigErrChan := make(chan error)

	cr := config.Reloader{
		ConfigFilePath: cnf.Filepath,
		ReloadInterval: 5 * time.Second,
		CurrentConfig:  cnf,
	}

	mainLoopInputs := loop.MainInputs{
		InitialConfig: cnf,
		Config:        newConfigChan,
		ConfigErr:     newConfigErrChan,
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go cr.ListenAndReload(newConfigChan, newConfigErrChan)
	go loop.Main(mainLoopInputs)

	wg.Wait()
}
