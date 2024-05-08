package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/1f349/daisy"
	"github.com/1f349/violet/utils"
	"github.com/google/subcommands"
	"github.com/mrmelon54/exit-reload"
	"log"
	"os"
	"path/filepath"
)

type serveCmd struct{ configPath string }

func (s *serveCmd) Name() string { return "serve" }

func (s *serveCmd) Synopsis() string { return "Serve contacts service" }

func (s *serveCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&s.configPath, "conf", "", "/path/to/config.json : path to the config file")
}

func (s *serveCmd) Usage() string {
	return `serve [-conf <config file>]
  Serve contacts service using information from the config file
`
}

func (s *serveCmd) Execute(_ context.Context, _ *flag.FlagSet, _ ...any) subcommands.ExitStatus {
	log.Println("[Daisy] Starting...")

	if s.configPath == "" {
		log.Println("[Daisy] Error: config flag is missing")
		return subcommands.ExitUsageError
	}

	openConf, err := os.Open(s.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("[Daisy] Error: missing config file")
		} else {
			log.Println("[Daisy] Error: open config file: ", err)
		}
		return subcommands.ExitFailure
	}

	var config daisy.Conf
	err = json.NewDecoder(openConf).Decode(&config)
	if err != nil {
		log.Println("[Daisy] Error: invalid config file: ", err)
		return subcommands.ExitFailure
	}

	configPathAbs, err := filepath.Abs(s.configPath)
	if err != nil {
		log.Fatal("[Daisy] Failed to get absolute config path")
	}
	wd := filepath.Dir(configPathAbs)
	normalLoad(config, wd)
	return subcommands.ExitSuccess
}

func normalLoad(startUp daisy.Conf, wd string) {
	srv := daisy.NewHttpServer(startUp, wd)
	log.Printf("[Daisy] Starting HTTP server on '%s'\n", srv.Addr)
	go utils.RunBackgroundHttp("HTTP", srv)

	exit_reload.ExitReload("Daisy", func() {}, func() {
		// stop http server
		_ = srv.Close()
	})
}
