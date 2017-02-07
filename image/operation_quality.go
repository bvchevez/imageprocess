package image

import (
	"fmt"
)

// QualityOperation represents all information necessay to perform Quality.
type QualityOperation struct {
	NewQuality int64
	Image      *MutableImage
}

// Do executes the actual Quality operation.
func (i *QualityOperation) Do() error {
	img := *i.Image
	return img.Quality(i)
}

// IsValid verifies that new quality is between 0 and 100 (100 inclusive)
func (i *QualityOperation) IsValid() bool {
	return (i.NewQuality > int64(0) && i.NewQuality <= int64(100))
}

func (i *QualityOperation) String() string {
	return fmt.Sprint("Quality")
}
