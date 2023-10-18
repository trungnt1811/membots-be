package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatNotiAmount(t *testing.T) {
	nums := []float64{12435475679.12921, 312372137, 234324.23, 43.2}
	expectResults := []string{"12,435,475,679.13", "312,372,137.00", "234,324.23", "43.20"}
	for idx, n := range nums {
		assert.Equal(t, expectResults[idx], FormatNotiAmt(n))
	}
}
