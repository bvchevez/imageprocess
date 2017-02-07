package main

// pipeline.go handles the image transformations.
import (
	"fmt"
	"net/http"
	"time"

	"github.com/bvchevez/imageprocess/helper"
	"github.com/bvchevez/imageprocess/image"

	log "github.com/Sirupsen/logrus"
	"github.com/newrelic/go-agent"
)

// Pipeline represents a single pipeline to be used to process a single image
type Pipeline struct {
	id       string
	site     string
	path     string
	rawQuery string
	imgObj   image.MutableImage
}

// Process processes our pipeline
// Takes in the new relic transaction of this particular handler.
func (p *Pipeline) Process(txn newrelic.Transaction) (*image.Image, *Response) {
	p.id = helper.GetPipelineID(p.site, p.path, p.rawQuery)
	defer helper.Timer(helper.TimerPayload{
		Start: time.Now(),
		Name:  "(" + p.id + ") Total Time",
	})

	log.WithFields(log.Fields{
		"site":  p.site,
		"path":  p.path,
		"query": p.rawQuery,
	}).Info("(" + p.id + ") Processing Request")

	// Download image.
	// Returns 403 on failure.
	bytes, err := p.downloadImage(txn)
	if err != nil {
		return nil, &Response{
			Code:  http.StatusForbidden,
			Data:  [1]string{err.Error()},
			Image: nil,
		}
	}

	// Initializes image.
	// Returns 400 on failure.
	err = p.initImage(bytes, txn)
	if err != nil {
		return nil, &Response{
			Code:  http.StatusBadRequest,
			Data:  [1]string{err.Error()},
			Image: nil,
		}
	}

	// Make and validates operations.
	// Returns 400 on failure.
	ops, err := image.MakeOperations(p.rawQuery, p.imgObj)
	if err != nil {
		return nil, &Response{
			Code:  http.StatusBadRequest,
			Data:  [1]string{err.Error()},
			Image: nil,
		}
	}

	// Performs operations.
	// Returns 400 on failure.
	for _, op := range ops {
		err := p.doOp(txn, op)
		if err != nil {
			return nil, &Response{
				Code:  http.StatusBadRequest,
				Data:  [1]string{err.Error()},
				Image: nil,
			}
		}
	}

	return p.imgObj.GetImage(), nil
}

func (p *Pipeline) downloadImage(txn newrelic.Transaction) ([]byte, error) {
	defer newrelic.Segment{
		Name:      "Download Image",
		StartTime: newrelic.StartSegmentNow(txn),
	}.End()

	bytes, err := GetImage(p.site, p.path, p.id)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (p *Pipeline) initImage(bytes []byte, txn newrelic.Transaction) error {
	defer newrelic.Segment{
		Name:      "Initialize Image",
		StartTime: newrelic.StartSegmentNow(txn),
	}.End()
	var err error

	p.imgObj, err = image.MakeImage(bytes, p.id, p.rawQuery)
	if err != nil {
		return err
	}

	p.imgObj.SetDefaults(image.Options{
		Quality:          helper.String2Int64(*config.defaultQuality),
		BicubicThreshold: helper.String2Int64(*config.bicubicThreshold),
	})

	return nil
}

func (p *Pipeline) doOp(txn newrelic.Transaction, op image.Operations) error {
	defer newrelic.Segment{
		Name:      fmt.Sprintf("%s", op),
		StartTime: newrelic.StartSegmentNow(txn),
	}.End()

	return op.Do()
}
