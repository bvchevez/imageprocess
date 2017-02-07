//imagegif.go is responsible for all transformations of a gif.
package image

import (
	"bytes"
	"fmt"
	"time"

	"image/gif"
	"os/exec"

	"github.com/bvchevez/imageprocess/helper"
	"github.com/bvchevez/imageprocess/point"
	log "github.com/Sirupsen/logrus"
)

// ImageGIF is a struct to hold gif data.
type ImageGIF struct {
	gifDecoded *gif.GIF // GifDecoded contains the result after running the bytes buffer through gif.DecodeAll(..).
	ImageData  *Image   // ImageData contains binary/dimensions of the gif.
	PipelineID string   // PipelineID stores the id of this call.

	QualityOp bool  // QualityOp triggers quality
	Colors    int64 // Color we want to display.

	ResizeOp     bool  // ResizeOp triggers resizing
	ResizeWidth  int64 // ResizeWidth stores the width we want resize to be
	ResizeHeight int64 // ResizeHeight stores the height we want resize to be.

	CropOp       bool         // CropOp triggers cropping
	CropWidth    int64        // CropWidth stores the width of the crop
	CropHeight   int64        // CropHeight stores the height of the crop
	CropPosition *point.Point // CropPosition stores the top left point of the new crop
}

// SetDimensions finds and sets the width/height of our image.
func (i *ImageGIF) SetDimensions() error {
	ff_bounds := i.gifDecoded.Image[0].Bounds()
	i.ImageData.Width = int64(ff_bounds.Dx())
	i.ImageData.Height = int64(ff_bounds.Dy())
	return nil
}

// set default data.
func (i *ImageGIF) SetDefaults(o Options) {}

// ApplyChanges applies the changes on the gif (currently acts as a stub to satisfy interface)
func (i *ImageGIF) ApplyChanges() error {
	var out bytes.Buffer

	args := []string{}

	// --crop=x1,y1+WxH
	if i.CropOp == true {
		args = append(
			args,
			fmt.Sprintf("--crop=%d,%d+%dx%d", i.CropPosition.X, i.CropPosition.Y, i.CropWidth, i.CropHeight),
		)
	}

	// --resize=WxH
	if i.ResizeOp == true {
		args = append(
			args,
			fmt.Sprintf("--resize=%dx%d", i.ResizeWidth, i.ResizeHeight),
		)
	}

	// --colors=2-256
	if i.QualityOp == true {
		args = append(
			args,
			fmt.Sprintf("--colors=%d", i.Colors),
		)
	}

	defer helper.Timer(helper.TimerPayload{
		Start: time.Now(),
		Name:  fmt.Sprintf("(%s) GIF ApplyChanges %s", i.PipelineID, args),
	})

	cmd := exec.Command(
		"gifsicle",
		args...,
	)

	cmd.Stdin = bytes.NewReader(i.ImageData.Data)
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Warn("Error when transforming gif.")
	}

	// Figure out the final size of this gif.
	i.ImageData.Data = out.Bytes()
	i.ImageData.Size = int64(len(i.ImageData.Data))
	i.ImageData.Type = GIF

	gifdec, err := gif.DecodeAll(bytes.NewBuffer(i.ImageData.Data))
	if err != nil {
		return fmt.Errorf("gif decode error [%s]", err)
	}
	ff_bounds := gifdec.Image[0].Bounds()
	i.ImageData.Width = int64(ff_bounds.Dx())
	i.ImageData.Height = int64(ff_bounds.Dy())

	return nil
}

// Resize takes in resize operation and performs resize on the image.
func (i *ImageGIF) Resize(o *ResizeOperation) error {
	defer helper.Timer(helper.TimerPayload{
		Start: time.Now(),
		Name:  "(" + i.PipelineID + ") GIF Resize",
	})

	// in this case, we simply ignore resize and return current dimensions.
	if o.IsValid() == false {
		return nil
	}

	i.ResizeOp = true
	i.ResizeWidth = o.NewWidth
	i.ResizeHeight = o.NewHeight

	return nil
}

// Crop takes in a crop operation object and performs the actual crop on the image.
func (i *ImageGIF) Crop(o *CropOperation) error {
	defer helper.Timer(helper.TimerPayload{
		Start: time.Now(),
		Name:  "(" + i.PipelineID + ") GIF Crop",
	})

	if o.IsValid() == false {
		return nil
	}

	i.CropOp = true
	i.CropWidth = o.NewWidth
	i.CropHeight = o.NewHeight
	i.CropPosition = o.Position

	return nil
}

// Density does not apply to gifs.
func (i *ImageGIF) Density(o *DensityOperation) error {
	return nil
}

// Quality determines the gif quality by the amount of colors its using.
func (i *ImageGIF) Quality(o *QualityOperation) error {
	i.QualityOp = true

	// Quality can be between 1 and 100. We need to conver thtat into colors
	// which can be anywhere between 2 and 256.
	i.Colors = int64((float64(o.NewQuality) * 2.56))

	return nil
}

// GetImage returns the image data
func (i *ImageGIF) GetImage() *Image {
	return i.ImageData
}

// Shutdown represents all clean up actions necessary after transformation is complete.
func (i *ImageGIF) Shutdown() {}
