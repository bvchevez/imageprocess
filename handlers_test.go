package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

// go test -run Test_Healthcheck__no_errors -v
// test that health check isn't breaking.
func Test_Healthcheck__no_errors(t *testing.T) {
	res := HandleHealthCheck()
	assert.Equal(t, http.StatusOK, res.Code)
}
