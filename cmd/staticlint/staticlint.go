package main

import (
	"strings"

	"github.com/kisielk/errcheck/errcheck"
	"github.com/screamsoul/go-metrics-tpl/internal/analysis/exitinmain"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

type StaticcheckConfig struct {
	prefixAnalizers []string
}

func NewStaticcheckConfig() *StaticcheckConfig {
	return &StaticcheckConfig{
		prefixAnalizers: []string{
			"SA",
			"S1030",
			"ST1000",
			"QF1001",
		},
	}
}

func (conf StaticcheckConfig) CheckPattern(analizerName string) bool {
	for _, p := range conf.prefixAnalizers {
		if strings.HasPrefix(analizerName, p) {
			return true
		}
	}
	return false
}

func (conf StaticcheckConfig) GetAnalizers() (analizers []*analysis.Analyzer) {
	for _, v := range staticcheck.Analyzers {
		if conf.CheckPattern(v.Analyzer.Name) {
			analizers = append(analizers, v.Analyzer)
		}
	}
	return
}

func main() {
	var mychecks []*analysis.Analyzer

	mychecks = append(mychecks, printf.Analyzer)
	mychecks = append(mychecks, shadow.Analyzer)
	mychecks = append(mychecks, structtag.Analyzer)
	mychecks = append(mychecks, errcheck.Analyzer)
	mychecks = append(mychecks, exitinmain.ExitInMainCheckAnalyzer)

	mychecks = append(mychecks, NewStaticcheckConfig().GetAnalizers()...)

	multichecker.Main(
		mychecks...,
	)
}
