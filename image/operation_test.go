package image

import (
	"reflect"
	"testing"

	"github.com/bvchevez/imageprocess/point"
	"github.com/stretchr/testify/assert"
)

//go test -run Test_ImageOperation_Make_Quality -v
func Test_ImageOperation_Make_Quality(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  500,
		ImageHeight: 300,
	}

	op, err := opMaker.Make([]string{"98"}, "output-quality")
	if err != nil {
		t.Errorf("Error not expected, [%v]", err)
	}

	switch typeOp := op.(type) {
	case *QualityOperation:
		assert.Equal(t, "*image.QualityOperation", reflect.TypeOf(typeOp).String())
		assert.Equal(t, int64(98), typeOp.NewQuality)
	default:
		t.Errorf("Other img operations not expected.")
	}
}

//go test -run Test_ImageOperation_Make_Quality__invalidInput -v
func Test_ImageOperation_Make_Quality__invalidInput(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  500,
		ImageHeight: 300,
	}

	expectedError := "invalid quality [asdf]"
	op, err := opMaker.Make([]string{"asdf"}, "output-quality")
	assert.Equal(t, expectedError, err.Error())
	assert.Equal(t, op, nil)

	expectedError = "invalid quality [asdf]"
	_, err = opMaker.Make([]string{"asdf"}, "output-quality")
	assert.Equal(t, expectedError, err.Error())
}

//go test -run Test_ImageOperation_Make_Density -v
func Test_ImageOperation_Make_Density(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  500,
		ImageHeight: 300,
	}

	op, err := opMaker.Make([]string{"2"}, "density")
	if err != nil {
		t.Errorf("Error not expected, [%v]", err)
	}

	switch typeOp := op.(type) {
	case *DensityOperation:
		assert.Equal(t, "*image.DensityOperation", reflect.TypeOf(typeOp).String())
		assert.Equal(t, int64(2), typeOp.NewDensity)
	default:
		t.Errorf("Other img operations not expected.")
	}
}

//go test -run Test_ImageOperation_Make_Density__invalidInput -v
func Test_ImageOperation_Make_Density__invalidInput(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  500,
		ImageHeight: 300,
	}

	expectedError := "invalid density [asdf]"
	op, err := opMaker.Make([]string{"asdf"}, "density")
	assert.Equal(t, expectedError, err.Error())
	assert.Equal(t, op, nil)

	expectedError = "invalid density [4]"
	_, err = opMaker.Make([]string{"4"}, "density")
	assert.Equal(t, expectedError, err.Error())
}

//go test -run Test_ImageOperation_Make_Resize__withWildcardOnHeight -v
func Test_ImageOperation_Make_Resize__withWildcardOnHeight(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 400,
	}

	op, err := opMaker.Make([]string{"100:*"}, "resize")
	if err != nil {
		t.Errorf("Error not expected, [%v]", err)
	}

	switch typeOp := op.(type) {
	case *ResizeOperation:
		assert.Equal(t, "*image.ResizeOperation", reflect.TypeOf(typeOp).String())
		assert.Equal(t, int64(100), typeOp.NewWidth)
		assert.Equal(t, int64(200), typeOp.NewHeight)
	default:
		t.Errorf("Other img operations not expected.")
	}
}

//go test -run Test_ImageOperation_Make_Resize__withWildcardOnWidth -v
func Test_ImageOperation_Make_Resize__withWildcardOnWidth(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 400,
	}

	op, err := opMaker.Make([]string{"*:100"}, "resize")
	if err != nil {
		t.Errorf("Error not expected, [%v]", err)
	}

	switch typeOp := op.(type) {
	case *ResizeOperation:
		assert.Equal(t, "*image.ResizeOperation", reflect.TypeOf(typeOp).String())
		assert.Equal(t, int64(50), typeOp.NewWidth)
		assert.Equal(t, int64(100), typeOp.NewHeight)
	default:
		t.Errorf("Other img operations not expected.")
	}
}

//go test -run Test_ImageOperation_Make_Resize__withRatioOnWidth -v
func Test_ImageOperation_Make_Resize__withRatioOnWidth(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 400,
	}

	op, err := opMaker.Make([]string{"0.5xw:400"}, "resize")
	if err != nil {
		t.Errorf("Error not expected, [%v]", err)
	}

	switch typeOp := op.(type) {
	case *ResizeOperation:
		assert.Equal(t, "*image.ResizeOperation", reflect.TypeOf(typeOp).String())
		assert.Equal(t, int64(100), typeOp.NewWidth)
		assert.Equal(t, int64(400), typeOp.NewHeight)
	default:
		t.Errorf("Other img operations not expected.")
	}
}

//go test -run Test_ImageOperation_Make_Resize__withRatioOnHeight -v
func Test_ImageOperation_Make_Resize__withRatioOnHeight(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 400,
	}

	op, err := opMaker.Make([]string{"100:0.75xh"}, "resize")
	if err != nil {
		t.Errorf("Error not expected, [%v]", err)
	}

	switch typeOp := op.(type) {
	case *ResizeOperation:
		assert.Equal(t, "*image.ResizeOperation", reflect.TypeOf(typeOp).String())
		assert.Equal(t, int64(100), typeOp.NewWidth)
		assert.Equal(t, int64(300), typeOp.NewHeight)
	default:
		t.Errorf("Other img operations not expected.")
	}
}

//go test -run Test_ImageOperation_Make_Resize__withRatioOnWidthAndHeight -v
func Test_ImageOperation_Make_Resize__withRatioOnWidthAndHeight(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 400,
	}

	op, err := opMaker.Make([]string{"0.75xw:0.50xh"}, "resize")
	if err != nil {
		t.Errorf("Error not expected, [%v]", err)
	}

	switch typeOp := op.(type) {
	case *ResizeOperation:
		assert.Equal(t, "*image.ResizeOperation", reflect.TypeOf(typeOp).String())
		assert.Equal(t, int64(150), typeOp.NewWidth)
		assert.Equal(t, int64(200), typeOp.NewHeight)
	default:
		t.Errorf("Other img operations not expected.")
	}
}

//go test -run Test_ImageOperation_Make_Resize__withTwoWildcards__invalid -v
func Test_ImageOperation_Make_Resize__withTwoWildcards__invalid(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 400,
	}

	expectedError := "height and width cannot both be '*'"
	op, err := opMaker.Make([]string{"*:*"}, "resize")
	assert.Equal(t, expectedError, err.Error())
	assert.Equal(t, op, nil)
}

//go test -run Test_ImageOperation_Make_Resize__invalidInputs__tooFewDimensions -v
func Test_ImageOperation_Make_Resize__invalidInputs__tooFewDimensions(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 400,
	}

	expectedError := "resize must have two dimensions"
	op, err := opMaker.Make([]string{"500"}, "resize")
	assert.Equal(t, expectedError, err.Error())
	assert.Equal(t, op, nil)
}

//go test -run Test_ImageOperation_Make_Resize__invalidInputs__tooManyDimensions -v
func Test_ImageOperation_Make_Resize__invalidInputs__tooManyDimensions(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 400,
	}

	expectedError := "resize must have two dimensions"
	op, err := opMaker.Make([]string{"500:200:300"}, "resize")
	assert.Equal(t, expectedError, err.Error())
	assert.Equal(t, op, nil)
}

//go test -run Test_ImageOperation_Make_Resize__invalidInputs__tooManyParameters -v
func Test_ImageOperation_Make_Resize__invalidInputs__tooManyParameters(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 400,
	}

	expectedError := "resize must have dimensions parameter"
	op, err := opMaker.Make([]string{"500:200", "100:300"}, "resize")
	assert.Equal(t, expectedError, err.Error())
	assert.Equal(t, op, nil)
}

//go test -run Test_ImageOperation_Make_Resize__invalidInputs__invalidInputs -v
func Test_ImageOperation_Make_Resize__invalidInputs__invalidInputs(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 400,
	}

	op, err := opMaker.Make([]string{"500:asdf123"}, "resize")
	assert.Equal(t, "dimension must be ratio, number, or '*', not 'asdf123'", err.Error())
	assert.Equal(t, op, nil)

	_, err = opMaker.Make([]string{"500asdf:123"}, "resize")
	assert.Equal(t, "dimension must be ratio, number, or '*', not '500asdf'", err.Error())

	_, err = opMaker.Make([]string{"asdf:123asdf"}, "resize")
	assert.Equal(t, "dimension must be ratio, number, or '*', not 'asdf'", err.Error())
}

//go test -run Test_ImageOperation_Make_Crop__withWildcardOnWidth -v
func Test_ImageOperation_Make_Crop__withWildcardOnWidth(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 400,
	}

	op, err := opMaker.Make([]string{"*:100", "center,center"}, "crop")
	if err != nil {
		t.Errorf("Error not expected, [%v]", err)
	}

	expectedPoint := &point.Point{X: 75, Y: 150}

	switch typeOp := op.(type) {
	case *CropOperation:
		assert.Equal(t, "*image.CropOperation", reflect.TypeOf(typeOp).String())
		assert.Equal(t, int64(50), typeOp.NewWidth)
		assert.Equal(t, int64(100), typeOp.NewHeight)
		assert.Equal(t, expectedPoint, typeOp.Position)
	default:
		t.Errorf("Other img operations not expected.")
	}
}

//go test -run Test_ImageOperation_Make_Crop__cropOutOfBound -v
func Test_ImageOperation_Make_Crop__cropOutOfBound(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 400,
	}

	op, err := opMaker.Make([]string{"*:800", "center,center"}, "crop")
	if err != nil {
		t.Errorf("Error not expected, [%v]", err)
	}

	expectedPoint := &point.Point{X: 0, Y: 0}

	switch typeOp := op.(type) {
	case *CropOperation:
		assert.Equal(t, "*image.CropOperation", reflect.TypeOf(typeOp).String())
		assert.Equal(t, int64(400), typeOp.NewWidth)
		assert.Equal(t, int64(800), typeOp.NewHeight)
		assert.Equal(t, expectedPoint, typeOp.Position)
	default:
		t.Errorf("Other img operations not expected.")
	}
}

//go test -run Test_ImageOperation_Make_Crop__withCustomDimension -v
func Test_ImageOperation_Make_Crop__withCustomDimension(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 400,
	}

	op, err := opMaker.Make([]string{"*:100", "10,30"}, "crop")
	if err != nil {
		t.Errorf("Error not expected, [%v]", err)
	}

	expectedPoint := &point.Point{X: 10, Y: 30}

	switch typeOp := op.(type) {
	case *CropOperation:
		assert.Equal(t, "*image.CropOperation", reflect.TypeOf(typeOp).String())
		assert.Equal(t, int64(50), typeOp.NewWidth)
		assert.Equal(t, int64(100), typeOp.NewHeight)
		assert.Equal(t, expectedPoint, typeOp.Position)
	default:
		t.Errorf("Other img operations not expected.")
	}
}

//go test -run Test_ImageOperation_Make_Crop__withMixedDimension -v
func Test_ImageOperation_Make_Crop__withMixedDimension(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 400,
	}

	op, err := opMaker.Make([]string{"*:100", "left,30"}, "crop")
	if err != nil {
		t.Errorf("Error not expected, [%v]", err)
	}

	expectedPoint := &point.Point{X: 0, Y: 30}

	switch typeOp := op.(type) {
	case *CropOperation:
		assert.Equal(t, "*image.CropOperation", reflect.TypeOf(typeOp).String())
		assert.Equal(t, int64(50), typeOp.NewWidth)
		assert.Equal(t, int64(100), typeOp.NewHeight)
		assert.Equal(t, expectedPoint, typeOp.Position)
	default:
		t.Errorf("Other img operations not expected.")
	}
}

//go test -run Test_ImageOperation_Make_Crop__withInvalidPosition -v
func Test_ImageOperation_Make_Crop__withInvalidPosition(t *testing.T) {
	var err error
	opMaker := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 400,
	}

	_, err = opMaker.Make([]string{"*:100", "0.5xw,asdfxh"}, "crop")
	assert.Equal(t,
		"crop Y position not 'top', 'center', 'bottom' or 'auto', and ratio must be non-zero float, not 'asdfxh'",
		err.Error())

	_, err = opMaker.Make([]string{"*:100", "0.5xw,asdfg"}, "crop")
	assert.Equal(t,
		"crop Y position not 'top', 'center', 'bottom' or 'auto', and explicit position must be ratio or number, not 'asdfg'",
		err.Error())

	_, err = opMaker.Make([]string{"*:100", "0.5xw,3xh"}, "crop")
	assert.Equal(t, "crop Y position is outside image: 1200", err.Error())

	_, err = opMaker.Make([]string{"*:100", "3xw,1xh"}, "crop")
	assert.Equal(t, "crop X position is outside image: 600", err.Error())

	// try with auto and outside image
	_, err = opMaker.Make([]string{"*:100", "auto,3xh"}, "crop")
	assert.Equal(t, "crop Y position is outside image: 1200", err.Error())

	_, err = opMaker.Make([]string{"*:100", "3xw,auto"}, "crop")
	assert.Equal(t, "crop X position is outside image: 600", err.Error())

	_, err = opMaker.Make([]string{"*:100", "0.2xw,0.3xh"}, "crop")
	if err != nil {
		t.Errorf("Error not expected, [%v]", err)
	}
}

//go test -run Test_ImageOperation_Make_Crop__withValidPosition -v
func Test_ImageOperation_Make_Crop__withValidPosition(t *testing.T) {
	var err error
	opMaker := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 400,
	}

	_, err = opMaker.Make([]string{"*:100", "0.2xw,0.3xh"}, "crop")
	if err != nil {
		t.Errorf("Error not expected, [%v]", err)
	}
}

//go test -run Test_ImageOperation_Make_Fill__PUX_5483 -v
// An example of a broken image
// https://hips.hearstapps.com/mp-eui-test.s3.amazonaws.com/images/6823205-large-1467990101.jpg?fill=16:9
// Will attempt to crop at 1922px width, but the actual width of this image is only 1920px.
// That's because it was trying to perform an aspect ratio of 1.78
// Solution is to round down, while increasing the significant digit to 4, so we get 1.7777, giving us 1919 px.
func Test_ImageOperation_Make_Fill__PUX_5483(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  1920,
		ImageHeight: 1080,
	}

	op, err := opMaker.Make([]string{"16:9"}, "fill")
	if err != nil {
		t.Errorf("Error not expected, [%v]", err)
	}

	switch typeOp := op.(type) {
	case *CropOperation:
		assert.Equal(t, "*image.CropOperation", reflect.TypeOf(typeOp).String())
		assert.Equal(t, int64(1919), typeOp.NewWidth)
		assert.Equal(t, int64(1080), typeOp.NewHeight)
		assert.Equal(t, &point.Point{X: 0, Y: 0}, typeOp.Position)
	default:
		t.Errorf("Other img operations not expected.")
	}
}

//go test -run Test_ImageOperation_Make_Fill__Shorter -v
func Test_ImageOperation_Make_Fill__Shorter(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  400,
		ImageHeight: 300,
	}

	op, err := opMaker.Make([]string{"16:9", "left,30"}, "fill")
	if err != nil {
		t.Errorf("Error not expected, [%v]", err)
	}

	switch typeOp := op.(type) {
	case *CropOperation:
		assert.Equal(t, "*image.CropOperation", reflect.TypeOf(typeOp).String())
		assert.Equal(t, int64(400), typeOp.NewWidth)
		assert.Equal(t, int64(225), typeOp.NewHeight)
		assert.Equal(t, &point.Point{X: 0, Y: 30}, typeOp.Position)
	default:
		t.Errorf("Other img operations not expected.")
	}
}

//go test -run Test_ImageOperation_Make_Fill__Narrower -v
func Test_ImageOperation_Make_Fill__Narrower(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  1600,
		ImageHeight: 900,
	}

	op, err := opMaker.Make([]string{"4:3", "30,top"}, "fill")
	if err != nil {
		t.Errorf("Error not expected, [%v]", err)
	}

	switch typeOp := op.(type) {
	case *CropOperation:
		assert.Equal(t, "*image.CropOperation", reflect.TypeOf(typeOp).String())
		assert.Equal(t, int64(1199), typeOp.NewWidth)
		assert.Equal(t, int64(900), typeOp.NewHeight)
		assert.Equal(t, &point.Point{X: 30, Y: 0}, typeOp.Position)
	default:
		t.Errorf("Other img operations not expected.")
	}
}

//go test -run Test_ImageOperation_Make_Fill__Bad_Height -v
func Test_ImageOperation_Make_Fill__Bad_Height(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  1600,
		ImageHeight: 900,
	}

	_, err := opMaker.Make([]string{"4:3x", "30,top"}, "fill")
	assert.Equal(t, "fill aspect height must be integer, not 3x", err.Error())
}

//go test -run Test_ratio -v
func Test_ratio(t *testing.T) {
	assert := assert.New(t)

	// <1 ratio
	result, err := ratio(15, 200)
	assert.Equal(nil, err)
	assert.Equal(float64(0.075000), result)

	// >1 ratio
	result, err = ratio(1000, 200)
	assert.Equal(nil, err)
	assert.Equal(float64(5), result)

	// 0 is allowed for numerator
	result, err = ratio(0, 200)
	assert.Equal(nil, err)
	assert.Equal(float64(0), result)
}

//go test -run Test_ratio_zero_denominator -v
func Test_ratio_zero_denominator(t *testing.T) {
	assert := assert.New(t)

	// 0 is not allowed for denumerator
	result, err := ratio(15, 0)
	assert.Equal("denominator of ratio can't be zero", err.Error())
	assert.Equal(float64(0), result)
}

//go test -run Test_ImageOperation_dims2px_normalRun -v
func Test_ImageOperation_dims2px_normalRun(t *testing.T) {
	op := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 300,
	}
	assert := assert.New(t)

	//wildcard on width
	op.dims2px("*", "100")
	assert.Equal(int64(66), op.NewWidth)
	assert.Equal(int64(100), op.NewHeight)

	//wildcard on height
	op.dims2px("900", "*")
	assert.Equal(int64(900), op.NewWidth)
	assert.Equal(int64(1350), op.NewHeight)

	//0 as height
	op.dims2px("*", "0")
	assert.Equal(int64(0), op.NewWidth)
	assert.Equal(int64(0), op.NewHeight)

	//no wildcard.
	op.dims2px("100", "200")
	assert.Equal(int64(100), op.NewWidth)
	assert.Equal(int64(200), op.NewHeight)
}

//go test -run Test_ImageOperation_dims2px_invalidHeightInt -v
func Test_ImageOperation_dims2px_invalidHeightInt(t *testing.T) {
	op := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 300,
	}

	op.dims2px("*", "asdf")
	assert.Equal(t, int64(0), op.NewWidth)
	assert.Equal(t, int64(0), op.NewHeight)
}

//go test -run Test_ImageOperation_dims2px_invalidWidthInt -v
func Test_ImageOperation_dims2px_invalidWidthInt(t *testing.T) {
	op := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 300,
	}

	op.dims2px("asdf", "*")
	assert.Equal(t, int64(0), op.NewWidth)
	assert.Equal(t, int64(0), op.NewHeight)
}

//go test -run Test_ImageOperation_dims2px_badRatioWidth -v
func Test_ImageOperation_dims2px_badRatioWidth(t *testing.T) {
	op := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 0,
	}

	err := op.dims2px("*", "100")
	assert.Equal(t, int64(0), op.NewWidth)
	assert.Equal(t, int64(100), op.NewHeight)
	assert.Equal(t, "denominator of ratio can't be zero", err.Error())
}

//go test -run Test_ImageOperation_dims2px_badRatioHeight -v
func Test_ImageOperation_dims2px_badRatioHeight(t *testing.T) {
	op := ImageOperation{
		ImageWidth:  0,
		ImageHeight: 1000,
	}

	err := op.dims2px("100", "*")
	assert.Equal(t, int64(100), op.NewWidth)
	assert.Equal(t, int64(0), op.NewHeight)
	assert.Equal(t, "denominator of ratio can't be zero", err.Error())
}

//go test -run Test_ImageOperation_dims2px_allRatio -v
func Test_ImageOperation_dims2px_allRatio(t *testing.T) {
	op := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 300,
	}

	err := op.dims2px("1.5xw", "1xh")
	assert.Equal(t, int64(300), op.NewWidth)
	assert.Equal(t, int64(300), op.NewHeight)
	assert.Equal(t, nil, err)
}

//go test -run Test_ImageOperation_dims2px_allRatioReverseWAndH -v
func Test_ImageOperation_dims2px_allRatioReverseWAndH(t *testing.T) {
	op := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 300,
	}

	err := op.dims2px("1.5xh", "1xw")
	assert.Equal(t, int64(450), op.NewWidth)
	assert.Equal(t, int64(200), op.NewHeight)
	assert.Equal(t, nil, err)
}

//go test -run Test_ImageOperation_dims2px_reverseRatioWithWildcard -v
func Test_ImageOperation_dims2px_reverseRatioWithWildcard(t *testing.T) {
	op := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 300,
	}

	err := op.dims2px("1.5xh", "*") //shuld be 1.5xw
	assert.Equal(t, int64(450), op.NewWidth)
	assert.Equal(t, int64(675), op.NewHeight)
	assert.Equal(t, nil, err)
}

//go test -run Test_ImageOperation_dims2px_ratioWithWildcard -v
func Test_ImageOperation_dims2px_ratioWithWildcard(t *testing.T) {
	op := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 300,
	}

	err := op.dims2px("1.5xw", "*")
	assert.Equal(t, int64(300), op.NewWidth)
	assert.Equal(t, int64(450), op.NewHeight)
	assert.Equal(t, nil, err)
}

//go test -run Test_ImageOperation_dims2px_withTwoWildcard -v
func Test_ImageOperation_dims2px_withTwoWildcard(t *testing.T) {
	op := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 300,
	}

	err := op.dims2px("*", "*")
	assert.Equal(t, int64(0), op.NewWidth)
	assert.Equal(t, int64(0), op.NewHeight)
	assert.Equal(t, "height and width cannot both be '*'", err.Error())
}

//go test -run Test_ImageOperation_dims2px_pixelAndRatio -v
func Test_ImageOperation_dims2px_pixelAndRatio(t *testing.T) {
	op := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 300,
	}

	err := op.dims2px("240", "1.5xh")
	assert.Equal(t, int64(240), op.NewWidth)
	assert.Equal(t, int64(450), op.NewHeight)
	assert.Equal(t, nil, err)
}

//go test -run Test_ImageOperation_dims2px_pixelAndWildcard -v
func Test_ImageOperation_dims2px_pixelAndWildcard(t *testing.T) {
	op := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 300,
	}

	err := op.dims2px("240", "*")
	assert.Equal(t, int64(240), op.NewWidth)
	assert.Equal(t, int64(360), op.NewHeight)
	assert.Equal(t, nil, err)
}

//go test -run Test_ImageOperation_dims2px_invalidRatioHeight -v
func Test_ImageOperation_dims2px_invalidRatioHeight(t *testing.T) {
	op := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 300,
	}

	err := op.dims2px("240", "asdfxh")
	assert.Equal(t, int64(240), op.NewWidth)
	assert.Equal(t, int64(0), op.NewHeight)
	assert.Equal(t, "ratio must be non-zero float, not 'asdfxh'", err.Error())
}

//go test -run Test_ImageOperation_dims2px_invalidRatioWidth -v
func Test_ImageOperation_dims2px_invalidRatioWidth(t *testing.T) {
	op := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 300,
	}

	err := op.dims2px("axw", "asdfxh")
	assert.Equal(t, int64(0), op.NewWidth)
	assert.Equal(t, int64(0), op.NewHeight)
	assert.Equal(t, "ratio must be non-zero float, not 'axw'", err.Error())
}

//go test -run Test_ImageOperation_setCropPosition_center_normalValues -v
func Test_ImageOperation_setCropPosition_center_normalValues(t *testing.T) {
	op := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 300,
	}
	assert := assert.New(t)

	op.ImageHeight = 600
	op.ImageWidth = 300
	op.NewHeight = 400
	op.NewWidth = 200

	op.setCropPosition("center", "bottom")
	assert.Equal(int64(50), op.Position.X)
	assert.Equal(int64(200), op.Position.Y)

	op.ImageHeight = 600
	op.ImageWidth = 350
	op.NewHeight = 400
	op.NewWidth = 200

	op.setCropPosition("center", "bottom")
	assert.Equal(int64(75), op.Position.X)
	assert.Equal(int64(200), op.Position.Y)

	//two centers
	op.ImageHeight = 600
	op.ImageWidth = 350
	op.NewHeight = 400
	op.NewWidth = 200

	op.setCropPosition("center", "center")
	assert.Equal(int64(75), op.Position.X)
	assert.Equal(int64(100), op.Position.Y)
}

//go test -run Test_ImageOperation_setCropPosition_center_oddValues -v
func Test_ImageOperation_setCropPosition_center_oddValues(t *testing.T) {
	op := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 300,
	}
	assert := assert.New(t)

	//crop width is equal to actual width.
	op.ImageHeight = 600
	op.ImageWidth = 350
	op.NewHeight = 400
	op.NewWidth = 350

	op.setCropPosition("center", "bottom")
	assert.Equal(int64(0), op.Position.X)
	assert.Equal(int64(200), op.Position.Y)

	//crop height is equal to actual height.
	op.ImageHeight = 600
	op.ImageWidth = 350
	op.NewHeight = 600
	op.NewWidth = 350

	op.setCropPosition("center", "center")
	assert.Equal(int64(0), op.Position.X)
	assert.Equal(int64(0), op.Position.Y)

	//crop values greater to actual values.
	op.ImageHeight = 600
	op.ImageWidth = 350
	op.NewHeight = 680
	op.NewWidth = 380

	op.setCropPosition("center", "center")
	assert.Equal(int64(0), op.Position.X)
	assert.Equal(int64(0), op.Position.Y)
}

//go test -run Test_ImageOperation_setCropPosition_top_right_normalValues -v
func Test_ImageOperation_setCropPosition_top_right_normalValues(t *testing.T) {
	op := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 300,
	}
	assert := assert.New(t)

	op.ImageWidth = 350
	op.ImageHeight = 600
	op.NewWidth = 300
	op.NewHeight = 400

	op.setCropPosition("right", "top")
	assert.Equal(int64(50), op.Position.X)
	assert.Equal(int64(0), op.Position.Y)

	op.ImageWidth = 350
	op.ImageHeight = 600
	op.NewWidth = 300
	op.NewHeight = 400

	op.setCropPosition("right", "center")
	assert.Equal(int64(50), op.Position.X)
	assert.Equal(int64(100), op.Position.Y)

	op.ImageWidth = 350
	op.ImageHeight = 600
	op.NewWidth = 300
	op.NewHeight = 400

	op.setCropPosition("center", "top")
	assert.Equal(int64(25), op.Position.X)
	assert.Equal(int64(0), op.Position.Y)
}

//go test -run Test_ImageOperation_setCropPosition_top_right_oddValues -v
func Test_ImageOperation_setCropPosition_top_right_oddValues(t *testing.T) {
	op := ImageOperation{
		ImageWidth:  200,
		ImageHeight: 300,
	}
	assert := assert.New(t)

	op.ImageWidth = 350
	op.ImageHeight = 600
	op.NewWidth = 400
	op.NewHeight = 680

	op.setCropPosition("right", "top")
	assert.Equal(int64(0), op.Position.X)
	assert.Equal(int64(0), op.Position.Y)
}

//go test -run Test_Make_InvalidOperation -v
func Test_Make_InvalidOperation(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  500,
		ImageHeight: 300,
	}

	_, err := opMaker.Make([]string{"12"}, "test")
	assert.Equal(t, "invalid operation test", err.Error())
}

//go test -run Test_Make_NoParams -v
func Test_Make_NoParams(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  500,
		ImageHeight: 300,
	}

	_, err := opMaker.Make([]string{}, "resize")
	assert.Equal(t, "at least one dimension is required", err.Error())
}

//go test -run Test_Make_NoInputs -v
func Test_Make_NoInputs(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  500,
		ImageHeight: 300,
	}

	_, err := opMaker.Make([]string{}, "")
	assert.Equal(t, "at least one dimension is required", err.Error())
}

//go test -run Test_setOutputQuality_invalidDimension -v
func Test_setOutputQuality_invalidDimension(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  500,
		ImageHeight: 300,
	}

	err := opMaker.setOutputQuality([]string{"1", "2"})
	assert.Equal(t, "too many dimensions. Maximum number of dimensions for quality is 1", err.Error())
}

//go test -run Test_setOutputQuality_invalidQualityNonNumeric -v
func Test_setOutputQuality_invalidQualityNonNumeric(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  500,
		ImageHeight: 300,
	}

	err := opMaker.setOutputQuality([]string{"asdf"})
	assert.Equal(t, "invalid quality [asdf]", err.Error())
}

//go test -run Test_setOutputQuality_invalidQualityOutOfBounds -v
func Test_setOutputQuality_invalidQualityOutOfBounds(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  500,
		ImageHeight: 300,
	}

	err := opMaker.setOutputQuality([]string{"200"})
	assert.Equal(t, "invalid quality [200]", err.Error())
}

//go test -run Test_setFrame_invalidDimension -v
func Test_setFrame_invalidDimension(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  500,
		ImageHeight: 300,
	}

	err := opMaker.setFrame([]string{"1", "2"})
	assert.Equal(t, "too many dimensions. Maximum number of dimensions for frame is 1", err.Error())
}

//go test -run Test_setFrame_SetToTrue -v
func Test_setFrame_SetToTrue(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  500,
		ImageHeight: 300,
	}

	opMaker.setFrame([]string{"1"})
	assert.Equal(t, true, opMaker.NewFrame)
}

//go test -run Test_setFrame_SetToFalse -v
func Test_setFrame_SetToFalse(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  500,
		ImageHeight: 300,
	}

	opMaker.setFrame([]string{"0"})
	assert.Equal(t, false, opMaker.NewFrame)
}

//go test -run Test_setDensity_invalidDimension -v
func Test_setDensity_invalidDimension(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  500,
		ImageHeight: 300,
	}

	err := opMaker.setDensity([]string{"1", "2"})
	assert.Equal(t, "too many dimensions. Maximum number of dimensions for density is 1", err.Error())
}

//go test -run Test_setCrop_invalidDimension -v
func Test_setCrop_invalidDimension(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  500,
		ImageHeight: 300,
	}

	err := opMaker.setCrop([]string{"1", "2", "3"})
	assert.Equal(t, "too many parameters for crop", err.Error())

	err = opMaker.setCrop([]string{"1"})
	assert.Equal(t, "crop size must have two dimensions", err.Error())
}

//go test -run Test_setCrop_badInputBothWildcard -v
func Test_setCrop_badInputBothWildcard(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  500,
		ImageHeight: 300,
	}

	err := opMaker.setCrop([]string{"*:*", "center,center"})
	assert.Equal(t, "height and width cannot both be '*'", err.Error())
}

//go test -run Test_setCrop_badInputWidth -v
func Test_setCrop_badInputWidth(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  500,
		ImageHeight: 300,
	}

	err := opMaker.setCrop([]string{"0.23xx:*", "center,center"})
	assert.Equal(t, "dimension must be 'w', 'h', 'g', or 'l', not 'x'", err.Error())
}

//go test -run Test_setCrop_badInputHeight -v
func Test_setCrop_badInputHeight(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  500,
		ImageHeight: 300,
	}

	err := opMaker.setCrop([]string{"*:0.23xx", "center,center"})
	assert.Equal(t, "dimension must be 'w', 'h', 'g', or 'l', not 'x'", err.Error())
}

//go test -run Test_setCrop_incorrectDim -v
func Test_setCrop_incorrectDim(t *testing.T) {
	opMaker := ImageOperation{
		ImageWidth:  500,
		ImageHeight: 300,
	}

	err := opMaker.setCrop([]string{"0.23xx", "center,center"})
	assert.Equal(t, "crop size must have two dimensions", err.Error())

	err = opMaker.setCrop([]string{"0.23xw:*", "center"})
	assert.Equal(t, "crop position must have two coordinates", err.Error())
}

//go test -run Test_setNewDimensions_bothWildcard -v
func Test_setNewDimensions_bothWildcard(t *testing.T) {
	op := ImageOperation{
		ImageWidth:  500,
		ImageHeight: 300,
	}

	err := op.setNewDimensions("*", "*")
	assert.Equal(t, int64(0), op.NewHeight)
	assert.Equal(t, int64(0), op.NewWidth)
	assert.Equal(t, "height and width cannot both be '*'", err.Error())
}
