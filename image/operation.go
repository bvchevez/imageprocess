// operation.go is responsible for translating an array of dimension string into the ImageOperation object.
package image

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bvchevez/imageprocess/helper"
	"github.com/bvchevez/imageprocess/point"
)

// Operations represents actual operations and operation information to be performed on image.
type Operations interface {
	Do() error
	IsValid() bool
}

// ImageOperation represents a single image command recieved by the pipeline.
type ImageOperation struct {
	ImageWidth  int64 // actual image width (before operation)
	ImageHeight int64 // actual image height (before operation)
	NewWidth    int64 // new width for resize/crop
	NewHeight   int64 // new height for resize/crop
	NewQuality  int64 // Quality of the outputted image (defaults to 75)
	NewDensity  int64 // Pixel density of the image after processing (defaults to 1)
	NewFrame    bool  // If true, we load first frame only (only valid for gifs)

	Position *point.Point //(x, y) are coordinates representing bottom left corner of our rectangle.

	Image *MutableImage
}

// Make turns parameters into operation types
func (i *ImageOperation) Make(params []string, action string) (Operations, error) {
	var err error

	if len(params) == 0 {
		return nil, fmt.Errorf("at least one dimension is required")
	}

	switch action {
	case "crop":
		if err = i.setCrop(params); err != nil {
			return nil, err
		}

		return &CropOperation{
			NewWidth:  i.NewWidth,
			NewHeight: i.NewHeight,
			Position:  i.Position,
			Image:     i.Image,
		}, nil

	case "fill":
		if err = i.setFill(params); err != nil {
			return nil, err
		}

		return &CropOperation{
			NewWidth:  i.NewWidth,
			NewHeight: i.NewHeight,
			Position:  i.Position,
			Image:     i.Image,
		}, nil

	case "resize":
		if err = i.setResize(params); err != nil {
			return nil, err
		}

		return &ResizeOperation{
			NewWidth:  i.NewWidth,
			NewHeight: i.NewHeight,
			Image:     i.Image,
		}, nil

	case "output-quality":
		if err = i.setOutputQuality(params); err != nil {
			return nil, err
		}

		return &QualityOperation{
			NewQuality: i.NewQuality,
			Image:      i.Image,
		}, nil

	case "density":
		if err = i.setDensity(params); err != nil {
			return nil, err
		}

		return &DensityOperation{
			NewDensity: i.NewDensity,
			Image:      i.Image,
		}, nil

	case "frame":
		return nil, nil
	}

	return nil, fmt.Errorf("invalid operation %v", action)
}

// setOutputQuality sets quality, must be of numeric type.
func (i *ImageOperation) setOutputQuality(dimensions []string) error {
	if len(dimensions) != 1 {
		return fmt.Errorf("too many dimensions. Maximum number of dimensions for quality is 1")
	}

	quality := dimensions[0]
	if helper.IsNumeric(quality) == false {
		return fmt.Errorf("invalid quality [%v]", quality)
	}

	i.NewQuality = helper.String2Int64(quality)
	if i.NewQuality < 0 || i.NewQuality > 100 {
		i.NewQuality = 0
		return fmt.Errorf("invalid quality [%v]", quality)
	}

	return nil
}

// setFrame sets frame attribute.
func (i *ImageOperation) setFrame(dimensions []string) error {
	if len(dimensions) != 1 {
		return fmt.Errorf("too many dimensions. Maximum number of dimensions for frame is 1")
	}

	frame := dimensions[0]
	if frame == "1" {
		i.NewFrame = true
		return nil
	}

	i.NewFrame = false
	return nil

}

// setDensity sets density, must be of numeric type.
func (i *ImageOperation) setDensity(dimensions []string) error {
	if len(dimensions) != 1 {
		return fmt.Errorf("too many dimensions. Maximum number of dimensions for density is 1")
	}

	density := dimensions[0]
	if density != "1" && density != "2" {
		return fmt.Errorf("invalid density [%v]", density)
	}

	i.NewDensity = helper.String2Int64(density)
	return nil

}

// setResize sets all parameters necessary for resize.
// dimensions looks like this
//  {"400:300"}
//  {"400:*"} (where * is wildcard, such that should be maintained.)
func (i *ImageOperation) setResize(params []string) error {
	if len(params) != 1 {
		return fmt.Errorf("resize must have dimensions parameter")
	}

	dimensions := strings.Split(params[0], ":")
	if len(dimensions) != 2 {
		return fmt.Errorf("resize must have two dimensions")
	}
	inputWidth := dimensions[0]
	inputHeight := dimensions[1]

	//calculate the dimensions of our new resize.
	if err := i.setNewDimensions(inputWidth, inputHeight); err != nil {
		return err
	}

	return nil
}

// setCrop sets all parameters necessary for cropping.
// dimensions looks like this
// must always have 2 dimensions and 2 custom parts.
//  {"0.6xw:0.7xh", "center,center"} //with and height ratio.
//  {"400:300", "center,center"}
//  {"400:*", "center,center"} (where * is wildcard, such that should be maintained.)
func (i *ImageOperation) setCrop(params []string) error {
	if len(params) > 2 {
		return fmt.Errorf("too many parameters for crop")
	}

	dimensions := strings.Split(params[0], ":")
	if len(dimensions) != 2 {
		return fmt.Errorf("crop size must have two dimensions")
	}
	inputWidth := dimensions[0]
	inputHeight := dimensions[1]

	coords := []string{"center", "top"}
	if len(params) == 2 {
		coords = strings.Split(params[1], ",")
	}
	if len(coords) != 2 {
		return fmt.Errorf("crop position must have two coordinates")
	}
	xPosition := coords[0]
	yPosition := coords[1]

	if err := i.setNewDimensions(inputWidth, inputHeight); err != nil {
		return err
	}

	if err := i.setCropPosition(xPosition, yPosition); err != nil {
		return err
	}

	return nil
}

// setFill sets all parameters necessary to crop to an aspect ratio.
// dimensions looks like this
// must have 2 relative dimensions, and optionally 2 position parts, like:
//  {"16:9", "center,top"}
//  {"4:3"}
func (i *ImageOperation) setFill(params []string) error {
	if len(params) > 2 {
		return fmt.Errorf("too many parameters for fill")
	}

	dimensions := strings.Split(params[0], ":")
	if len(dimensions) != 2 {
		return fmt.Errorf("aspect ratio to fill must have two dimensions")
	}
	aspectWidth := dimensions[0]
	aspectHeight := dimensions[1]

	coords := []string{"center", "top"}
	if len(params) == 2 {
		coords = strings.Split(params[1], ",")
	}
	if len(coords) != 2 {
		return fmt.Errorf("fill position must have two coordinates")
	}
	xPosition := coords[0]
	yPosition := coords[1]

	if err := i.setFillDimensions(aspectWidth, aspectHeight); err != nil {
		return err
	}

	if err := i.setCropPosition(xPosition, yPosition); err != nil {
		return err
	}

	return nil
}

// setNewDimensions calculates our new dimention based on inputWidth or heights.
// it will automatically process ratio and wildcard and set the newWidth and newHeight in pixels.
// inputWidth/inputHeight can be any mixture of the following...
//  1. integer  300, 450... etc
//  2. ratio    1.5xw, 2.3xh...
//  3. wildcard *
//  4. CAN NOT be BOTH wildcards (*)
func (i *ImageOperation) setNewDimensions(inputWidth, inputHeight string) error {
	if err := i.dims2px(inputWidth, inputHeight); err != nil {
		return err
	}

	return nil
}

// setFillDimensions takes two strings representing a relative width and height, converts them
// to integers and then to an aspect ratio, and calls setNewDimensions with crop params that will
// yield the largest possible crop of that aspect ratio that fits in the source image.
func (i *ImageOperation) setFillDimensions(aspectWidth, aspectHeight string) error {
	var width, height int64
	var aspectRatio, imageRatio float64
	var err error

	if width, err = strconv.ParseInt(aspectWidth, 10, 64); err != nil {
		return fmt.Errorf("fill aspect width must be integer, not %s", aspectWidth)
	}
	if height, err = strconv.ParseInt(aspectHeight, 10, 64); err != nil {
		return fmt.Errorf("fill aspect height must be integer, not %s", aspectHeight)
	}
	if aspectRatio, err = ratio(width, height); err != nil {
		return err
	}

	if imageRatio, err = ratio(i.ImageWidth, i.ImageHeight); err != nil {
		return err
	}

	// Rounding up can potentially return a value a few pixels larger than the actual image size.
	// Therefore we're rounding down to be safe, but upping the number of significant digits to
	// increase accuracy.
	if aspectRatio > imageRatio {
		return i.setNewDimensions("1xw", fmt.Sprintf("%fxw", helper.RoundDown(1/aspectRatio, 4)))
	}
	return i.setNewDimensions(fmt.Sprintf("%fxh", helper.RoundDown(aspectRatio, 4)), "1xh")
}

// setCropPosition takes two strings representing a horizontal and vertical position within the
// source image from which a crop should be taken, and sets the position in pixels of the left
// and top of the crop area, either by recognizing the special terms "left", "center", "top", &c.,
// or by calling pos2px to handle numeric or ratio formats.
func (i *ImageOperation) setCropPosition(xPos, yPos string) error {
	p := &point.Point{X: 0, Y: 0}

	switch {
	case xPos == "left":
		p.X = 0
	case xPos == "center":
		p.X = (i.ImageWidth - i.NewWidth) / 2
	case xPos == "right":
		p.X = i.ImageWidth - i.NewWidth
	case xPos == "auto":
		// deferred below

	default:
		px, err := i.pos2px(xPos)
		if err != nil {
			return fmt.Errorf(
				"crop X position not 'left', 'center', 'right' or 'auto', and %s", err)
		}
		if px > i.ImageWidth {
			return fmt.Errorf("crop X position is outside image: %d", px)
		}
		p.X = px
	}

	switch {
	case yPos == "top":
		p.Y = 0
	case yPos == "center":
		p.Y = (i.ImageHeight - i.NewHeight) / 2
	case yPos == "bottom":
		p.Y = i.ImageHeight - i.NewHeight
	case yPos == "auto":
		// deferred below

	default:
		px, err := i.pos2px(yPos)
		if err != nil {
			return fmt.Errorf(
				"crop Y position not 'top', 'center', 'bottom' or 'auto', and %s", err)
		}
		if px > i.ImageHeight {
			return fmt.Errorf("crop Y position is outside image: %d", px)
		}
		p.Y = px
	}

	if p.X < 0 {
		p.X = 0
	}
	if p.Y < 0 {
		p.Y = 0
	}

	// signal auto positioning
	if xPos == "auto" {
		p.X = -1
	}
	if yPos == "auto" {
		p.Y = -1
	}

	i.Position = p
	return nil
}

// dims2px converts both dimensions in place to pixels, if they're in valid forms of
// wildcard, ratio, or pixels.
// the only restriction is that we cannot have two wildcards.
func (i *ImageOperation) dims2px(inputWidth, inputHeight string) error {
	var pixelWidth, pixelHeight int64
	var widthIsWildcard, heightIsWildcard bool
	var err error

	if pixelWidth, widthIsWildcard, err = i.dim2px(inputWidth); err != nil {
		return err
	}
	i.NewWidth = pixelWidth
	if pixelHeight, heightIsWildcard, err = i.dim2px(inputHeight); err != nil {
		return err
	}
	i.NewHeight = pixelHeight

	if widthIsWildcard && heightIsWildcard {
		return fmt.Errorf("height and width cannot both be '*'")
	}

	if widthIsWildcard {
		heightRatio, err := ratio(i.NewHeight, i.ImageHeight)
		if err != nil {
			return err
		}
		i.NewWidth = int64(float64(i.ImageWidth) * heightRatio)
	}
	if heightIsWildcard {
		widthRatio, err := ratio(i.NewWidth, i.ImageWidth)
		if err != nil {
			return err
		}
		i.NewHeight = int64(float64(i.ImageHeight) * widthRatio)
	}

	return nil
}

// dim2px takes a string representing either dimension of an area within an image, and returns
// either an integer pixel coordinate or a flag indicating that a "*" wildcard was passed in.
func (i *ImageOperation) dim2px(dim string) (int64, bool, error) {
	if dim == "*" {
		return 0, true, nil
	} else if len(dim) > 2 && dim[len(dim)-2] == 'x' { // e.g. "1.123xh"
		px, err := i.ratio2px(dim)
		if err != nil {
			return 0, false, err
		}
		return px, false, nil
	} else if intDim, err := strconv.ParseInt(dim, 10, 64); err == nil {
		return intDim, false, nil
	}
	return 0, false, fmt.Errorf("dimension must be ratio, number, or '*', not '%s'", dim)
}

// pos2px takes a string representing either dimension of an area within an image, and returns
// an integer pixel coordinate.
func (i *ImageOperation) pos2px(dim string) (int64, error) {
	if len(dim) > 2 && dim[len(dim)-2] == 'x' { // e.g. "1.123xh"
		px, err := i.ratio2px(dim)
		if err != nil {
			return 0, err
		}
		return px, nil
	} else if intDim, err := strconv.ParseInt(dim, 10, 64); err == nil {
		return intDim, nil
	}
	return 0, fmt.Errorf("explicit position must be ratio or number, not '%s'", dim)
}

// ratio2px takes a string representing a decimal fraction of one dimension, like ".5xh", and
// returns the integer pixel length of that fraction of that side of the source image.
func (i *ImageOperation) ratio2px(dim string) (int64, error) {
	ratio, err := strconv.ParseFloat(dim[0:len(dim)-2], 64)
	if err != nil || ratio == 0.0 {
		return 0, fmt.Errorf("ratio must be non-zero float, not '%s'", dim)
	}

	length, err := i.getLength(dim[len(dim)-1])
	if err != nil {
		return 0, err
	}

	return int64(ratio * float64(length)), nil
}

// getLength takes a single 1-byte character indicating one dimension of an image and returns the
// length in pixels of that side of the source image. The character can be 'w' or 'h' to indicate
// the width or height directly, or 'g' or 'l' to indicate whichever of those is greater or lesser.
func (i *ImageOperation) getLength(dir byte) (int64, error) {
	switch {
	case dir == 'l' && i.ImageWidth < i.ImageHeight, dir == 'g' && i.ImageWidth > i.ImageHeight:
		return i.ImageWidth, nil
	case dir == 'w':
		return i.ImageWidth, nil
	case dir == 'h', dir == 'l', dir == 'g':
		return i.ImageHeight, nil
	default:
		return 0, fmt.Errorf("dimension must be 'w', 'h', 'g', or 'l', not '%c'", dir)
	}
}

// ratio takes two integers representing the numerator and denominator of a ratio and returns the
// float representation of that ratio.
func ratio(num, denom int64) (float64, error) {
	if denom == 0 {
		return 0.0, fmt.Errorf("denominator of ratio can't be zero")
	}
	return float64(num) / float64(denom), nil
}
