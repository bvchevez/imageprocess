package main

import (
	"fmt"
	"net/http"

	"github.com/bvchevez/imageprocess/helper"
	log "github.com/Sirupsen/logrus"
	"github.com/getsentry/raven-go"
	"github.com/newrelic/go-agent"
	"gopkg.in/throttled/throttled.v2"
	"gopkg.in/throttled/throttled.v2/store/memstore"
)

var (
	newRelicApp newrelic.Application
)

// Middleware installs middlewares on our http handler.
func Middleware(route string, fn func(http.ResponseWriter, *http.Request)) (string, http.Handler) {
	next := http.Handler(http.HandlerFunc(fn))

	// attach rate limiter.
	if *config.throttle == "1" {
		next = addRateLimiter(next)
		log.WithFields(log.Fields{
			"route":       route,
			"concurrency": *config.concurrency,
			"burst":       *config.burst,
		}).Info("Rate Limiter added")
	}

	// attach panic recovery + logger.
	if *config.sentryKey != "" && *config.sentrySecret != "" && *config.sentryProjectId != "" {
		next = addPanicLogger(next)
		log.WithFields(log.Fields{
			"route": route,
		}).Info("Panic Logger added")
	}

	return route, next
}

// InitMontoring initializes monitoring handlers.
func InitMontoring() {
	var err error

	if newRelicKey == "" {
		log.Info("New Relic Disabled")
		return
	}

	newRelicApp, err = newrelic.NewApplication(newrelic.NewConfig(newRelicAppName, newRelicKey))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("New Relic failed to initialize.")
	}

	log.WithFields(log.Fields{
		"application name": newRelicAppName,
	}).Info("New Relic Enabled")
}

// httpStatsMonitor monitors http statuses using new relic.
func httpStatsMonitor(route string, handle http.Handler) (string, http.Handler) {
	if newRelicApp == nil {
		return route, handle
	}

	// This works because wraphandle returns (string, http.Handler)
	return newrelic.WrapHandle(newRelicApp, route, handle)
}

// addPanicLogger integrates a panic logger (raven) into our http handler to catch panics automatically.
func addPanicLogger(next http.Handler) http.Handler {
	raven.SetDSN(
		fmt.Sprintf("https://%s:%s@app.getsentry.com/%s",
			*config.sentryKey, *config.sentrySecret, *config.sentryProjectId))
	return http.Handler(http.HandlerFunc(raven.RecoveryHandler(next.ServeHTTP)))
}

// addRateLimiter throttles our request rate to a maximum of *config.concurrency per second.
// with a buffer of *config.burst per method.
func addRateLimiter(next http.Handler) http.Handler {
	store, err := memstore.New(65536)
	if err != nil {
		eventId := raven.CaptureError(err, nil)
		log.WithFields(log.Fields{
			"error":          err.Error(),
			"sentry_eventID": eventId,
		}).Fatal("Fatal Error! Unable to initialize memstore.")
	}

	quota := throttled.RateQuota{
		MaxRate:  throttled.PerSec(helper.String2Int(*config.concurrency)),
		MaxBurst: helper.String2Int(*config.burst),
	}

	rateLimiter, err := throttled.NewGCRARateLimiter(store, quota)
	if err != nil {
		eventId := raven.CaptureError(err, nil)
		log.WithFields(log.Fields{
			"error":          err.Error(),
			"sentry_eventID": eventId,
		}).Fatal("Fatal Error! Unable to initialize NewGCRARateLimiter.")
	}

	httpRateLimiter := throttled.HTTPRateLimiter{
		RateLimiter: rateLimiter,
		VaryBy:      &throttled.VaryBy{Method: true},
	}

	return httpRateLimiter.RateLimit(next)
}
