package main

import (
	"github.com/bvchevez/imageprocess/image"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// go test -run Test_ExtractSiteFromPath_succsess -v
func Test_ExtractSiteFromPath_succsess(t *testing.T) {
	s, err := ExtractInfoFromPath("/hdm-dev/images/fear-and-loathing-in-las-vegas-1448994091.jpg")
	assert.Equal(t, "hdm-dev", s[0])
	assert.Equal(t, "/images/fear-and-loathing-in-las-vegas-1448994091.jpg", s[1])
	assert.Equal(t, nil, err)
}

// go test -run Test_ExtractSiteFromPath_failureEdgeCases -v
func Test_ExtractSiteFromPath_failureEdgeCases(t *testing.T) {
	s, err := ExtractInfoFromPath("/hdm-dev/")
	assert.Equal(t, nil, err)
	assert.Equal(t, "hdm-dev", s[0])
	assert.Equal(t, "/", s[1])

	_, err2 := ExtractInfoFromPath("/hdm-dev")
	assert.Equal(t, "Invalid path [/hdm-dev]", err2.Error())

	_, err3 := ExtractInfoFromPath("/")
	assert.Equal(t, "Invalid path [/]", err3.Error())

	_, err4 := ExtractInfoFromPath("")
	assert.Equal(t, "Invalid path []", err4.Error())
}

// go test -run Test_GetImageFromUrl_badUrl -v
func Test_GetImageFromUrl_badUrl(t *testing.T) {
	b, err := GetImageFromUrl("some-bad-url")

	assert.Equal(t, `Error getting image: Get some-bad-url: unsupported protocol scheme ""`, err.Error())
	assert.Equal(t, []byte(nil), b)
}

// go test -run Test_JsonWriter__200Response -v
func Test_JsonWriter__200Response(t *testing.T) {
	res := &Response{Code: http.StatusOK, Data: "Test Data"}

	w := httptest.NewRecorder()
	JsonWriter(w, res)
	assert.Equal(t, "\"Test Data\"", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

// go test -run Test_JsonWriter__400Response -v
func Test_JsonWriter__400Response(t *testing.T) {
	res := &Response{Code: http.StatusBadRequest, Data: "Test Data"}

	w := httptest.NewRecorder()
	JsonWriter(w, res)
	assert.Equal(t, "\"Test Data\"", w.Body.String())
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

// go test -run Test_JsonWriter__FailedJsonMarhsal -v
func Test_JsonWriter__FailedJsonMarhsal(t *testing.T) {
	res := &Response{
		Code: http.StatusBadRequest,
		Data: map[float64]int{2.5: 1},
	}

	w := httptest.NewRecorder()
	JsonWriter(w, res)
	assert.Equal(t, "{error: \"json: unsupported type: map[float64]int\"}", w.Body.String())
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

// go test -run Test_ImageWriter -v
func Test_ImageWriter(t *testing.T) {
	if (config.port == nil) {
		config.Init()
	}
	*config.surrogateControl = "max-age=12345"
	*config.cacheControl = "max-age=54321"

	res := &Response{
		Code: http.StatusOK,
		Data: nil,
		Image: &image.Image{
			Data:   nil,
			Type:   "mock/type",
			Size:   100,
			Width:  15,
			Height: 20,
		},
	}

	w := httptest.NewRecorder()
	ImageWriter(w, res)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "", w.Body.String())

	//verifying headers are correctly set.
	assert.Equal(t, "mock/type", w.Header().Get("Content-Type"))
	assert.Equal(t, "max-age=12345", w.Header().Get("Surrogate-Control"))
	assert.Equal(t, "max-age=54321", w.Header().Get("Cache-Control"))
	assert.Equal(t, "100", w.Header().Get("Content-Length"))
}
