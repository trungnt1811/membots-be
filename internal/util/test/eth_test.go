package test

import (
	"github.com/astraprotocol/affiliate-system/internal/util"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UserAmountToWei(t *testing.T) {
	asserts := assert.New(t)

	s := "0.1"
	i, err := util.UserAmountToWei(s, big.NewInt(18))
	asserts.Nil(err)

	asserts.Equal("100000000000000000", i.String())
}
