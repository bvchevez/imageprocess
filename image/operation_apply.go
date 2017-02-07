package image

import (
	"fmt"
)

// ApplyOperation represents all data necessary to apply changes to an image.
type ApplyOperation struct {
	Image *MutableImage
}

// Do applies the changes to the image.
func (i *ApplyOperation) Do() error {
	img := *i.Image
	return img.ApplyChanges()
}

// IsValid is a stub function that satisifies OperationInterface.
func (i *ApplyOperation) IsValid() bool {
	return true
}

func (i *ApplyOperation) String() string {
	return fmt.Sprintf("Apply Changes")
}
