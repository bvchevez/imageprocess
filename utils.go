package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/bvchevez/imageprocess/helper"
	"github.com/bvchevez/imageprocess/image"
)

// Response represents a JSON response that a handler would use to output
// JSON to the client per request
type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
	Image *image.Image
}

// CustomWriter writes bytes to http response writer, caller can pass whatever content type they want.
func CustomWriter(w http.ResponseWriter, data []byte, ct string) {
	w.Header().Set("Content-Type", ct)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// JsonWriter writes json response from response to response writer
func JsonWriter(w http.ResponseWriter, res *Response) {
	b, err := json.Marshal(res.Data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{error: "%s"}`, err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(res.Code)
	w.Write(b)
}

// ImageHeaderWriter writes the header info for a given image
func ImageHeaderWriter(w http.ResponseWriter, res *Response) {
	w.Header().Set("Content-Type", res.Image.Type)
	w.Header().Set("Surrogate-Control", *config.surrogateControl)
	w.Header().Set("Cache-Control", *config.cacheControl)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", res.Image.Size))
	w.Header().Set("X-Image-Dimensions", fmt.Sprintf("%d:%d", res.Image.Width, res.Image.Height))
	w.Header().Set("X-Source-Image-Dimensions",
		fmt.Sprintf("%d:%d", res.Image.SourceWidth, res.Image.SourceHeight))

	// If this image is gif, we mark this as animated.
	// Currently the only test for animation is GIF.
	// We can add more criterias as we move along.
	if res.Image.Animated {
		w.Header().Set("X-Animated", "1")
	} else {
		w.Header().Set("X-Animated", "0")
	}

	w.WriteHeader(http.StatusOK)
}

// ImageWriter writes image data from the response to response writer.
func ImageWriter(w http.ResponseWriter, res *Response) {
	ImageHeaderWriter(w, res)
	w.Write(res.Image.Data)
}

// extract site/path info from given path.
// Always returns a slice of exactly two strings, site/path.
// return two blank strings on error.
func ExtractInfoFromPath(path string) ([]string, error) {
	s := strings.SplitN(strings.TrimLeft(path, "/"), "/", 2)
	if len(s) != 2 {
		return nil, fmt.Errorf("Invalid path [%s]", path)
	}

	//make sure the path portion always have a single left slash.
	s[1] = "/" + strings.TrimLeft(s[1], "/")
	return s, nil

}

// getImage takes site and path strings and attempts to return an Image struct (defined in pipeline.go)
// Will first check for a matching site within the map of S3 configs, and fall back to checking site strings
func GetImage(site, path, pipelineID string) ([]byte, error) {
	defer helper.Timer(helper.TimerPayload{
		Start: time.Now(),
		Name:  "(" + pipelineID + ") getImage",
	})

	conifgSite := config.GetSite(site)
	if len(conifgSite) != 0 {
		return GetImageFromUrl(conifgSite + path)
	} else {
		return nil, fmt.Errorf("A proper source destination was not found. Source was: " + site)
	}
}

// getImageFromUrl takes a URL string to be retrived
// It will attempt to retrieve a resource (in this case an image) from the URL
func GetImageFromUrl(url string) ([]byte, error) {
	// Make request for resource
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error getting image: %v", err)
	}

	// Close the response body after this function returns
	defer res.Body.Close()

	// Check for 200 status code
	// S3 does not always send 404s only
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Source returned a status code other than 200: %d", res.StatusCode)
	}

	// Read in data from response body
	data, err := ioutil.ReadAll(res.Body)

	// Check for errors reading response body
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %v", err)
	}

	return data, nil
}
