package point

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go test -run Test_Point_String -v
func Test_Point_String(t *testing.T) {
	p := &Point{X: 1, Y: 2}

	assert.Equal(t, "(1,2)", fmt.Sprintf("%s", p))
}
