package image

import (
	"io/ioutil"
	"testing"

	"github.com/bvchevez/imageprocess/helper"
	"github.com/bvchevez/imageprocess/point"
	"github.com/h2non/bimg"
	"github.com/stretchr/testify/assert"
)

//benchmark transforming.
//avg 1818424 ns/op (1.8 ms)
//go test -bench=Benchmark_JPEG_Resize_only
func Benchmark_JPEG_Resize_only(b *testing.B) {
	helper.Logging = false
	for i := 0; i < b.N; i++ {
		img := getMockImageJPEG()
		op := &ResizeOperation{
			NewWidth:  int64(200),
			NewHeight: int64(300),
			Image:     &img,
		}
		op.Do()
	}
}

//benchmark transforming.
//avg 10818948 ns/op (10.8 ms)
//go test -bench=Benchmark_JPEG_Resize_And_ApplyChange
func Benchmark_JPEG_Resize_And_ApplyChange(b *testing.B) {

	for i := 0; i < b.N; i++ {
		img := getMockImageJPEG()
		img.Resize(&ResizeOperation{
			NewWidth:  int64(200),
			NewHeight: int64(300),
			Image:     &img,
		})
		img.ApplyChanges()
	}
}

func getMockImageJPEG() MutableImage {
	data, _ := ioutil.ReadFile("test/test.jpg")
	img, _ := MakeImage(data, "1", "")
	img.SetDimensions()
	return img
}

func getMockImageTIFF() MutableImage {
	data, _ := ioutil.ReadFile("test/test.tif")
	img, _ := MakeImage(data, "1", "")
	img.SetDimensions()
	return img
}

func getMockImagePNG() MutableImage {
	data, _ := ioutil.ReadFile("test/test.png")
	img, _ := MakeImage(data, "1", "")
	img.SetDimensions()
	return img
}

func getMockImageCMYK() MutableImage {
	data, _ := ioutil.ReadFile("test/cmyk.jpg")
	img, _ := MakeImage(data, "1", "")
	img.SetDimensions()
	return img
}

//go test -run Test_ImageJPEG_ResizeCMYK -v
func Test_ImageJPEG_ResizeCMYK(t *testing.T) {
	img := getMockImageCMYK()

	err := img.Resize(&ResizeOperation{
		NewWidth:  int64(200),
		NewHeight: int64(300),
		Image:     &img,
	})
	if err != nil {
		t.Errorf("Error not expected! [%v]", err)
	}

	assert.Equal(t, int64(200), img.GetImage().Width)
	assert.Equal(t, int64(300), img.GetImage().Height)
	assert.Equal(t, int64(2000), img.GetImage().SourceWidth)
	assert.Equal(t, int64(3000), img.GetImage().SourceHeight)
}

//go test -run Test_ImageJPEG_SetDimension -v
func Test_ImageJPEG_SetDimension(t *testing.T) {
	data, _ := ioutil.ReadFile("test/test.jpg")
	img, _ := MakeImage(data, "1", "")
	err := img.SetDimensions()

	if err != nil {
		t.Errorf("Error not expected! [%v]", err)
	}

	assert.Equal(t, int64(375), img.GetImage().Width)
	assert.Equal(t, int64(500), img.GetImage().Height)
	assert.Equal(t, false, img.GetImage().Animated) // jpeg, not animated.
}

//go test -run Test_ImageTIFF_Resize_only -v
func Test_ImageTIFF_Resize_only(t *testing.T) {
	img := getMockImageTIFF()

	err := img.Resize(&ResizeOperation{
		NewWidth:  int64(50),
		NewHeight: int64(100),
		Image:     &img,
	})
	if err != nil {
		t.Errorf("Error not expected! [%v]", err)
	}

	//this returns jpeg by default.
	assert.Equal(t, "jpeg", bimg.NewImage(img.GetImage().Data).Type())
	assert.Equal(t, int64(50), img.GetImage().Width)
	assert.Equal(t, int64(100), img.GetImage().Height)
}

//go test -run Test_ImageTIFF_Resize_ApplyChange -v
func Test_ImageTIFF_Resize_ApplyChange(t *testing.T) {
	img := getMockImageTIFF()

	err := img.Resize(&ResizeOperation{
		NewWidth:  int64(50),
		NewHeight: int64(100),
		Image:     &img,
	})
	if err != nil {
		t.Errorf("Error not expected! [%v]", err)
	}

	assert.Equal(t, int64(50), img.GetImage().Width)
	assert.Equal(t, int64(100), img.GetImage().Height)
}

//go test -run Test_ImageJPEG_Resize_only -v
func Test_ImageJPEG_Resize_only(t *testing.T) {
	img := getMockImageJPEG()

	err := img.Resize(&ResizeOperation{
		NewWidth:  int64(200),
		NewHeight: int64(300),
		Image:     &img,
	})
	if err != nil {
		t.Errorf("Error not expected! [%v]", err)
	}

	assert.Equal(t, int64(200), img.GetImage().Width)
	assert.Equal(t, int64(300), img.GetImage().Height)
}

//go test -run Test_ImageJPEG_Resize_twoBadDimensions -v
func Test_ImageJPEG_Resize_largerDimensions(t *testing.T) {
	img := getMockImageJPEG()

	err := img.Resize(&ResizeOperation{
		NewWidth:  int64(4000),
		NewHeight: int64(4000),
		Image:     &img,
	})

	// assert error is nill.
	assert.Nil(t, err)
	assert.Equal(t, int64(375), img.GetImage().Width)
	assert.Equal(t, int64(500), img.GetImage().Height)
}

//go test -run Test_ImageJPEG_Crop -v
func Test_ImageJPEG_Crop(t *testing.T) {
	img := getMockImageJPEG()

	err := img.Crop(&CropOperation{
		Image:     &img,
		NewWidth:  int64(150),
		NewHeight: int64(200),
		Position: &point.Point{
			X: 1,
			Y: 2,
		},
	})

	if err != nil {
		t.Errorf("Error not expected! [%v]", err)
	}

	assert.Equal(t, int64(150), img.GetImage().Width)
	assert.Equal(t, int64(200), img.GetImage().Height)
}

//go test -run Test_ImageJPEG_Crop_ZeroWHDimensions -v
func Test_ImageJPEG_Crop_ZeroWHDimensions(t *testing.T) {
	img := getMockImageJPEG()

	err := img.Crop(&CropOperation{
		Image:     &img,
		NewWidth:  int64(0),
		NewHeight: int64(0),
		Position: &point.Point{
			X: 1,
			Y: 2,
		},
	})

	assert.Equal(t, "Crop [0x0] @ (1, 2) is out of bounds.", err.Error())
}

//go test -run Test_ImageJPEG_Crop_LargeWHDimensions -v
func Test_ImageJPEG_Crop_LargeWHDimensions(t *testing.T) {
	img := getMockImageJPEG()

	err := img.Crop(&CropOperation{
		Image:     &img,
		NewWidth:  int64(100000),
		NewHeight: int64(90000),
		Position: &point.Point{
			X: 1,
			Y: 2,
		},
	})

	assert.Equal(t, "Crop [100000x90000] @ (1, 2) is out of bounds.", err.Error())
}

//go test -run Test_ImageJPEG_Crop_CropOutOfBounds -v
func Test_ImageJPEG_Crop_CropOutOfBounds(t *testing.T) {
	img := getMockImageJPEG()

	err := img.Crop(&CropOperation{
		Image:     &img,
		NewWidth:  int64(375),
		NewHeight: int64(500),
		Position: &point.Point{
			X: 1,
			Y: 2,
		},
	})

	assert.Equal(t, "extract_area: bad extract area\n", err.Error())
}

//go test -run Test_ImageJPEG_Crop_CropDimSameAsRealDim -v
func Test_ImageJPEG_Crop_CropDimSameAsRealDim(t *testing.T) {
	img := getMockImageJPEG()

	err := img.Crop(&CropOperation{
		Image:     &img,
		NewWidth:  int64(375),
		NewHeight: int64(500),
		Position: &point.Point{
			X: 0,
			Y: 0,
		},
	})

	if err != nil {
		t.Errorf("Error not expected! [%v]", err)
	}

	// crop is one pixel smaller, because vips will break if they're exactly the same.
	assert.Equal(t, int64(375), img.GetImage().Width)
	assert.Equal(t, int64(500), img.GetImage().Height)
}

//go test -run Test_ImageJPEG_Crop_badPosition -v
func Test_ImageJPEG_Crop_badPosition(t *testing.T) {
	img := getMockImageJPEG()

	err := img.Crop(&CropOperation{
		Image:     &img,
		NewWidth:  int64(120),
		NewHeight: int64(200),
		Position: &point.Point{
			X: 3000, //these positions are out of bound.
			Y: 4000,
		},
	})

	assert.Equal(t, "extract_area: bad extract area\n", err.Error())
}

//go test -run Test_ApplyChanges_AutoConvertTiffToJPEG -v
func Test_ApplyChanges_AutoConvertTiffToJPEG(t *testing.T) {
	img := getMockImageTIFF()

	err := img.Resize(&ResizeOperation{
		Image:     &img,
		NewWidth:  int64(20),
		NewHeight: int64(30),
	})
	if err != nil {
		t.Errorf("Error not expected! [%v]", err)
	}

	img.ApplyChanges()
	assert.Equal(t, JPEG, GetFileType(img.GetImage().Data))
}

//go test -run Test_ApplyChanges_Density2 -v
func Test_ApplyChanges_Density2(t *testing.T) {
	img := getMockImageJPEG()

	assert.Equal(t, int64(375), img.GetImage().Width)
	assert.Equal(t, int64(500), img.GetImage().Height)

	err := img.Density(&DensityOperation{
		Image:      &img,
		NewDensity: int64(2),
	})
	if err != nil {
		t.Errorf("Error not expected! [%v]", err)
	}
	img.ApplyChanges()

	assert.Equal(t, int64(750), img.GetImage().Width)
	assert.Equal(t, int64(1000), img.GetImage().Height)
}

//go test -run Test_Quality_Default -v
func Test_Quality_Default(t *testing.T) {
	data, _ := ioutil.ReadFile("test/test.jpg")

	img := &ImageFixed{
		ImageData: &Image{
			Data:     data,
			Type:     GetFileType(data),
			Size:     int64(len(data)),
			Animated: false,
		},
		PipelineID: "1",
		Type:       GetFileType(data),
	}

	img.SetDimensions()

	//sets default to 70
	img.SetDefaults(Options{
		Quality: 70,
		Density: 2,
	})

	//expect 70 w/o doing anything else.
	assert.Equal(t, int64(70), img.NewQuality)

	img.NewQuality = 80
	assert.Equal(t, int64(80), img.NewQuality)
}

//go test -run Test_BicubicThreshold_Default -v
func Test_BicubicThreshold_Default(t *testing.T) {
	data, _ := ioutil.ReadFile("test/test.jpg")

	img := &ImageFixed{
		ImageData: &Image{
			Data:     data,
			Type:     GetFileType(data),
			Size:     int64(len(data)),
			Animated: false,
		},
		PipelineID: "1",
		Type:       GetFileType(data),
	}

	img.SetDimensions()

	img.SetDefaults(Options{
		BicubicThreshold: 250,
	})

	//expect 70 w/o doing anything else.
	assert.Equal(t, int64(250), img.BicubicThreshold)
}

//go test -run Test_Quality_Setting -v
func Test_Quality_Setting(t *testing.T) {
	img := getMockImageJPEG()
	prev := img.GetImage().Size
	err := img.Quality(&QualityOperation{
		Image:      &img,
		NewQuality: int64(2),
	})
	if err != nil {
		t.Errorf("Error not expected! [%v]", err)
	}
	img.ApplyChanges()
	current := img.GetImage().Size

	//size should decrease when quality is decreased.
	assert.True(t, prev > current)
}

//go test -run Test_ImageJPEG_Crop_PUX_5115 -v
func Test_ImageJPEG_Crop_PUX_5115(t *testing.T) {
	img := getMockImageJPEG()

	err := img.Crop(&CropOperation{
		Image:     &img,
		NewWidth:  int64(200),
		NewHeight: int64(200),
		Position: &point.Point{
			X: 0,
			Y: 0,
		},
	})

	if err != nil {
		t.Errorf("Error not expected! [%v]", err)
	}

	// crop is one pixel smaller, because vips will break if they're exactly the same.
	assert.Equal(t, int64(200), img.GetImage().Width)
	assert.Equal(t, int64(200), img.GetImage().Height)
}
