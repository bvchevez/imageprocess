package image

import (
	"fmt"
)

// DensityOperation represents the data needed to perform Density
type DensityOperation struct {
	NewDensity int64
	Image      *MutableImage
}

// Do performs the actual density action.
func (i *DensityOperation) Do() error {
	img := *i.Image
	return img.Density(i)
}

// IsValid checks if density is int64 1 or 2, it cannot be anything else.
func (i *DensityOperation) IsValid() bool {
	return (i.NewDensity == int64(1) || i.NewDensity == int64(2))
}

func (i *DensityOperation) String() string {
	return fmt.Sprint("Density")
}
