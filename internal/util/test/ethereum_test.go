package test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/flexstack.ai/membots-be/internal/util"
)

func TestWeiToEther(t *testing.T) {
	tests := []struct {
		wei      string
		expected string
	}{
		{"1000000000000000000", "1"},          // 1 Ether
		{"500000000000000000", "0.5"},         // 0.5 Ether
		{"100300000000000000", "0.1003"},      // 0.1003 Ether
		{"1234567890000000000", "1.23456789"}, // 1.23456789 Ether
	}

	for _, test := range tests {
		result := util.WeiStrToEtherStr(test.wei)
		require.Equal(t, test.expected, result)
	}
}
