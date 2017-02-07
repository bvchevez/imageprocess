package main

// Config deals with all things configuration.
import (
	"flag"
	"fmt"
	"os"

	cnf "github.com/bvchevez/imageprocess/config"
	log "github.com/Sirupsen/logrus"
	"github.com/rakyll/globalconf"
)

var (
	healthcheckToken string
	useSSL           bool
	useCDN           bool
	config           *Config

	// New Relic variables.
	newRelicKey     string
	newRelicAppName string
)

// InitConfigurations is called after the flags has been initialized.
// It takes the config file and binds the configuration to the flags.
func InitConfigurations(configFile string) error {
	// Initialize configuration, reading from environment variables using a 'HIPS_' prefix first,
	// then moving to a static configuration file, usually located in '/etc/hips/hips.conf'.
	conf, err := globalconf.NewWithOptions(&globalconf.Options{
		Filename:  configFile,
		EnvPrefix: "HIPS_",
	})
	if err != nil {
		return err
	}

	conf.ParseAll()

	// Below are for special cases where env var cannot start with "HIPS_"
	newRelicKey = os.Getenv("NEW_RELIC_LICENSE_KEY")
	newRelicAppName = os.Getenv("NEW_RELIC_APP_NAME")
	healthcheckToken = os.Getenv("HEALTHCHECK_TOKEN")
	useSSL = (os.Getenv("USE_SSL") == "1")
	useCDN = (os.Getenv("USE_CDN") == "1")

	cnf.Init(useCDN)

	return nil
}

// Defines a singular configuration struct
// Used to set up all supported sites
// Should only be instantiated/called once
type Config struct {
	//our port number we're listening from.
	port *string

	// Configuration flags
	logLevel         *string
	sentryKey        *string
	sentrySecret     *string
	sentryProjectId  *string
	surrogateControl *string
	cacheControl     *string
	throttle         *string
	concurrency      *string
	burst            *string
	defaultQuality   *string
	bicubicThreshold *string

	//server options
	serverReadTimeout  *string
	serverWriteTimeout *string
}

func (c *Config) Init() {
	c.port = flag.String("port", "6116", "Port we're listening off")
	c.surrogateControl = flag.String("surrogate-control", "max-age=31536000", "Cache control for Fastly")
	c.cacheControl = flag.String("cache-control", "max-age=31536000", "Cache control for browser")
	c.logLevel = flag.String("log-level", "development", "Log level")
	c.sentryKey = flag.String("sentry-key", "", "Sentry account key")
	c.sentrySecret = flag.String("sentry-secret", "", "Sentry account secret")
	c.sentryProjectId = flag.String("sentry-project-id", "", "Sentry project id")
	c.throttle = flag.String("throttle", "0", "Throttle swith. '1' means on, everything else menas off.")
	c.concurrency = flag.String("concurrency", "20", "Throttle concurrency limit per second")
	c.burst = flag.String("burst", "100", "Throttle max burst size.")
	c.serverReadTimeout = flag.String("server-read-timeout", "60", "Throttle max burst size.")
	c.serverWriteTimeout = flag.String("server-write-timeout", "60", "Throttle max burst size.")
	c.defaultQuality = flag.String("default-quality", "95", "Default output-quality for images.")
	c.bicubicThreshold = flag.String("bicubic-threshold", "300", "Minimum pixels in width we want before converting to bicubic.")
}

// getSite takes a string representing the site to get
// and returns it's corresponding site string match (if any)
func (c *Config) GetSite(site string) string {
	scheme := "http"
	if useSSL {
		scheme = "https"
	}

	for k, v := range cnf.AllowedSites {
		if v == site {
			return fmt.Sprintf("%s://%s", scheme, k)
		}
	}

	return fmt.Sprintf("%s://%s", scheme, site)
}

// InitLogLevel initializes log level. For development, we're logging debug, else we only log a minimum of info.
func InitLogLevel() {
	switch *config.logLevel {
	case "development":
		log.SetLevel(log.DebugLevel)
	case "staging":
		log.SetLevel(log.InfoLevel)
	case "production":
		log.SetLevel(log.WarnLevel)
	}
}

// IsSupportedSite makes sure that routeSite is in the whitelist of sites supported.
func IsSupportedSite(routeSite string) bool {

	for k, v := range cnf.AllowedSites {
		// if route is a domain, we check if it exists as a map key.
		if k == routeSite {
			return true
		}

		// if route is a hips path, we check if it exists as a map value.
		if v == routeSite {
			return true
		}
	}
	return false
}

func init() {
	config = &Config{}
}
