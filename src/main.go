package main

import (
	"flag"
)

func main() {
	downloadsPath := flag.String("d", "~/Downloads", "A path do downloads directory")
	logIntoFile := flag.Bool("lf", false, "Logs will be saved into file. Default is .dorglog")
	logFile := flag.String("logFile", ".dorglog", "Path to log file")
	flag.Parse()

	setupLog(*logIntoFile, *logFile)
	cnf := setupConfig(*downloadsPath, *downloadsPath)

	runMainServices(cnf)
}
