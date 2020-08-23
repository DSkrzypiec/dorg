// Package loop contains implementation of the main program loop which scans
// and moves files from "downloads" catalog.
package loop

import (
	"dorg/config"
	"fmt"
	"log"
	"time"
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
		case err := <-inputs.ConfigErr:
			log.Fatal(err)
		default:
			time.Sleep(1 * time.Second)
			fmt.Printf("Waiting for new files to organize...[%s]\n", cnf.DownloadsPath)
		}
	}
}
