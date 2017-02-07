package image

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go test -run Test_Image_MakeImage_noData -v
func Test_Image_MakeImage_noData(t *testing.T) {
	var data []byte
	data = nil
	img, err := MakeImage(data, "1", "frame=1")

	if err == nil {
		t.Errorf("Error expected")
		return
	}

	assert.Equal(t, nil, img)
	assert.Equal(t, "No image data retrieved.", err.Error())
}

//go test -run Test_Image_MakeImage_InvalidImageType -v
func Test_Image_MakeImage_InvalidImageType(t *testing.T) {
	data, _ := ioutil.ReadFile("test/test.webp")
	img, err := MakeImage(data, "", "")

	if err == nil {
		t.Errorf("Error expected")
		return
	}

	assert.Equal(t, nil, img)
	assert.Equal(t, "Invalid image type [application/octet-stream]", err.Error())
}

//go test -run Test_Image_MakeOperation_noQuery -v
func Test_Image_MakeOperation_noQuery(t *testing.T) {
	img := getMockImageJPEG()
	op, err := MakeOperations("", img)

	expected := make([]Operations, 0, 5)
	assert.Equal(t, expected, op)
	assert.Equal(t, nil, err)
}

//go test -run Test_Image_MakeOperation_TooManyOp -v
func Test_Image_MakeOperation_TooManyOp(t *testing.T) {
	img := getMockImageJPEG()
	op, err := MakeOperations(
		"resize=100:*&resize=100:*&resize=100:*&resize=100:*&resize=100:*&resize=100:*&resize=100:*&resize=100:*",
		img)

	expected := []Operations(nil)
	assert.Equal(t, expected, op)
	assert.Equal(t, "too many operations [8]", err.Error())
}

//go test -run Test_Image_MakeOperation_BadSplit -v
func Test_Image_MakeOperation_BadSplit(t *testing.T) {
	img := getMockImageJPEG()
	op, err := MakeOperations(
		"resize=100=bad:*",
		img)

	expected := []Operations(nil)
	assert.Equal(t, expected, op)
	assert.Equal(t, "invalid parameter [resize=100=bad:*]", err.Error())
}

//go test -run Test_Image_MakeOperation_InvalidOperation -v
func Test_Image_MakeOperation_InvalidOperation(t *testing.T) {
	img := getMockImageJPEG()
	op, err := MakeOperations(
		"bad=9000:9000",
		img)

	expected := []Operations(nil)
	assert.Equal(t, expected, op)
	assert.Equal(t, "invalid operation bad", err.Error())
}

//go test -run Test_Image_MakeOperation_SkippedValidOperation -v
func Test_Image_MakeOperation_SkippedValidOperation(t *testing.T) {
	img := getMockImageJPEG()

	//frame will get skipped, since it does not need to be processed here.
	op, err := MakeOperations("frame=1&resize=200:*", img)

	//verify the first operation is a resize operation and not frame.
	assert.Equal(t, "*image.ResizeOperation", reflect.TypeOf(op[0]).String())
	assert.Equal(t, "*image.ApplyOperation", reflect.TypeOf(op[1]).String())
	assert.Equal(t, nil, err)
}

//go test -run Test_Image_MakeOperation_NormalCropResize -v
func Test_Image_MakeOperation_NormalCropResize(t *testing.T) {
	img := getMockImageJPEG()
	op, err := MakeOperations(
		"crop=200:200;0,0&resize=200:*",
		img)

	assert.Equal(t, "*image.CropOperation", reflect.TypeOf(op[0]).String())
	assert.Equal(t, "*image.ResizeOperation", reflect.TypeOf(op[1]).String())
	assert.Equal(t, nil, err)
}

//go test -run Test_Image_DoTransformation_NoOperations -v
func Test_Image_DoTransformation_NoOperations(t *testing.T) {
	op := []Operations(nil)
	err := DoTransformation(op)
	assert.Equal(t, nil, err)
}

//go test -run Test_Image_DoTransformation_ErrorDuringTransform -v
func Test_Image_DoTransformation_ErrorDuringTransform(t *testing.T) {
	img := getMockImageJPEG()
	op, _ := MakeOperations(
		"crop=200:200;200,0",
		img)
	err := DoTransformation(op)
	assert.Equal(t, int64(375), img.GetImage().Width)
	assert.Equal(t, int64(500), img.GetImage().Height)
	assert.Equal(t, "extract_area: bad extract area\n", err.Error())
}

//go test -run Test_Image_DoTransformation_normalTransform -v
func Test_Image_DoTransformation_normalTransform(t *testing.T) {
	img := getMockImageJPEG()
	op, _ := MakeOperations(
		"crop=200:200;0,0&resize=200:100",
		img)
	err := DoTransformation(op)

	assert.Equal(t, int64(200), img.GetImage().Width)
	assert.Equal(t, int64(100), img.GetImage().Height)
	assert.Equal(t, nil, err)
}

//go test -run Test_Image_DoTransformation_normalTransform_NoUpscaleJPEG -v
func Test_Image_DoTransformation_normalTransform_NoUpscaleJPEG(t *testing.T) {
	img := getMockImageJPEG()

	//initial image dimensions.
	assert.Equal(t, int64(375), img.GetImage().Width)
	assert.Equal(t, int64(500), img.GetImage().Height)

	//upscale image.
	op, _ := MakeOperations(
		"resize=1000:*",
		img)
	err := DoTransformation(op)
	assert.Equal(t, nil, err)

	// verify the result image dimensions does not change.
	imgScaled, _ := MakeImage(img.GetImage().Data, "1", "")
	imgScaled.SetDimensions()
	assert.Equal(t, int64(375), imgScaled.GetImage().Width)
	assert.Equal(t, int64(500), imgScaled.GetImage().Height)
}

//go test -run Test_Image_DoTransformation_normalTransform_NoUpscalePNG -v
func Test_Image_DoTransformation_normalTransform_NoUpscalePNG(t *testing.T) {
	img := getMockImagePNG()

	//initial image dimensions.
	assert.Equal(t, int64(172), img.GetImage().Width)
	assert.Equal(t, int64(200), img.GetImage().Height)

	//upscale image.
	op, _ := MakeOperations(
		"resize=1000:*",
		img)
	err := DoTransformation(op)
	assert.Equal(t, nil, err)

	// verify the result image dimensions does not change.
	imgScaled, _ := MakeImage(img.GetImage().Data, "1", "")
	imgScaled.SetDimensions()
	assert.Equal(t, int64(172), imgScaled.GetImage().Width)
	assert.Equal(t, int64(200), imgScaled.GetImage().Height)
}

//go test -run Test_Image_DoTransformation_normalTransform_NoUpscaleGIF -v
func Test_Image_DoTransformation_normalTransform_NoUpscaleGIF(t *testing.T) {
	img := getMockImageGIF()

	//initial image dimensions.
	assert.Equal(t, int64(900), img.GetImage().Width)
	assert.Equal(t, int64(450), img.GetImage().Height)

	//upscale image.
	op, _ := MakeOperations(
		"resize=1000:*",
		img)
	err := DoTransformation(op)
	assert.Equal(t, nil, err)

	// verify the result image dimensions does not change.
	imgScaled, _ := MakeImage(img.GetImage().Data, "1", "")
	imgScaled.SetDimensions()
	assert.Equal(t, int64(900), imgScaled.GetImage().Width)
	assert.Equal(t, int64(450), imgScaled.GetImage().Height)
}

//go test ./image -run Test_Crop_AutoPosition -v
func Test_Crop_AutoPosition(t *testing.T) {
	img := getMockImageJPEG()
	op, _ := MakeOperations( "crop=200:100;auto,auto", img )
	err := DoTransformation(op)

	assert.Equal(t, int64(200), img.GetImage().Width)
	assert.Equal(t, int64(100), img.GetImage().Height)
	assert.Equal(t, nil, err)
}

//go test ./image -run Test_Resize_Crop_AutoPosition -v
func Test_Resize_Crop_AutoPosition(t *testing.T) {
	img := getMockImageJPEG()
	op, _ := MakeOperations( "resize=250:*&crop=250:100;auto,auto", img )
	err := DoTransformation(op)

	assert.Equal(t, int64(250), img.GetImage().Width)
	assert.Equal(t, int64(100), img.GetImage().Height)
	assert.Equal(t, nil, err)
}
