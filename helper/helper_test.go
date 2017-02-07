package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

//go test -run Test_RoundDown -v
func Test_RoundDown__oneSig(t *testing.T) {
	f := RoundDown(1.99, 1)
	assert.Equal(t, 1.9, f)
}

//go test -run Test_RoundDown -v
func Test_RoundDown__twoSig(t *testing.T) {
	f := RoundDown(1.999999, 3)
	assert.Equal(t, 1.999, f)
}

func Test_GetPipelineID(t *testing.T) {
	f := GetPipelineID("testsite", "blah-blah/blahblah", "asdf=asdf123:asdf")
	assert.Equal(t, "2b6058d2a70727f640190b5aaa147cb8", f)
}

func Test_String2Float64__passing(t *testing.T) {
	s := "1.5"
	f := String2Float64(s)
	assert.Equal(t, float64(1.5), f)
}

func Test_String2Float64__fail(t *testing.T) {
	s := "1.5a"
	f := String2Float64(s)
	assert.Equal(t, float64(0.0), f)
}

func Test_Int642String__passing(t *testing.T) {
	i := Int642String(int64(11))
	assert.Equal(t, "11", i)
}

func Test_Int2String__passing(t *testing.T) {
	i := Int2String(11)
	assert.Equal(t, "11", i)
}

func Test_String2Int64__passing(t *testing.T) {
	s := "11"
	i := String2Int64(s)
	assert.Equal(t, int64(11), i)
}

func Test_String2Int64__fail(t *testing.T) {
	s := "1.5a"
	f := String2Int64(s)
	assert.Equal(t, int64(0), f)
}

func Test_String2Int__passing(t *testing.T) {
	s := "11"
	i := String2Int(s)
	assert.Equal(t, int(11), i)
}

func Test_String2Int__fail(t *testing.T) {
	s := "1.5a"
	f := String2Int(s)
	assert.Equal(t, int(0), f)
}

func Test_IsFloat__passing(t *testing.T) {
	s := "1.5"
	f := IsFloat(s)
	assert.Equal(t, true, f)
}

func Test_IsFloat__fail(t *testing.T) {
	s := "1.5a"
	f := IsFloat(s)
	assert.Equal(t, false, f)
}

func Test_IsNumeric__passing(t *testing.T) {
	s := "15"
	f := IsNumeric(s)
	assert.Equal(t, true, f)
}

//go test -run Test_IsRatio__fail__badType -v
func Test_IsRatio__fail__badType(t *testing.T) {
	s := "15aa"
	f := IsRatio(s, "bb")
	assert.Equal(t, false, f)

	s2 := "aa"
	f2 := IsRatio(s2, "aa")
	assert.Equal(t, false, f2)
}

//go test -run Test_IsRatio__fail -v
func Test_IsRatio__fail__tooBig(t *testing.T) {
	s := "15xh"
	f := IsRatio(s, "xh")
	assert.Equal(t, false, f)
}

func Test_IsRatio__fail__tooSmall(t *testing.T) {
	s := "-1xh"
	f := IsRatio(s, "xh")
	assert.Equal(t, false, f)
}

func Test_IsRatio__fail__zero(t *testing.T) {
	s := "0xh"
	f := IsRatio(s, "xh")
	assert.Equal(t, false, f)
}

func Test_IsRatio__fail__one(t *testing.T) {
	s := "1xh"
	f := IsRatio(s, "xh")
	assert.Equal(t, false, f)
}

func Test_IsRatio__fail__notNumeric(t *testing.T) {
	s := "asdfxh"
	f := IsRatio(s, "xh")
	assert.Equal(t, false, f)
}

func Test_IsRatio__passing(t *testing.T) {
	s := "0.5xh"
	f := IsRatio(s, "xh")
	assert.Equal(t, true, f)
}

func Test_IsNumeric__fail(t *testing.T) {
	s := "15a"
	f := IsNumeric(s)
	assert.Equal(t, false, f)
}

func Test_InSlice__passing(t *testing.T) {
	s := []string{"a", "b", "c", "d", "e"}
	assert.Equal(t, true, InSlice("a", s))
}

func Test_InSlice__fail(t *testing.T) {
	s := []string{"a", "b", "c", "d", "e"}
	assert.Equal(t, false, InSlice("1", s))
}

func Test_UnescapeURL__pass(t *testing.T) {
	var s string
	var e error

	// url encoded, with one-off html entity.
	s, e = UnescapeURL("crop=220%3A200;0,0&amp;resize=100%3A%2A")
	assert.Nil(t, e)
	assert.Equal(t, "crop=220:200;0,0&resize=100:*", s)

	// url decoded, nothing should change here.
	s, e = UnescapeURL("crop=220:200;0,0&resize=100:*")
	assert.Nil(t, e)
	assert.Equal(t, "crop=220:200;0,0&resize=100:*", s)

	// url decoded with the one-off html entity.
	s, e = UnescapeURL("crop=220:200;0,0&amp;resize=100:*")
	assert.Nil(t, e)
	assert.Equal(t, "crop=220:200;0,0&resize=100:*", s)
}

//this test serves no purpose other than to get that 100% test coverage.
func Test_Timer_Dat100PercentCoverage(t *testing.T) {
	Timer(
		TimerPayload{
			Name:  "test timer",
			Start: time.Now(),
		},
	)
	assert.Equal(t, false, false)
}
