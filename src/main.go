package main

import (
	"flag"
	"log"
	"sync"
	"time"

	"dorg/config"
	"dorg/loop"
)

func main() {
	downloadsPath := flag.String("d", "~/Downloads", "A path do downloads directory")
	flag.Parse()

	cnf := config.Config{"", *downloadsPath, *downloadsPath}
	configFile, err := cnf.CreateFile()
	if err != nil {
		log.Fatal(err)
	}

	cfErr := cnf.SaveToFile(configFile)
	if cfErr != nil {
		log.Fatal(cfErr)
	}

	newConfigChan := make(chan config.Config)
	newConfigErrChan := make(chan error)

	cr := config.Reloader{cnf.Filepath, 5 * time.Second, cnf}
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
