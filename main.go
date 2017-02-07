package main

import (
	"os"
	"os/signal"
	"runtime"

	log "github.com/Sirupsen/logrus"
)

// Entry point for HIPS. Starts up everything.
func main() {
	//set maxproc to number of cpu.
	runtime.GOMAXPROCS(runtime.NumCPU())

	config.Init()
	configFile := "config/hips.conf"

	// initialize stuff.
	err := InitConfigurations("config/hips.conf")
	if err != nil {
		//can't log into sentry yet... we don't even have config values.
		log.WithFields(log.Fields{
			"config_file": configFile,
			"error":       err.Error(),
		}).Fatal("Fatal Error! Failed to load configuration.")
	}

	InitLogLevel()
	InitMontoring()
	InitRoutes()
	StartServer()

	// Listen for and terminate HIPS on SIGKILL or SIGINT signals.
	sigStop := make(chan os.Signal)
	signal.Notify(sigStop, os.Interrupt, os.Kill)

	select {
	case <-sigStop:
		log.Info("Shutting down server...")
	}
}
