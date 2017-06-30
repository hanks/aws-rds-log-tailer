package main

import (
	"aws-rds-log-tailer/pkg/tailer"
	"aws-rds-log-tailer/pkg/version"
	"flag"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Version: ", version.VERSION)

	outPtr := flag.String("out", "postgresql.log", "path to output esjson file")
	dbIDPtr := flag.String("dbID", "", "db identifier")
	flag.Parse()

	if *outPtr == "" || *dbIDPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	tailer.Execute(*outPtr, *dbIDPtr)
}
