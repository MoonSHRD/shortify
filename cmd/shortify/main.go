package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MoonSHRD/logger"
	"github.com/MoonSHRD/shortify/app"
	"github.com/MoonSHRD/shortify/config"
	"github.com/MoonSHRD/shortify/router"
)

func main() {
	var configPath string
	var generateConfig bool
	var verboseLogging bool
	var syslogLogging bool
	flag.StringVar(&configPath, "config", "", "Path to server config")
	flag.BoolVar(&generateConfig, "genconfig", false, "Generate new config")
	flag.BoolVar(&verboseLogging, "verbose", true, "Verbose logging")
	flag.BoolVar(&syslogLogging, "syslog", false, "Log to system logging daemon")
	flag.Parse()
	defer logger.Init("shortify", verboseLogging, syslogLogging, ioutil.Discard).Close() // TODO Make ability to use file for log output
	if generateConfig {
		confStr, err := config.Generate()
		if err != nil {
			log.Fatalf("Cannot generate config! %s", err.Error())
		}
		fmt.Print(confStr)
		os.Exit(0)
	}
	logger.Info("Starting Shortify...")
	if configPath == "" {
		logger.Fatal("Path to config isn't specified!")
	}

	cfg, err := config.Parse(configPath)
	if err != nil {
		logger.Fatal(err)
	}
	app, err := app.NewApp(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	router, err := router.NewRouter(app)
	if err != nil {
		logger.Fatalf("Failed to initialize router: %s", err.Error())
	}

	// CTRL+C handler.
	signalHandler := make(chan os.Signal, 1)
	shutdownDone := make(chan bool, 1)
	signal.Notify(signalHandler, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalHandler
		logger.Info("CTRL+C or SIGTERM received, shutting down Shortify...")
		app.Destroy()
		shutdownDone <- true
	}()

	app.Run(router)

	<-shutdownDone
	os.Exit(0)
}
