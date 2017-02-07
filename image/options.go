package image

import (
    "os"
)

const (
	JPEG = "image/jpeg" // jpeg mime
	PNG  = "image/png"  // png mime
	GIF  = "image/gif"  // gif mime
	TIFF = "image/tiff" // tiff mime

	HaarCascadesPath = "/data/haarcascades/"
	HaarCascadeFrontalFaceAlt = HaarCascadesPath + "haarcascade_frontalface_alt.xml"
)

var (
	// Map of file types to byte slices
	// Used to determine file type when not getting from s3
	fileTypes = map[string][]byte{
		JPEG: []byte{0xff, 0xd8},
		PNG:  []byte{0x89, 0x50},
		GIF:  []byte{0x47, 0x49},
		TIFF: []byte{0x49, 0x49},
	}

	// allowedCustomValues reprents all allowed cropping
	allowedCustomValues = []string{
		"top",
		"bottom",
		"left",
		"right",
		"center",
	}
	// allowedOperations represents all allowed operations that HIPS supports
	allowedOperations = []string{
		"resize",
		"crop",
		"output-quality",
		"density",
		"frame",
	}
	// maxOperations represents the maximum operations allowed per request
	maxOperations int = 5

	// interlace represents the Interlace option of libvips.
	interlace bool = true

	// maximum dimensions we want to set our width and height to
	maxWidth  int64 = 3000
	maxHeight int64 = 3000

	goPath string = os.Getenv("GOPATH")
	appPath string = goPath + "/src/github.com/bvchevez/imageprocess"
)

// Options represents different image options available.
type Options struct {
	Quality          int64
	Density          int64
	BicubicThreshold int64
}
