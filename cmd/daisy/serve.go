package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/1f349/daisy"
	"github.com/1f349/violet/utils"
	"github.com/charmbracelet/log"
	"github.com/google/subcommands"
	"github.com/mrmelon54/exit-reload"
	"os"
	"path/filepath"
)

type serveCmd struct {
	configPath string
	debugLog   bool
}

func (s *serveCmd) Name() string { return "serve" }

func (s *serveCmd) Synopsis() string { return "Serve contacts service" }

func (s *serveCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&s.configPath, "conf", "", "/path/to/config.json : path to the config file")
	f.BoolVar(&s.debugLog, "debug", false, "enable debug logging")
}

func (s *serveCmd) Usage() string {
	return `serve [-conf <config file>]
  Serve contacts service using information from the config file
`
}

func (s *serveCmd) Execute(_ context.Context, _ *flag.FlagSet, _ ...any) subcommands.ExitStatus {
	if s.debugLog {
		daisy.Logger.SetLevel(log.DebugLevel)
	}
	daisy.Logger.Info("Starting...")

	if s.configPath == "" {
		daisy.Logger.Error("Config flag is missing")
		return subcommands.ExitUsageError
	}

	openConf, err := os.Open(s.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			daisy.Logger.Error("Missing config file")
		} else {
			daisy.Logger.Error("Open config file", "err", err)
		}
		return subcommands.ExitFailure
	}

	var config daisy.Conf
	err = json.NewDecoder(openConf).Decode(&config)
	if err != nil {
		daisy.Logger.Error("Invalid config file", "err", err)
		return subcommands.ExitFailure
	}

	configPathAbs, err := filepath.Abs(s.configPath)
	if err != nil {
		daisy.Logger.Error("Failed to get absolute config path")
		return subcommands.ExitFailure
	}
	wd := filepath.Dir(configPathAbs)
	normalLoad(config, wd)
	return subcommands.ExitSuccess
}

func normalLoad(startUp daisy.Conf, wd string) {
	srv := daisy.NewHttpServer(startUp, wd)
	daisy.Logger.Infof("Starting HTTP server on '%s'", srv.Addr)
	go utils.RunBackgroundHttp(daisy.Logger, srv)

	exit_reload.ExitReload("Daisy", func() {}, func() {
		// stop http server
		_ = srv.Close()
	})
}
