package main

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Postgres struct {
	DatabaseDNS string `env:"DATABASE_DSN"`
}

type Config struct {
	Postgres
	ListenAddress   string `env:"ADDRESS"`
	LogLevel        string `env:"LOG_LEVEL"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	flag.StringVar(&cfg.ListenAddress, "a", "localhost:8080", "Адрес и порт сервера")
	flag.StringVar(&cfg.LogLevel, "ll", "INFO", "Уровень логирования")
	flag.IntVar(&cfg.StoreInterval, "i", 300, "Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/metrics-db.json", "Полное имя файла, куда сохраняются текущие значения")
	flag.BoolVar(&cfg.Restore, "r", true, "Загружать или нет ранее сохранённые значения из указанного файла при старте сервера")
	flag.StringVar(&cfg.DatabaseDNS, "d", "", "Строка подключения к базе")

	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
