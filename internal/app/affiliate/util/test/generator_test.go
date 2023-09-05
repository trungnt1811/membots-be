package test

import (
	"fmt"
	"testing"

	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/util/vcgen"
)

func Test_generator(t *testing.T) {
	vc, _ := vcgen.NewWithOptions(
		vcgen.SetCount(2), // number of vouchers wanted to generate
		vcgen.SetPattern("########"),
		vcgen.SetPrefix("ACB-"),
		vcgen.SetSuffix("-XYZ"),
	)
	result, err := vc.Run()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result[0])
}
