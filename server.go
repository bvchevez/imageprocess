package main

// server.go initializes our listener and starts our server.
import (
	"net/http"
	"time"

	"github.com/bvchevez/imageprocess/helper"
	log "github.com/Sirupsen/logrus"
	"github.com/getsentry/raven-go"
)

var (
	server *http.Server
)

// StartServer sets up and starts our server's listener.
func StartServer() {
	//set up our server.
	server.Addr = ":" + *config.port
	server.Handler = mux
	server.ReadTimeout = time.Duration(helper.String2Int(*config.serverReadTimeout)) * time.Second
	server.WriteTimeout = time.Duration(helper.String2Int(*config.serverWriteTimeout)) * time.Second

	//initialize our server.
	log.WithFields(log.Fields{
		"port": *config.port,
	}).Info("Server listening...")

	err := server.ListenAndServe()
	if err != nil {
		eventId := raven.CaptureError(err, nil)
		log.WithFields(log.Fields{
			"port":           *config.port,
			"error":          err.Error(),
			"sentry_eventID": eventId,
		}).Fatal("Fatal Error! Failed to listen on port.")
	}
}

func init() {
	server = &http.Server{}
}
