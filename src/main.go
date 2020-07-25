package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"dorg/config"
)

func main() {
	downloadsPath := flag.String("d", "~/Downloads", "A path do downloads directory")
	flag.Parse()

	var cnf config.Config
	cnf = config.Default(*downloadsPath)
	cfErr := cnf.SaveToFile()
	if cfErr != nil {
		log.Fatal(cfErr)
	}

	newConfigChan := make(chan config.Config)
	newConfigErrChan := make(chan error)

	cr := config.Reloader{cnf.Filepath, 5 * time.Second, cnf}
	go cr.ListenAndReload(newConfigChan, newConfigErrChan)

	for {
		select {
		case newCnf := <-newConfigChan:
			cnf = newCnf
		case err := <-newConfigErrChan:
			log.Fatal(err)
		default:
			time.Sleep(1 * time.Second)
			fmt.Printf("Waiting for new files to organize...[%s]\n", cnf.DownloadsPath)
		}
	}
}
