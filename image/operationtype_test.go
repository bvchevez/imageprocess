package image

import (
	"fmt"
	"testing"

	"github.com/bvchevez/imageprocess/point"
	"github.com/stretchr/testify/assert"
)

type MockedMutableImage struct{}

func (m MockedMutableImage) SetDimensions() error {
	return nil
}

// set default data.
func (i MockedMutableImage) SetDefaults(o Options) {}

func (m MockedMutableImage) ApplyChanges() error {
	return fmt.Errorf("Applying Changes!")
}

func (m MockedMutableImage) Crop(i *CropOperation) error {
	return fmt.Errorf("Crop called!")
}

func (m MockedMutableImage) Resize(i *ResizeOperation) error {
	return fmt.Errorf("Resize called!")
}

func (m MockedMutableImage) Quality(i *QualityOperation) error {
	return fmt.Errorf("Quality called!")
}

func (m MockedMutableImage) Density(i *DensityOperation) error {
	return fmt.Errorf("Density called!")
}

func (m MockedMutableImage) GetImage() *Image {
	return nil
}

func (m MockedMutableImage) Shutdown() {}

func MakeMockMutableImage() MutableImage {
	return MockedMutableImage{}
}

//go test -run Test_CropOperation_String -v
func Test_CropOperation_String(t *testing.T) {
	op := &CropOperation{
		NewWidth:  500,
		NewHeight: 300,
		Position:  &point.Point{X: 1, Y: 1},
	}

	assert.Equal(t, "Crop", fmt.Sprintf("%s", op))
}

//go test -run Test_CropOperation_Do -v
func Test_CropOperation_Do(t *testing.T) {
	img := MakeMockMutableImage()
	op := &CropOperation{
		NewWidth:  500,
		NewHeight: 300,
		Position:  &point.Point{X: 1, Y: 1},
		Image:     &img,
	}

	err := op.Do()
	assert.Equal(t, "Crop called!", err.Error())
}

// go test -run Test_ResizeOperation_String -v
func Test_ResizeOperation_String(t *testing.T) {
	op := &ResizeOperation{
		NewWidth:  500,
		NewHeight: 300,
	}

	assert.Equal(t, "Resize", fmt.Sprintf("%s", op))
}

// go test -run Test_ResizeOperation_MaxWidthIsValidForJPEG -v
func Test_ResizeOperation_MaxWidthIsValidForJPEG(t *testing.T) {
	img := getMockImageJPEG()
	op := &ResizeOperation{
		NewWidth:  5000,
		NewHeight: 300,
		Image:     &img,
	}

	assert.Equal(t, false, op.IsValid())
}

// go test -run Test_ResizeOperation_MaxWidthIsValidForGIF -v
func Test_ResizeOperation_MaxWidthIsValidForGIF(t *testing.T) {
	img := getMockImageGIF()
	op := &ResizeOperation{
		NewWidth:  5000,
		NewHeight: 300,
		Image:     &img,
	}

	assert.Equal(t, false, op.IsValid())
}

// go test -run Test_ResizeOperation_MaxHeightIsValidJPEG -v
func Test_ResizeOperation_MaxHeightIsValidJPEG(t *testing.T) {
	img := getMockImageJPEG()
	op := &ResizeOperation{
		NewWidth:  500,
		NewHeight: 5000,
		Image:     &img,
	}

	assert.Equal(t, false, op.IsValid())
}

// go test -run Test_ResizeOperation_MaxHeightIsValidGIF -v
func Test_ResizeOperation_MaxHeightIsValidGIF(t *testing.T) {
	img := getMockImageGIF()
	op := &ResizeOperation{
		NewWidth:  500,
		NewHeight: 5000,
		Image:     &img,
	}

	assert.Equal(t, false, op.IsValid())
}

//go test -run Test_ResizeOperation_Do -v
func Test_ResizeOperation_Do(t *testing.T) {
	img := MakeMockMutableImage()
	op := &ResizeOperation{
		NewWidth:  500,
		NewHeight: 300,
		Image:     &img,
	}

	err := op.Do()
	assert.Equal(t, "Resize called!", err.Error())
}

//go test -run Test_DensityOperation_String -v
func Test_DensityOperation_String(t *testing.T) {
	img := MakeMockMutableImage()
	op := &DensityOperation{
		Image:      &img,
		NewDensity: 2,
	}

	assert.Equal(t, "Density", fmt.Sprintf("%s", op))
}

//go test -run Test_DensityOperation_IsValidTrue -v
func Test_DensityOperation_IsValidTrue(t *testing.T) {
	img := MakeMockMutableImage()
	op := &DensityOperation{
		Image:      &img,
		NewDensity: 2,
	}

	assert.Equal(t, true, op.IsValid())
}

//go test -run Test_DensityOperation_IsValidFalse -v
func Test_DensityOperation_IsValidFalse(t *testing.T) {
	img := MakeMockMutableImage()
	op := &DensityOperation{
		Image:      &img,
		NewDensity: 0,
	}

	assert.Equal(t, false, op.IsValid())

	op2 := &DensityOperation{
		Image:      &img,
		NewDensity: 5,
	}

	assert.Equal(t, false, op2.IsValid())
}

//go test -run Test_DensityOperation_Do -v
func Test_DensityOperation_Do(t *testing.T) {
	img := MakeMockMutableImage()
	op := &DensityOperation{
		Image:      &img,
		NewDensity: 2,
	}

	err := op.Do()
	assert.Equal(t, "Density called!", err.Error())
}

//go test -run Test_QualityOperation_String -v
func Test_QualityOperation_String(t *testing.T) {
	img := MakeMockMutableImage()
	op := &QualityOperation{
		Image:      &img,
		NewQuality: 95,
	}

	assert.Equal(t, "Quality", fmt.Sprintf("%s", op))
}

//go test -run Test_QualityOperation_IsValidTrue -v
func Test_QualityOperation_IsValidTrue(t *testing.T) {
	img := MakeMockMutableImage()
	op := &QualityOperation{
		Image:      &img,
		NewQuality: 90,
	}

	assert.Equal(t, true, op.IsValid())
}

//go test -run Test_QualityOperation_IsValidFalse -v
func Test_QualityOperation_IsValidFalse(t *testing.T) {
	img := MakeMockMutableImage()
	op := &QualityOperation{
		Image:      &img,
		NewQuality: -12,
	}

	assert.Equal(t, false, op.IsValid())

	op2 := &QualityOperation{
		Image:      &img,
		NewQuality: 500,
	}

	assert.Equal(t, false, op2.IsValid())
}

//go test -run Test_QualityOperation_Do -v
func Test_QualityOperation_Do(t *testing.T) {
	img := MakeMockMutableImage()
	op := &QualityOperation{
		Image:      &img,
		NewQuality: 2,
	}

	err := op.Do()
	assert.Equal(t, "Quality called!", err.Error())
}

//go test -run Test_ApplyOperation_String -v
func Test_ApplyOperation_String(t *testing.T) {
	op := &ApplyOperation{}
	assert.Equal(t, "Apply Changes", fmt.Sprintf("%s", op))
}

//go test -run Test_ApplyOperation_IsValidTrue -v
func Test_ApplyOperation_IsValidTrue(t *testing.T) {
	img := MakeMockMutableImage()
	op := &ApplyOperation{
		Image: &img,
	}
	assert.Equal(t, true, op.IsValid())
}

//go test -run Test_ApplyOperation_Do -v
func Test_ApplyOperation_Do(t *testing.T) {
	img := MakeMockMutableImage()
	op := &ApplyOperation{
		Image: &img,
	}

	err := op.Do()
	assert.Equal(t, "Applying Changes!", err.Error())
}
