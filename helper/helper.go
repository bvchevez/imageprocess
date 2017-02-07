package helper

import (
	"crypto/md5"
	"fmt"
	"io"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

var (
	Logging = true // Logging turns off logging, useful for doing stuff like benchmark, where we don't want logging to get in the way.
)

// TimerPayload represents the data being passed into the timer helper.
type TimerPayload struct {
	Name  string
	Start time.Time
}

func RoundDown(input float64, places int) float64 {
	pow := math.Pow(10, float64(places))
	digit := pow * input
	return math.Floor(digit) / pow
}

//PipelineID creates a unique pipeline id using site/path/rawQuery.
func GetPipelineID(site, path, rawQuery string) string {
	h := md5.New()
	io.WriteString(h, fmt.Sprintf("%s %s %s", site, path, rawQuery))
	return fmt.Sprintf("%x", h.Sum(nil))
}

//String2Float64 converts string to float64, does not return error, returns float64(0) instead.
func String2Float64(v string) float64 {
	i, err := strconv.ParseFloat(v, 64)
	if err != nil {
		log.WithFields(log.Fields{
			"string": v,
			"error":  err.Error(),
		}).Warn("Error while converting string to float.")
		return 0.0
	}

	return float64(i)
}

//Int642String converts int64 to string
func Int642String(v int64) string {
	return strconv.Itoa(int(v))
}

//Int642String converts int64 to string
func Int2String(v int) string {
	return strconv.Itoa(v)
}

//String2Int64 converts string to int64, does not return error, returns int64(0) instead.
func String2Int64(v string) int64 {
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		log.WithFields(log.Fields{
			"string": v,
			"error":  err.Error(),
		}).Warn("Error while converting string to int64.")
		return 0
	}

	return i
}

//String2Int64 converts string to int64, does not return error, returns int64(0) instead.
func String2Int(v string) int {
	return int(String2Int64(v))
}

//IsFloat returns true if string is in float format, false if not. does not return error.
func IsFloat(v string) bool {
	if _, err := strconv.ParseFloat(v, 64); err == nil {
		return true
	}

	return false
}

//IsNumeric returns true if string is numeric, false if not. does not return error.
func IsNumeric(v string) bool {
	if _, err := strconv.Atoi(v); err == nil {
		return true
	}

	return false
}

//IsRatio determines if our input is a ratio and if it is must be less than 1 and greater than 0
//valid ratios are "0.5xw" where "xw" is type.
func IsRatio(r, typ string) bool {
	if len(r) < 2 || r[len(r)-2:] != typ {
		return false
	}

	rNum := String2Float64(strings.TrimRight(r, typ))
	return rNum > 0 && rNum < 1
}

// InSlice verifies that a specified string is present in a slice of string
func InSlice(a string, list []string) bool {
	// Loop through list looking for a match
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}

// UnescapeURL decodes any urlencoded substrings.
// In addition, it also supports customized one-off replacements as requested by specific tickets.
//     maps "&amp;" to "&"
func UnescapeURL(input string) (string, error) {
	ueURL, err := url.QueryUnescape(input)
	if err != nil {
		return "", err
	}

	// Custom unescaping stuff.
	ueURL = strings.Replace(ueURL, "&amp;", "&", -1) // PUX-5392

	return ueURL, nil
}

//Timer takes in a Time parameter and a name and logs down the difference between Time parameter and now.
func Timer(p TimerPayload) {
	if Logging == false {
		return
	}
	log.WithFields(log.Fields{
		"time": time.Now().Sub(p.Start),
	}).Info(p.Name)
}
