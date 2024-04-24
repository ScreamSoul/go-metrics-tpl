package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/alexflint/go-arg"
)

type Server struct {
	ListenServerHost string          `arg:"-a,env:ADDRESS" default:"localhost:8080" help:"Адрес и порт сервера"`
	CompressRequest  bool            `arg:"-z,env:COMPRESS_REQUEST" default:"true" help:"compress body request"`
	BackoffIntervals []time.Duration `arg:"--b-intervals,env:BACKOFF_INTERVALS" help:"Интервалы повтора запроса (default=1s,3s,5s)"`
	BackoffRetries   bool            `arg:"--backoff,env:BACKOFF_RETRIES" default:"true" help:"Повтор запроса при разрыве соединения"`
	HashBodyKey      string          `arg:"-k,env:KEY" default:"" help:"hash key"`
}

type Config struct {
	Server
	ReportInterval int    `arg:"-r,env:REPORT_INTERVAL" default:"10" help:"the frequency of sending metrics to the server"`
	PollInterval   int    `arg:"-p,env:POLL_INTERVAL" default:"2" help:"the frequency of polling metrics from the runtime package"`
	LogLevel       string `arg:"--ll,env:LOG_LEVEL" default:"INFO" help:"log level"`
}

func (c *Config) GetServerURL() string {
	return strings.TrimRight(fmt.Sprintf("http://%s", c.Server.ListenServerHost), "/")

}

func (c *Config) GetUpdateMetricURL() string {
	return fmt.Sprintf("%s/updates/", c.GetServerURL())
}

func NewConfig() (*Config, error) {
	var cfg Config

	arg.MustParse(&cfg)

	if cfg.Server.BackoffIntervals == nil && cfg.Server.BackoffRetries {
		cfg.Server.BackoffIntervals = []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
	} else if !cfg.Server.BackoffRetries {
		cfg.Server.BackoffIntervals = nil
	}

	return &cfg, nil
}
