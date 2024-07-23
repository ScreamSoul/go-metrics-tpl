package client

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/screamsoul/go-metrics-tpl/pkg/utils"
)

type CryptoPublicKey struct {
	Key *rsa.PublicKey
}

type Server struct {
	ListenServerHost string          `arg:"-a,env:ADDRESS" default:"localhost:8080" help:"Адрес и порт сервера" json:"address"`
	CompressRequest  bool            `arg:"-z,env:COMPRESS_REQUEST" default:"true" help:"compress body request"`
	BackoffIntervals []time.Duration `arg:"--b-intervals,env:BACKOFF_INTERVALS" help:"Интервалы повтора запроса (default=1s,3s,5s)"`
	BackoffRetries   bool            `arg:"--backoff,env:BACKOFF_RETRIES" default:"true" help:"Повтор запроса при разрыве соединения"`
	HashBodyKey      string          `arg:"-k,env:KEY" default:"" help:"hash key"`
	CryptoKey        CryptoPublicKey `arg:"--crypto-key,env:CRYPTO_KEY" default:"" help:"the path to the file with the public key" josn:"crypto_key"`
}

func (cpk *CryptoPublicKey) UnmarshalText(b []byte) error {
	keyData, err := os.ReadFile(string(b))
	if err != nil {
		return err
	}

	block, _ := pem.Decode(keyData)
	if block == nil || !strings.Contains(block.Type, "PUBLIC KEY") {
		return errors.New("not find public key in file")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	rsaPub, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return errors.New("key type is not RSA")
	}
	cpk.Key = rsaPub
	return nil
}

type Client struct {
	RateLimit      int    `arg:"-l,env:RATE_LIMIT" default:"1" help:"the number of simultaneous outgoing requests to the server"`
	ReportInterval int    `arg:"-r,env:REPORT_INTERVAL" default:"10" help:"the frequency of sending metrics to the server" json:"report_interval"`
	PollInterval   int    `arg:"-p,env:POLL_INTERVAL" default:"2" help:"the frequency of polling metrics from the runtime package" json:"poll_interval"`
	LogLevel       string `arg:"--ll,env:LOG_LEVEL" default:"INFO" help:"log level"`
}
type Config struct {
	Server
	Client
}

func (c *Config) GetServerURL() string {
	return strings.TrimRight(fmt.Sprintf("http://%s", c.Server.ListenServerHost), "/")

}

func (c *Config) GetUpdateMetricURL() string {
	return fmt.Sprintf("%s/updates/", c.GetServerURL())
}

func NewConfig() (*Config, error) {
	var cfg Config

	utils.FillFromFile(&cfg)

	arg.MustParse(&cfg)

	if cfg.Server.BackoffIntervals == nil && cfg.Server.BackoffRetries {
		cfg.Server.BackoffIntervals = []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
	} else if !cfg.Server.BackoffRetries {
		cfg.Server.BackoffIntervals = nil
	}

	if cfg.RateLimit <= 0 {
		cfg.RateLimit = 1
	}

	return &cfg, nil
}
