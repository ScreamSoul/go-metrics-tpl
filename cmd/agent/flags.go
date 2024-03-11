package main

import (
	"flag"
)

type AppFlags struct {
	listenServerHost string
	reportInterval   int
	pollInterval     int
}

var appFlags = AppFlags{}

func init() {
	flag.StringVar(&appFlags.listenServerHost, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&appFlags.reportInterval, "r", 10, "the frequency of sending metrics to the server")
	flag.IntVar(&appFlags.pollInterval, "p", 2, "the frequency of polling metrics from the runtime package")
}
