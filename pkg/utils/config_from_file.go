package utils

import (
	"encoding/json"
	"os"

	"github.com/alexflint/go-arg"
)

type ConfigFile struct {
	Path string `arg:"-c,--config,env:CONFIG" default:"" help:"Path to the configuration file"`
}

func FillFromFile(cfg interface{}) {
	var cf ConfigFile

	arg.MustParse(&cf)

	if cf.Path == "" {
		return
	}

	data, err := os.ReadFile(cf.Path)

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return
	}
}
