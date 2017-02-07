package main

// handler.go handles requests parameters and writes the result to response writer.
import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	cnf "github.com/bvchevez/imageprocess/config"
	log "github.com/Sirupsen/logrus"
	"github.com/newrelic/go-agent"
)

var (
	start time.Time
)

// HealthCheckResponse represents the health check data response (in JSON)
// that this service should output when a health check is requested
type Health struct {
	GifsicleVersion string   `json:"gifsicle_info"`
	LibVipsVersion  string   `json:"libvips_info"`
	GolangVersion   string   `json:"go_info"`
	Goroutines      int      `json:"goroutines"`
	NumberOfCPUs    int      `json:"cpus"`
	UpTime          int64    `json:"uptime"`
	StatusUpdated   string   `json:"status_updated"`
	Errors          []string `json:"errors"`
}

// HandleImage handles image request and outputs a Response pointer.
func HandleImage(site, path, params string, txn newrelic.Transaction) *Response {
	site = cnf.NormalizeSite(site)
	if IsSupportedSite(site) == false {
		err := fmt.Errorf("Invalid site [%s].", site)
		log.WithFields(log.Fields{
			"error": err.Error(),
			"site":  site,
		}).Warn("Unsupported site.")

		return &Response{
			Code:  http.StatusBadRequest,
			Data:  [1]string{err.Error()},
			Image: nil,
		}
	}

	pipeline := &Pipeline{
		site:     site,
		path:     path,
		rawQuery: params,
	}
	image, resp := pipeline.Process(txn)
	if resp != nil {
		log.WithFields(log.Fields{
			"error":  resp.Data,
			"site":   site,
			"path":   path,
			"params": params,
		}).Warn("Image processing error.")

		return resp
	}

	return &Response{
		Code:  http.StatusOK,
		Data:  nil,
		Image: image,
	}
}

// HandleHealthCheck handles healthcheck requests and outputs a Response pointer.
func HandleHealthCheck() *Response {
	var (
		errors      []string
		err         error
		gifsicleCmd string
		libVipsCmd  string
		goCmd       string
	)

	// get Gifsicle information
	gifsicleCmd, err = exec.LookPath("gifsicle")
	if err != nil {
		errors = append(errors, "Gifsicle not installed")
	}

	gifsicleVersion, _ := exec.Command(
		gifsicleCmd,
		"--version",
	).Output()

	// Get LibVips information
	libVipsCmd, err = exec.LookPath("vips")
	if err != nil {
		errors = append(errors, "LibVips not installed")
	}

	libVipsVersion, _ := exec.Command(
		libVipsCmd,
		"--version",
	).Output()

	// Get Go information
	goCmd, err = exec.LookPath("go")
	if err != nil {
		errors = append(errors, "Go not installed")
	}

	goVersion, _ := exec.Command(
		goCmd,
		"version",
	).Output()

	// Get Uptime
	uptime := time.Now().Unix() - start.Unix()

	// Return new healthcheck response
	// Check for healthcheck in site
	return &Response{
		Code: http.StatusOK,
		Data: Health{
			GifsicleVersion: strings.Split(string(gifsicleVersion), "\n")[0],
			LibVipsVersion:  strings.Split(string(libVipsVersion), "\n")[0],
			GolangVersion:   strings.Split(string(goVersion), "\n")[0],
			UpTime:          uptime,
			Goroutines:      runtime.NumGoroutine(),
			NumberOfCPUs:    runtime.NumCPU(),
			StatusUpdated:   time.Now().Format("Mon Jan 2 15:04:05 MST 2006"),
			Errors:          errors,
		},
		Image: nil,
	}
}

func init() {
	start = time.Now()
}
