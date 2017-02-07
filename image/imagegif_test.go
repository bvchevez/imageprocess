package image

import (
	"bytes"
	"reflect"
	"testing"

	"image/gif"
	"io/ioutil"

	"github.com/bvchevez/imageprocess/point"
	"github.com/stretchr/testify/assert"
)

func mockGif(rawQuery string) MutableImage {
	data, _ := ioutil.ReadFile("test/test.gif")
	img, _ := MakeImage(data, "1", rawQuery)
	return img
}

func getMockImageGIF() MutableImage {
	data, _ := ioutil.ReadFile("test/test.gif")
	img, _ := MakeImage(data, "1", "")
	img.SetDimensions()
	return img
}

func mockImageGIF() ImageGIF {
	data, _ := ioutil.ReadFile("test/test.gif")

	img := &Image{
		Data:     data,
		Type:     GetFileType(data),
		Size:     int64(len(data)),
		Animated: false,
	}

	img.Animated = true
	gifdec, _ := gif.DecodeAll(bytes.NewBuffer(data))

	return ImageGIF{
		gifDecoded: gifdec,
		ImageData:  img,
		PipelineID: "1",
	}
}

//benchmark transforming resize only.
//avg 591936243 ns/op (591 ms)
//go test -bench=Benchmark_GIF_Resize_only
func Benchmark_GIF_Resize_only(b *testing.B) {
	for i := 0; i < b.N; i++ {
		img := mockGif("")
		img.SetDimensions()
		img.Resize(&ResizeOperation{
			NewWidth:  200,
			NewHeight: 300,
			Image:     &img,
		})
		img.ApplyChanges()
	}
}

//benchmark transforming crop + resize.
//avg 570523548 ns/op (570 ms)
//go test -bench=Benchmark_GIF_Crop_Resize
func Benchmark_GIF_Crop_Resize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		img := mockGif("")
		img.SetDimensions()

		img.Crop(&CropOperation{
			Image:     &img,
			NewWidth:  800,
			NewHeight: 400,
			Position: &point.Point{
				X: 1,
				Y: 2,
			},
		})

		img.Resize(&ResizeOperation{
			NewWidth:  200,
			NewHeight: 300,
			Image:     &img,
		})

		img.ApplyChanges()
	}
}

//go test -run Test_ImageGIF_SetDimension_oneFrameOnly -v
func Test_ImageGIF_SetDimension_oneFrameOnly(t *testing.T) {
	t.Log("Test that when one frame is loaded, JPEG is what's returned.")
	data, _ := ioutil.ReadFile("test/test.gif")
	img, _ := MakeImage(data, "1", "frame=1")

	//newImg returned here should be JPEG, since we only need one frame.
	err := img.SetDimensions()
	if err != nil {
		t.Errorf("Error not expected %s", err.Error())
		return
	}

	switch jpegImg := img.(type) {
	case *ImageFixed:
		assert.Equal(t, "*image.ImageFixed", reflect.TypeOf(jpegImg).String())
		assert.Equal(t, int64(900), jpegImg.GetImage().Width)
		assert.Equal(t, int64(450), jpegImg.GetImage().Height)
		assert.Equal(t, false, jpegImg.GetImage().Animated) // we're only returning first frame. Not animated.
	default:
		t.Errorf("Other img types not expected.")
	}
}

//go test -run Test_ImageGIF_LoadFramesFromGif_loadAllFrames -v
func Test_ImageGIF_LoadFramesFromGif_loadAllFrames(t *testing.T) {
	t.Log("Test loading all frames \"LoadFramesFromGif\" will load all frames.")
	data, err := ioutil.ReadFile("test/test.gif")
	if err != nil {
		t.Errorf("Error not expected at read %s", err.Error())
		return
	}

	img, err := MakeImage(data, "1", "")
	if err != nil {
		t.Errorf("Error not expected at MakeImage %s", err.Error())
		return
	}

	err = img.SetDimensions()
	if err != nil {
		t.Errorf("Error not expected at SetDimensions %s", err.Error())
		return
	}

	assert.Equal(t, int64(900), img.GetImage().Width)
	assert.Equal(t, int64(450), img.GetImage().Height)
	assert.Equal(t, true, img.GetImage().Animated) // returning the full gif, animated.
}

//go test -run Test_ImageGIF_Resize_multiFrame -v
func Test_ImageGIF_Resize_multiFrame(t *testing.T) {
	img := mockGif("")
	err := img.SetDimensions()
	if err != nil {
		t.Errorf("Error not expected at SetDimensions %s", err.Error())
		return
	}

	img.Resize(&ResizeOperation{
		NewWidth:  200,
		NewHeight: 100,
		Image:     &img,
	})

	img.ApplyChanges()

	assert.Equal(t, int64(200), img.GetImage().Width)
	assert.Equal(t, int64(100), img.GetImage().Height)
	assert.Equal(t, int64(900), img.GetImage().SourceWidth)
	assert.Equal(t, int64(450), img.GetImage().SourceHeight)

	switch img.(type) {
	case *ImageGIF:
	default:
		t.Errorf("Other img types not expected.")
	}

}

//go test -run Test_ImageGIF_Resize_singleFrame -v
func Test_ImageGIF_Resize_singleFrame(t *testing.T) {
	img := mockGif("resize=100:200&frame=1")
	err := img.SetDimensions()

	if err != nil {
		t.Errorf("Error not expected at SetDimensions %s", err.Error())
		return
	}

	img.Resize(&ResizeOperation{
		NewWidth:  200,
		NewHeight: 100,
		Image:     &img,
	})

	assert.Equal(t, "*image.ImageFixed", reflect.TypeOf(img).String())
	assert.Equal(t, int64(200), img.GetImage().Width)
	assert.Equal(t, int64(100), img.GetImage().Height)

}

//go test -run Test_ImageGIF_Crop_multiFrame -v
func Test_ImageGIF_Crop_multiFrame(t *testing.T) {
	img := mockGif("")
	err := img.SetDimensions()
	if err != nil {
		t.Errorf("Error not expected at SetDimensions %s", err.Error())
		return
	}

	img.Crop(&CropOperation{
		Image:     &img,
		NewWidth:  200,
		NewHeight: 100,
		Position: &point.Point{
			X: 1,
			Y: 2,
		},
	})

	img.ApplyChanges()

	assert.Equal(t, int64(200), img.GetImage().Width)
	assert.Equal(t, int64(100), img.GetImage().Height)

	switch img.(type) {
	case *ImageGIF:
	default:
		t.Errorf("Other img types not expected.")
	}

}

//go test -run Test_ImageGIF_Crop_singleFrame -v
func Test_ImageGIF_Crop_singleFrame(t *testing.T) {
	img := mockGif("resize=100:200&frame=1")
	err := img.SetDimensions()
	if err != nil {
		t.Errorf("Error not expected at SetDimensions %s", err.Error())
		return
	}

	img.Crop(&CropOperation{
		Image:     &img,
		NewWidth:  200,
		NewHeight: 100,
		Position: &point.Point{
			X: 1,
			Y: 2,
		},
	})

	// crop is one pixel smaller, to follow fixed image's logic.
	assert.Equal(t, "*image.ImageFixed", reflect.TypeOf(img).String())
	assert.Equal(t, int64(200), img.GetImage().Width)
	assert.Equal(t, int64(100), img.GetImage().Height)
}

//go test -run Test_ImageGIF_Crop_Resize_multiFrame -v
func Test_ImageGIF_Crop_Resize_multiFrame(t *testing.T) {
	img := mockGif("")
	err := img.SetDimensions()
	if err != nil {
		t.Errorf("Error not expected at SetDimensions %s", err.Error())
		return
	}

	img.Crop(&CropOperation{
		Image:     &img,
		NewWidth:  200,
		NewHeight: 100,
		Position: &point.Point{
			X: 1,
			Y: 2,
		},
	})

	img.Resize(&ResizeOperation{
		NewWidth:  50,
		NewHeight: 50,
		Image:     &img,
	})

	img.ApplyChanges()

	assert.Equal(t, int64(50), img.GetImage().Width)
	assert.Equal(t, int64(50), img.GetImage().Height)

	switch img.(type) {
	case *ImageGIF:
	default:
		t.Errorf("Other img types not expected.")
	}
}

//go test -run Test_ImageGIF_Resize_Crop_multiFrame -v
func Test_ImageGIF_Resize_Crop_multiFrame(t *testing.T) {
	img := mockGif("")
	err := img.SetDimensions()
	if err != nil {
		t.Errorf("Error not expected at SetDimensions %s", err.Error())
		return
	}

	img.Resize(&ResizeOperation{
		NewWidth:  50,
		NewHeight: 50,
		Image:     &img,
	})

	img.Crop(&CropOperation{
		Image:     &img,
		NewWidth:  200,
		NewHeight: 100,
		Position: &point.Point{
			X: 1,
			Y: 2,
		},
	})

	img.ApplyChanges()

	assert.Equal(t, int64(50), img.GetImage().Width)
	assert.Equal(t, int64(50), img.GetImage().Height)

	switch img.(type) {
	case *ImageGIF:
	default:
		t.Errorf("Other img types not expected.")
	}
}

//go test -run Test_ImageGIF_Crop_CropDimSameAsRealDim -v
func Test_ImageGIF_Crop_CropDimSameAsRealDim(t *testing.T) {
	img := mockGif("")
	img.SetDimensions()

	err := img.Crop(&CropOperation{
		Image:     &img,
		NewWidth:  int64(900),
		NewHeight: int64(450),
		Position: &point.Point{
			X: 0,
			Y: 0,
		},
	})

	if err != nil {
		t.Errorf("Error not expected! [%v]", err)
	}

	assert.Equal(t, int64(900), img.GetImage().Width)
	assert.Equal(t, int64(450), img.GetImage().Height)

	img.ApplyChanges()
}

//go test -run Test_ImageGIF_ApplyChanges_multiFrame -v
func Test_ImageGIF_ApplyChanges_multiFrame(t *testing.T) {
	img := mockGif("")
	err := img.SetDimensions()
	if err != nil {
		t.Errorf("Error not expected at SetDimensions %s", err.Error())
		return
	}

	img.Resize(&ResizeOperation{
		Image:     &img,
		NewWidth:  200,
		NewHeight: 100,
	})
	img.ApplyChanges()
	assert.Equal(t, int64(200), img.GetImage().Width)
	assert.Equal(t, int64(100), img.GetImage().Height)

	switch img.(type) {
	case *ImageGIF:
	default:
		t.Errorf("Other img types not expected.")
	}
}

//go test -run Test_ImageGIF_ApplyChanges_singleFrame -v
func Test_ImageGIF_ApplyChanges_singleFrame(t *testing.T) {
	img := mockGif("frame=1")
	err := img.SetDimensions()
	if err != nil {
		t.Errorf("Error not expected at SetDimensions %s", err.Error())
		return
	}

	img.Resize(&ResizeOperation{
		Image:     &img,
		NewWidth:  200,
		NewHeight: 100,
	})
	img.ApplyChanges()

	assert.Equal(t, "*image.ImageFixed", reflect.TypeOf(img).String())
	assert.Equal(t, int64(200), img.GetImage().Width)
	assert.Equal(t, int64(100), img.GetImage().Height)
}

//go test -run Test_ImageGIF_Resize_InvalidInput -v
func Test_ImageGIF_Resize_InvalidInput(t *testing.T) {
	img := mockGif("")
	err := img.SetDimensions()
	if err != nil {
		t.Errorf("Error not expected at SetDimensions %s", err.Error())
		return
	}

	err = img.Resize(&ResizeOperation{
		Image:     &img,
		NewWidth:  100000,
		NewHeight: 100000,
	})
	img.ApplyChanges()

	//error
	assert.Nil(t, err)
	assert.Equal(t, int64(900), img.GetImage().Width)
	assert.Equal(t, int64(450), img.GetImage().Height)
}

//go test -run Test_ImageGIF_Quality_ValidateSize -v
func Test_ImageGIF_Quality_ValidateSize(t *testing.T) {
	img := mockGif("")
	err := img.SetDimensions()
	if err != nil {
		t.Errorf("Error not expected at SetDimensions %s", err.Error())
		return
	}

	originalSize := img.GetImage().Size
	err = img.Quality(&QualityOperation{
		Image:      &img,
		NewQuality: 10,
	})
	img.ApplyChanges()

	newSize := img.GetImage().Size
	// validate that the new size gets small when we lower the quality.
	assert.True(t, (originalSize > newSize))
}

//go test -run Test_ImageGIF_Quality_ValidInput -v
func Test_ImageGIF_Quality_ValidInput(t *testing.T) {
	gif := mockImageGIF()

	gif.Quality(&QualityOperation{NewQuality: 100})
	assert.Equal(t, int64(256), gif.Colors)

	gif.Quality(&QualityOperation{NewQuality: 50})
	assert.Equal(t, int64(128), gif.Colors)

	gif.Quality(&QualityOperation{NewQuality: 1})
	assert.Equal(t, int64(2), gif.Colors)
}

//go test -run Test_ImageGIF_Resize_HandleErrorFromGifsicle -v
func Test_ImageGIF_Resize_HandleErrorFromGifsicle(t *testing.T) {
	img := mockGif("")
	img.SetDimensions()

	img.Resize(&ResizeOperation{
		Image:     &img,
		NewWidth:  0,
		NewHeight: 0,
	})

	err := img.ApplyChanges()
	assert.Nil(t, err) // no longer want to error out on gifsicle errors or warnings.
}

//go test -run Test_ImageGIF_Crop_InvalidInput -v
func Test_ImageGIF_Crop_InvalidInput(t *testing.T) {
	img := mockGif("")
	img.SetDimensions()

	img.Crop(&CropOperation{
		Image:     &img,
		NewWidth:  0,
		NewHeight: 0,
	})
	err := img.ApplyChanges()

	//error
	assert.Equal(t, nil, err)
	//no changes expected.
	assert.Equal(t, int64(900), img.GetImage().Width)
	assert.Equal(t, int64(450), img.GetImage().Height)
}
