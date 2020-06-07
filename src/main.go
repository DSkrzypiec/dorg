package main

import (
	"flag"
	"fmt"
	"time"

	"dorg/config"
)

func main() {
	downloadsPath := flag.String("d", "~/Downloads", "A path do downloads directory")
	flag.Parse()

	config := config.Default(*downloadsPath)

	for {
		time.Sleep(1 * time.Second)
		fmt.Printf("Waiting for new files to organize...[%s]\n", config.DownloadsPath)
	}
}
