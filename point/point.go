//represents a point.
package point

import (
	"fmt"
)

type Point struct {
	//the following is used to calculate graph coordinates, where (0,0) is the top left corner.
	X int64
	Y int64
}

func (p *Point) String() string {
	return fmt.Sprintf("(%v,%v)", p.X, p.Y)
}
