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

	config := config.Default(*downloadsPath)
	_, err := config.CreateFileIfNotExist()
	if err != nil {
		log.Fatal("Cannot create config file.", err.Error())
	}

	for {
		time.Sleep(1 * time.Second)
		fmt.Printf("Waiting for new files to organize...[%s]\n", config.DownloadsPath)
	}
}
