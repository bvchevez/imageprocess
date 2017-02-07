package image

import (
	"fmt"
	"bytes"
	"image"

	"github.com/nfnt/resize"
	"github.com/bvchevez/imageprocess/point"
	"github.com/bvchevez/imageprocess/smartcrop"
)

// CropOperation reprensents the information necessary to perform a Image Crop.
type CropOperation struct {
	NewWidth  int64
	NewHeight int64
	Position  *point.Point
	Image     *MutableImage
}

// FindBestCrop calls smartcrop to analyze the image and returns top crop parameter
func (i *CropOperation) FindBestCrop() (smartcrop.Crop, error) {

	mutable := *i.Image
	s := bytes.NewReader(mutable.GetImage().Data)
	img, _, err := image.Decode(s)
	if err != nil {
		return smartcrop.Crop{}, fmt.Errorf(err.Error())
	}

	settings := smartcrop.CropSettings{
		FaceDetection: true,
		FaceDetectionHaarCascadeFilepath: appPath + HaarCascadeFrontalFaceAlt,
		InterpolationType: resize.Bicubic,
		//DebugMode: true,
	}
	analyzer := smartcrop.NewAnalyzerWithCropSettings(settings)
	cropWidth := i.NewWidth
	cropHeight := i.NewHeight

	// crop returns something like
	// crop: {X:98 Y:0 Width:882 Height:441 Score:{Detail:-2.4122274199066016 Saturation:21.35539732757885 Skin:935.1400903412627 Total:0.006516884013613292}}
	crop, err := analyzer.FindBestCrop2(img, int(cropWidth), int(cropHeight))
	return crop, err
}

// Do performs the actual crop operation
func (i *CropOperation) Do() error {
	var crop smartcrop.Crop

	img := *i.Image
	pos := i.Position

	// handle auto positioning
	if pos.X < 0 || pos.Y < 0 {
		var err error
		crop, err = i.FindBestCrop()
		if err != nil {
			return fmt.Errorf(err.Error())
		}

		width := img.GetImage().Width
		height := img.GetImage().Height

		// set auto crop coordinates
		if pos.X < 0 {
			pos.X = int64(crop.X) + (int64(crop.Width) - i.NewWidth) / 2
		}
		if pos.Y < 0 {
			pos.Y = int64(crop.Y) + (int64(crop.Height) - i.NewHeight) / 2
		}

		// fit crop window if flowing outside the image
		if pos.X > width - i.NewWidth {
			pos.X = width - i.NewWidth
		}
		if pos.Y > height - i.NewHeight {
			pos.Y = height - i.NewHeight
		}

		if pos.X < 0 {
			pos.X = 0
		}
		if pos.Y < 0 {
			pos.Y = 0
		}
	}

	return img.Crop(i)
}

// IsValid checks if crop is even necessary. (if crop size is the same or greater than image size, we return false.)
func (i *CropOperation) IsValid() bool {
	img := *i.Image

	//if we're cropping something that has 0 for with or height, we return false.
	if i.NewWidth == 0 || i.NewHeight == 0 {
		return false
	}

	//if we're cropping something that's bigger than the image in width or height
	//we return false, that's not valid.
	if i.NewWidth > img.GetImage().Width || i.NewHeight > img.GetImage().Height {
		return false
	}

	return true
}

func (i *CropOperation) String() string {
	return fmt.Sprint("Crop")
}
