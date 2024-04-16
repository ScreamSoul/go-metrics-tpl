package main

import (
	"time"

	"github.com/alexflint/go-arg"
)

type Postgres struct {
	DatabaseDSN      string          `arg:"-d,env:DATABASE_DSN" default:"" help:"Строка подключения к базе Postgres"`
	BackoffIntervals []time.Duration `arg:"--b-intervals,env:BACKOFF_INTERVALS" help:"Интервалы повтора запроса (обязательно если (default=1s,3s,5s)"`
	BackoffRetries   bool            `arg:"--backoff,env:BACKOFF_RETRIES" default:"true" help:"Повтор запроса при разрыве соединения"`
}

type Config struct {
	Postgres
	ListenAddress   string `arg:"-a,env:ADDRESS" default:"localhost:8080" help:"Адрес и порт сервера"`
	LogLevel        string `arg:"--ll,env:LOG_LEVEL" default:"INFO" help:"Уровень логирования"`
	StoreInterval   int    `arg:"-i,env:STORE_INTERVAL" default:"300" help:"Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск"`
	FileStoragePath string `arg:"-f,env:FILE_STORAGE_PATH" default:"/tmp/metrics-db.json" help:"Полное имя файла, куда сохраняются текущие значения"`
	Restore         bool   `arg:"-r,env:RESTORE" default:"true" help:"Загружать или нет ранее сохранённые значения из указанного файла при старте сервера"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	if err := arg.Parse(&cfg); err != nil {
		return nil, err
	}

	if cfg.Postgres.BackoffIntervals == nil && cfg.Postgres.BackoffRetries {
		cfg.Postgres.BackoffIntervals = []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
	} else if !cfg.Postgres.BackoffRetries {
		cfg.Postgres.BackoffIntervals = nil
	}

	return &cfg, nil
}
