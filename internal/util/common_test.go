package util

import (
	"fmt"
	"testing"
)

func TestFormatNotiAmount(t *testing.T) {
	number := 12435475679.12921
	out := FormatNotiAmt(number)
	fmt.Println(out)
}
