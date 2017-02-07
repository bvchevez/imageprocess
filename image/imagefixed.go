// imagefixed.go contains all logic necessary for the transformation of a fixed image.
package image

import (
	"fmt"
	"time"

	"github.com/bvchevez/imageprocess/helper"
	"github.com/h2non/bimg"
)

// ImageFixed is a struct that represents a non-animated image, such as jpeg/png/tiff/webm
type ImageFixed struct {
	PipelineID       string
	ImageData        *Image
	NewQuality       int64  // Quality to save this to.
	NewDensity       int64  // Final density for this image.
	Type             string // Image MIME
	BicubicThreshold int64  // Minimum pixels we want before converting to bicubic
}

// SetDimensions initializes ImageData with actual image width and height.
func (i *ImageFixed) SetDimensions() error {
	size, err := bimg.NewImage(i.ImageData.Data).Size()
	if err != nil {
		return err
	}

	i.ImageData.Width = int64(size.Width)
	i.ImageData.Height = int64(size.Height)
	i.ImageData.Size = int64(len(i.ImageData.Data))

	return nil
}

// set default data.
func (i *ImageFixed) SetDefaults(o Options) {
	i.NewQuality = o.Quality
	i.NewDensity = o.Density
	i.BicubicThreshold = o.BicubicThreshold
}

// ApplyChanges applies anything other than Resize or Crop (such as Density, Quality, colorspace... etc)
func (i *ImageFixed) ApplyChanges() error {
	defer helper.Timer(helper.TimerPayload{
		Start: time.Now(),
		Name:  "(" + i.PipelineID + ") " + i.Type + " ApplyChanges",
	})

	opt := bimg.Options{
		Quality:   int(i.NewQuality),
		Interlace: interlace,
	}

	if i.NewDensity == 2 {
		opt.Width = int(i.GetImage().Width * 2)
		opt.Height = int(i.GetImage().Height * 2)
	}

	//make sure this image is sRGB colorspace.
	opt.Interpretation = bimg.InterpretationSRGB

	imgByte, err := bimg.Resize(i.ImageData.Data, opt)
	if err != nil {
		return err
	}

	i.ImageData.Data = imgByte
	i.SetDimensions()
	return nil
}

// Resize takes in resize operation and performs resize on the image.
func (i *ImageFixed) Resize(o *ResizeOperation) error {
	defer helper.Timer(helper.TimerPayload{
		Start: time.Now(),
		Name:  "(" + i.PipelineID + ") " + i.Type + " resize",
	})

	// in this case, we simply ignore resize and return current dimensions.
	if o.IsValid() == false {
		return nil
	}

	opt := bimg.Options{
		Width:        int(o.NewWidth),
		Height:       int(o.NewHeight),
		Quality:      100,
		Force:        true,
		Interpolator: bimg.Bilinear,
	}

	if o.NewWidth <= i.BicubicThreshold {
		opt.Interpolator = bimg.Bicubic
	}

	newByte, err := bimg.Resize(i.ImageData.Data, opt)
	if err != nil {
		return err
	}

	i.ImageData.Data = newByte
	i.SetDimensions()

	return nil
}

// Crop takes in a crop operation object and performs the actual crop on the image.
func (i *ImageFixed) Crop(o *CropOperation) error {
	defer helper.Timer(helper.TimerPayload{
		Start: time.Now(),
		Name:  "(" + i.PipelineID + ") " + i.Type + " crop",
	})

	if o.IsValid() == false {
		return fmt.Errorf("Crop [%dx%d] @ (%d, %d) is out of bounds.", o.NewWidth, o.NewHeight, o.Position.X, o.Position.Y)
	}

	opt := bimg.Options{
		Top:        int(o.Position.Y),
		Left:       int(o.Position.X),
		AreaWidth:  int(o.NewWidth),
		AreaHeight: int(o.NewHeight),
		Quality:    100,
	}

	if opt.Top == 0 && opt.Left == 0 {
		opt.Top = -1
	}

	newByte, err := bimg.Resize(i.ImageData.Data, opt)
	if err != nil {
		return err
	}

	i.ImageData.Data = newByte
	i.SetDimensions()
	return nil
}

//Density sets the density for our image but doesn't actually apply the density.
func (i *ImageFixed) Density(o *DensityOperation) error {
	i.NewDensity = o.NewDensity
	return nil
}

//Quality sets the quality for our image but doesn't actually apply the quality.
func (i *ImageFixed) Quality(o *QualityOperation) error {
	i.NewQuality = o.NewQuality
	return nil
}

// GetImage returns the image data
func (i *ImageFixed) GetImage() *Image {
	return i.ImageData
}

// Shutdown represents all clean up actions necessary after transformation is complete.
func (i *ImageFixed) Shutdown() {}
