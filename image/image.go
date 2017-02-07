// image.go is responsible for taking in an ImageOperation object and transforming an image.
package image

import (
	"bytes"
	"fmt"
	"strings"

	"image/gif"
	"image/jpeg"
)

// MutableImage represents the gif/jpeg/png...etc type that will be responsible for transforming the Image struct.
type MutableImage interface {
	SetDimensions() error              // SetDimensions initializes the image. Will load vips pointer and figure out the image size.
	SetDefaults(o Options)             // SetDefaults sets the default parameters.
	ApplyChanges() error               // ApplyChanges applies the changes we have been making from temp data back into original data.
	Crop(o *CropOperation) error       // Crop crops the image.
	Resize(o *ResizeOperation) error   // resizes the image
	Quality(o *QualityOperation) error // Quality sets quality (1-100)
	Density(o *DensityOperation) error // Density sets density (1,2).
	GetImage() *Image                  // GetImage returns the image binary data.
	Shutdown()                         // Shutdown shuts down the image, clears any memory ref.
}

// Image is a struct that holds the basic informations of any single image to be transformed upon.
type Image struct {
	Data         []byte // Image data got from S3/URL
	Type         string // Image type
	Animated     bool   // Is this animated
	Size         int64  // Image size
	Width        int64  // Image width
	Height       int64  // Image height
	SourceWidth  int64  // Source image width
	SourceHeight int64  // Source image height
}

func (i *Image) SetSourceDimensions() {
	i.SourceWidth = i.Width
	i.SourceHeight = i.Height
}

// MakeImage dynamically initializes an image object depends on image type.
// if firstFrame is true, we return a jpeg of the first frame for gif.
func MakeImage(data []byte, pipelineID, rawQuery string) (MutableImage, error) {

	if len(data) == 0 {
		return nil, fmt.Errorf("No image data retrieved.")
	}

	img := &Image{
		Data:     data,
		Type:     GetFileType(data),
		Size:     int64(len(data)),
		Animated: false,
	}

	var imageWrapper MutableImage

	switch img.Type {
	case GIF:
		img.Animated = true
		buf := bytes.NewBuffer(data)
		gifdec, err := gif.DecodeAll(buf)
		if err != nil {
			return nil, err
		}

		//if only first frame is requested, we convert first frame to jpeg and return *ImageJPEG
		if IsFirstFrame(rawQuery) {
			gifimg := gifdec.Image[0]
			b := new(bytes.Buffer)
			jpeg.Encode(b, gifimg, &jpeg.Options{Quality: 100})
			gifbytes := b.Bytes()

			jpegImg, err := MakeImage(gifbytes, pipelineID, rawQuery)
			if err != nil {
				return nil, err
			}

			return jpegImg, nil
		}

		imageWrapper = &ImageGIF{
			gifDecoded: gifdec,
			ImageData:  img,
			PipelineID: pipelineID,
		}

	case JPEG, PNG, TIFF:
		imageWrapper = &ImageFixed{
			ImageData:  img,
			PipelineID: pipelineID,
			Type:       img.Type,
		}

	default:
		return nil, fmt.Errorf("Invalid image type [%s]", img.Type)
	}

	imageWrapper.SetDimensions()
	imageWrapper.GetImage().SetSourceDimensions()
	return imageWrapper, nil
}

// GetFileType takes a slice of byte and attempts to determine its file type
// Will return "application/octet-stream" for any non-matching slices
func GetFileType(data []byte) string {
	for t, sig := range fileTypes {
		if bytes.Equal(data[:2], sig) {
			return t
		}
	}

	return "application/octet-stream"
}

// MakeOperations parses out a request's query parameters and attempts to convert them into valid Image operations
func MakeOperations(rawQuery string, imgObj MutableImage) ([]Operations, error) {
	//Break query string up
	bits := strings.SplitN(rawQuery, "&", -1)
	operations := make([]Operations, 0, maxOperations)

	// Check for one or more query arguments
	if rawQuery == "" || len(bits) == 0 {
		return operations, nil
	}

	// Make sure we don't have more than the allowed operations.
	if len(bits) > maxOperations {
		return nil, fmt.Errorf("too many operations [%v]", len(bits))
	}

	imgObj.SetDimensions()

	// set starting dimension.
	// NOTE: This width/height combo will change to reflect the new resulting width/height
	// of each operations.
	width := imgObj.GetImage().Width
	height := imgObj.GetImage().Height
	for _, bit := range bits {

		//initialize new operation.
		operation := ImageOperation{
			ImageWidth:  width,
			ImageHeight: height,
			Image:       &imgObj,
		}

		// Split query argument into key/value
		split := strings.Split(bit, "=")

		// Check for valid key/value pair and ensure the requested operation is allowed
		if len(split) != 2 {
			return nil, fmt.Errorf("invalid parameter [%v]", bit)
		}
		action := split[0]
		params := strings.Split(split[1], ";")

		newOp, err := operation.Make(params, action)
		if err != nil {
			return nil, err
		}

		if newOp == nil {
			continue
		}

		// if new width is zero, that means width isn't affected by this operation.
		if operation.NewWidth > 0 {
			width = operation.NewWidth
		}

		// if new height is zero, that means height isn't affected by this operation.
		if operation.NewHeight > 0 {
			height = operation.NewHeight
		}
		operations = append(operations, newOp)
	}

	//if we have more than one operation, we add an additional operation to apply changes.
	if len(operations) > 0 {
		operations = append(operations, &ApplyOperation{Image: &imgObj})
	}

	return operations, nil
}

// IsFirstFrame figures out if "frame=1" exists in our url parameter.
func IsFirstFrame(rawQuery string) bool {
	return strings.Contains(rawQuery, "frame=1")
}

// DoTransformation performs the transformation as defined by p.operations
func DoTransformation(o []Operations) error {
	if len(o) == 0 {
		return nil
	}

	for _, op := range o {
		err := op.Do()
		if err != nil {
			return err
		}
	}

	return nil
}
