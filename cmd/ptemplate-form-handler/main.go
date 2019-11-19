package main

import (
	"flag"
	"fmt"
	"github.com/Miguel-Dorta/logolang"
	"github.com/Miguel-Dorta/ptemplate-form-handler/internal"
	"github.com/Miguel-Dorta/ptemplate-form-handler/pkg/server"
	"os"
	"strconv"
)

var (
	configPath string
	port       int
	log *logolang.Logger
)

func init() {
	log = logolang.NewLogger()
	log.Level = logolang.LevelInfo

	var verbose, version bool
	flag.StringVar(&configPath, "config", "config.toml", "Path to config file")
	flag.IntVar(&port, "port", 8080, "Port to listen")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")
	flag.BoolVar(&version, "version", false, "Print version and exit")
	flag.Parse()

	if version {
		fmt.Println(internal.Version)
		os.Exit(0)
	}

	if verbose {
		log.Level = logolang.LevelDebug
	}

	if port < 1 || port > 65535 {
		log.Criticalf("invalid port")
		os.Exit(1)
	}
}

func main() {
	server.Log = log
	server.Run(configPath, strconv.Itoa(port))
}
