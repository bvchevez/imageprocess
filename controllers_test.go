package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// go test -run Test_hipsController_goodRequest -v
func Test_hipsController_goodRequest(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/hips/mock-unsupported-site/foo.jpeg?resize=100:*", nil)

	hipsController(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "[\"Invalid site [mock-unsupported-site].\"]", w.Body.String())
}

// go test -run Test_imageController_badRequest -v
func Test_hipsController_badRequest(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/hips/bad-site", nil)

	hipsController(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// go test -run Test_hipsController_methodNotAllowedPOST -v
func Test_hipsController_methodNotAllowedPOST(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/hips/bad-site", nil)

	hipsController(w, r)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

// go test -run Test_hipsController_methodNotAllowedPATCH -v
func Test_hipsController_methodNotAllowedPATCH(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/hips/bad-site", nil)

	hipsController(w, r)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

// go test -run Test_hipsController_methodNotAllowedPUT -v
func Test_hipsController_methodNotAllowedPUT(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "/hips/bad-site", nil)

	hipsController(w, r)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

// go test -run Test_hipsController_methodNotAllowedDELETE -v
func Test_hipsController_methodNotAllowedDELETE(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/hips/bad-site", nil)

	hipsController(w, r)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

// go test -run Test_indexController_robotsRequest -v
func Test_indexController_robotsRequest(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/robots.txt", nil)

	indexController(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))
}

// go test -run Test_indexController_faviconRequest -v
func Test_indexController_faviconRequest(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/favicon.ico", nil)

	indexController(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "image/vnd.microsoft.icon", w.Header().Get("Content-Type"))
}

// go test -run Test_indexController_invalidRootRequest -v
func Test_indexController_invalidRootRequest(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/favicon.bad", nil)

	indexController(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// go test -run Test_indexController_validRootRequest -v
func Test_indexController_validRootRequest(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)

	indexController(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}

// go test -run Test_indexController_goodRequest -v
func Test_indexController_goodRequest(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/mock-unsupported-site/foo.jpeg?resize=100:*", nil)

	indexController(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "[\"Invalid site [mock-unsupported-site].\"]", w.Body.String())
}

// go test -run Test_indexController_badRequest -v
func Test_indexController_badRequest(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/bad-site", nil)

	indexController(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// go test -run Test_indexController_methodNotAllowedPOST -v
func Test_indexController_methodNotAllowedPOST(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/hips/bad-site", nil)

	indexController(w, r)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

// go test -run Test_indexController_methodNotAllowedPATCH -v
func Test_indexController_methodNotAllowedPATCH(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/hips/bad-site", nil)

	indexController(w, r)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

// go test -run Test_indexController_methodNotAllowedPUT -v
func Test_indexController_methodNotAllowedPUT(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "/hips/bad-site", nil)

	indexController(w, r)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

// go test -run Test_indexController_methodNotAllowedDELETE -v
func Test_indexController_methodNotAllowedDELETE(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/hips/bad-site", nil)

	indexController(w, r)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

// go test -run Test_healthController_methodNotAllowedPOST -v
func Test_healthController_methodNotAllowedPOST(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/health?token=test-123", nil)

	healthController(w, r)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

// go test -run Test_healthController_methodNotAllowedPATCH -v
func Test_healthController_methodNotAllowedPATCH(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/health?token=test-123", nil)

	healthController(w, r)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

// go test -run Test_healthController_methodNotAllowedPUT -v
func Test_healthController_methodNotAllowedPUT(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "/health?token=test-123", nil)

	healthController(w, r)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

// go test -run Test_healthController_methodNotAllowedDELETE -v
func Test_healthController_methodNotAllowedDELETE(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/health?token=test-123", nil)

	healthController(w, r)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

// go test -run Test_healthControllerGET -v
func Test_healthControllerGET(t *testing.T) {
	healthcheckToken = "test-123"
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/health?token=test-123", nil)
	healthController(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

// go test -run Test_healthControllerHEAD -v
func Test_healthControllerHEAD(t *testing.T) {
	healthcheckToken = "test-123"
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("HEAD", "/health?token=test-123", nil)
	healthController(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

// go test -run Test_healthcheckController_badToken -v
func Test_healthcheckController_badToken(t *testing.T) {
	healthcheckToken = "test-1234"

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/health?token=test-123", nil)
	healthController(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// go test . -run Test_indexController_HEAD_goodRequest -v
func Test_indexController_HEAD_goodRequest(t *testing.T) {

	if (config.port == nil) {
		config.Init()
	}
	// initialize stuff.
	err := InitConfigurations("config/hips.conf")
	assert.Equal(t, nil, err)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("HEAD", "/hdm-dev/images/fear-and-loathing-in-las-vegas-1448994091.jpg?resize=768:*", nil)

	indexController(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "768:327", w.Header().Get("X-Image-Dimensions"))
	assert.Equal(t, "", w.Body.String())
}
