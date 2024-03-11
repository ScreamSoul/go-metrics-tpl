package main

import (
	"flag"
)

type AppFlags struct {
	listenHost string
}

var appFlags = AppFlags{}

func init() {
	flag.StringVar(&appFlags.listenHost, "a", "localhost:8080", "address and port to run server")
}
