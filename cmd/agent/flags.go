package main

import (
	"flag"
)

type Config struct {
	ListenServerHost string
	ReportInterval   int
	PollInterval     int
}

var cfg Config

func init() {
	flag.StringVar(&cfg.ListenServerHost, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&cfg.ReportInterval, "r", 10, "the frequency of sending metrics to the server")
	flag.IntVar(&cfg.PollInterval, "p", 2, "the frequency of polling metrics from the runtime package")
}

func parseConfig() {
	flag.Parse()
}
