package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/flexstack.ai/membots-be/internal/util"
)

func Test_ListJoin(t *testing.T) {
	asserts := assert.New(t)
	list := []int{1, 2, 3, 4}
	expectedStr := "1,2,3,4"
	s := util.JoinList(list, ",")
	asserts.Equal(expectedStr, s)

	s = util.JoinList([]int{}, ",")
	asserts.Equal("", s)
}

func Test_RemoveFromList(t *testing.T) {
	asserts := assert.New(t)
	list := []int{1, 2, 3, 4}
	expectedLen := 3
	ret := util.RemoveFromList(list, 2)
	asserts.Equal(expectedLen, len(ret))
}
