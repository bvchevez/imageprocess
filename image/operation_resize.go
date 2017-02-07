package image

import (
	"fmt"
)

// ResizeOperation holds all info necessary to perform a resize on an image.
type ResizeOperation struct {
	NewWidth  int64
	NewHeight int64
	Image     *MutableImage
}

// Do executes the actual Resize action.
func (i *ResizeOperation) Do() error {
	img := *i.Image
	return img.Resize(i)
}

// IsValid checks if resize is even necessary.
func (i *ResizeOperation) IsValid() bool {
	img := *i.Image

	//We cannot "resize" to something that's bigger than the image itself.
	if i.NewWidth > img.GetImage().Width || i.NewHeight > img.GetImage().Height {
		return false
	}

	return true
}

func (i *ResizeOperation) String() string {
	return fmt.Sprint("Resize")
}
