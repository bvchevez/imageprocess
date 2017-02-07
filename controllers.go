package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bvchevez/imageprocess/helper"

	log "github.com/Sirupsen/logrus"
	"github.com/newrelic/go-agent"
)

// default response for method not allowed.
var (
	methodNotAllowedResponse = &Response{
		Code:  http.StatusMethodNotAllowed,
		Data:  []string{"Invalid method. Only GET is allowed."},
		Image: nil,
	}

	// default response for bad request.
	badRequestResponse = &Response{
		Code:  http.StatusBadRequest,
		Data:  []string{"Invalid image path format, must be {hips-domain}/{image-source-domain}/{path-to-image-source}?{parameters}"},
		Image: nil,
	}

	// default response for root.
	defaultRootResponse = &Response{
		Code:  http.StatusOK,
		Data:  []string{"Image path format must be {hips-domain}/{image-source-domain}/{path-to-image-source}?{parameters}"},
		Image: nil,
	}

	// ctMapper is a content type mapper for the root file requested.
	// It also acts as a whitelist for allowed file access.
	ctMapper = map[string]string{
		"/robots.txt":  "text/plain",
		"/favicon.ico": "image/vnd.microsoft.icon",
	}
)

// indexController is controller for root request.
func indexController(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" && req.Method != "HEAD" {
		JsonWriter(w, methodNotAllowedResponse)
		return
	}

	if req.URL.Path == "/" {
		JsonWriter(w, defaultRootResponse)
		return
	}

	// logic to get root files, such as robots.txt or favicon.ico.
	if contentType, ok := ctMapper[req.URL.Path]; ok {
		data, _ := ioutil.ReadFile(fmt.Sprintf(".%s", req.URL.Path))
		CustomWriter(w, data, contentType)
		return
	}

	if res, err := imageControllerHelper(w, req.URL.Path, req.URL.RawQuery); err != nil {
		JsonWriter(w, res)
	} else if req.Method == "HEAD" {
		ImageHeaderWriter(w, res)
	} else {
		ImageWriter(w, res)
	}
}

// hipsController is a controller for legacy url, if an url starts with /hips/ we take care of it here. (will eventually get depricated)
func hipsController(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" && req.Method != "HEAD" {
		JsonWriter(w, methodNotAllowedResponse)
		return
	}

	reqURLPath := strings.Replace(req.URL.Path, "/hips", "", 1)
	if res, err := imageControllerHelper(w, reqURLPath, req.URL.RawQuery); err != nil {
		JsonWriter(w, res)
	} else if req.Method == "HEAD" {
		ImageHeaderWriter(w, res)
	} else {
		ImageWriter(w, res)
	}
}

// imageControllerHelper is a helper that handles all common stuff between hips and index controller.
func imageControllerHelper(w http.ResponseWriter, path, params string) (*Response, error) {
	pathInfo, err := ExtractInfoFromPath(path)
	if pathInfo == nil || len(pathInfo) != 2 || err != nil {
		return badRequestResponse, err
	}

	txn, txnOk := w.(newrelic.Transaction)
	ueParams, err := helper.UnescapeURL(params)
	if err != nil {
		return badRequestResponse, err
	}

	res := HandleImage(pathInfo[0], pathInfo[1], ueParams, txn)
	if res.Code != http.StatusOK {
		if txnOk {
			err := txn.AddAttribute("Original Path", fmt.Sprintf("%s?%s", path, ueParams))
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Warn("Failed to add attributes to NR Transaction.")
			}
		}
		return res, fmt.Errorf("Response status not OK")
	}

	return res, nil
}

// healthController is a controller for /health path.
func healthController(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" && req.Method != "HEAD" {
		JsonWriter(w, methodNotAllowedResponse)
		return
	}

	q := req.URL.Query()
	if (len(q["token"]) > 0) && (q["token"][0] == healthcheckToken) {
		res := HandleHealthCheck()
		JsonWriter(w, res)
		return
	}

	log.WithFields(log.Fields{
		"query": q,
	}).Warn("Invalid healthcheck token.")

	w.WriteHeader(http.StatusForbidden)
	w.Write(nil)
	return
}
