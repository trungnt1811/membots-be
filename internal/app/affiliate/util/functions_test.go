package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveUnicode(t *testing.T) {

	str := "Điện Máy Xanh Nè"

	out := ToAsciiString(str)

	fmt.Println("out", out)

	assert.Equal(t, out, "Dien May Xanh Ne", "Fail to remove unicode")
}
