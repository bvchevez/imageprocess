package main

import (
	cnf "github.com/bvchevez/imageprocess/config"
	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	config = &Config{}
}

func resetConfig() {
	useSSL = false
	cnf.AllowedSites = map[string]string{
		// cdn
		"bpc.h-cdn.test.co":  "bestproducts",
		"cad.h-cdn.test.co":  "caranddriver",
		"clv.h-cdn.test2.co": "countryliving",
		"clv.h-cdn.test1.co": "countryliving",
		"cos.h-cdn.co":       "cosmopolitan",
	}
}

// go test -run Test_Supported_site -v
func Test_Supported_site(t *testing.T) {
	resetConfig()
	assert.Equal(t, false, IsSupportedSite("elle"))
	assert.Equal(t, false, IsSupportedSite("google.com"))
	assert.Equal(t, true, IsSupportedSite("caranddriver"))
	assert.Equal(t, true, IsSupportedSite("clv.h-cdn.test1.co"))
	assert.Equal(t, true, IsSupportedSite("clv.h-cdn.test2.co"))
}

// go test -run Test_GetSite_useCDN -v
func Test_GetSite_useCDN(t *testing.T) {
	cnf.Init(true)
	assert.Equal(t, "http://cos.h-cdn.co", config.GetSite("cosmopolitan"))
	assert.Equal(t, "http://amv-prod-cos.s3.amazonaws.com", config.GetSite("amv-prod-cos.s3.amazonaws.com"))
	assert.Equal(t, "http://cos.h-cdn.co", config.GetSite("cos.h-cdn.co"))

	cnf.Init(false)
	assert.Equal(t, "http://amv-prod-cos.s3.amazonaws.com", config.GetSite("cosmopolitan"))
	assert.Equal(t, "http://amv-prod-cos.s3.amazonaws.com", config.GetSite("amv-prod-cos.s3.amazonaws.com"))
	assert.Equal(t, "http://cos.h-cdn.co", config.GetSite("cos.h-cdn.co"))
}

// go test -run Test_GetSite_exists -v
func Test_GetSite_exists(t *testing.T) {
	resetConfig()
	assert.Equal(t, "http://cad.h-cdn.test.co", config.GetSite("caranddriver"))
	assert.Equal(t, "http://cos.h-cdn.co", config.GetSite("cosmopolitan"))
}

// go test -run Test_SSL_toggle -v
func Test_SSL_toggle(t *testing.T) {
	resetConfig()

	// test ssl disabled
	useSSL = false
	assert.Equal(t, "http://cad.h-cdn.test.co", config.GetSite("caranddriver"))

	// test ssl enabled
	useSSL = true
	assert.Equal(t, "https://cad.h-cdn.test.co", config.GetSite("caranddriver"))
}

// go test -run Test_GetSite_assumeDomain -v
func Test_GetSite_assumeDomain(t *testing.T) {
	resetConfig()

	assert.Equal(t, "http://bpc.h-cdn.test.co", config.GetSite("bpc.h-cdn.test.co"))
	//elle does not exist in the sites array, we assume this is a domain.
	assert.Equal(t, "http://elle", config.GetSite("elle"))
}

// go test -run Test_InitLogLevel__production -v
func Test_InitLogLevel__production(t *testing.T) {
	config.logLevel = new(string)
	*config.logLevel = "production"
	InitLogLevel()
	assert.Equal(t, log.WarnLevel, log.GetLevel())
}

// go test -run Test_InitLogLevel__stage -v
func Test_InitLogLevel__stage(t *testing.T) {
	config.logLevel = new(string)
	*config.logLevel = "staging"
	InitLogLevel()
	assert.Equal(t, log.InfoLevel, log.GetLevel())
}

// go test -run Test_InitLogLevel__develop -v
func Test_InitLogLevel__develop(t *testing.T) {
	config.logLevel = new(string)
	*config.logLevel = "development"
	InitLogLevel()
	assert.Equal(t, log.DebugLevel, log.GetLevel())
}

// go test -run Test_InitConfiguration__badConfig -v
func Test_InitConfiguration__badConfig(t *testing.T) {
	err := InitConfigurations("/Bad/Path")
	assert.Equal(t, "open /Bad/Path: no such file or directory", err.Error())
}

// go test -run Test_InitConfiguration__goodConfig -v
func Test_InitConfiguration__goodConfig(t *testing.T) {
	err := InitConfigurations("fixtures/test.config")
	assert.Equal(t, nil, err)
}
