package main

// routes.go handles routing logic and starts our listener.
import (
	"net/http"
)

var (
	mux *http.ServeMux
)

// InitRoutes initializes all rotes and maps them to their respecitve handler.
func InitRoutes() {
	mux.Handle(httpStatsMonitor(Middleware("/", indexController)))
	mux.Handle(httpStatsMonitor(Middleware("/hips/", hipsController)))
	mux.Handle(Middleware("/health", healthController))
}

func init() {
	mux = http.NewServeMux()
}
