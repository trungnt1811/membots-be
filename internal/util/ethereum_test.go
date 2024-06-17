package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWeiToEther(t *testing.T) {
	tests := []struct {
		wei      string
		expected string
	}{
		{"1000000000000000000", "1.00000000"}, // 1 Ether
		{"500000000000000000", "0.50000000"},  // 0.5 Ether
		{"1234567890000000000", "1.23456789"}, // 1.23456789 Ether
	}

	for _, test := range tests {
		result := WeiStrToEtherStr(test.wei)
		require.Equal(t, test.expected, result)
	}
}
