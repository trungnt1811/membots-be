package validators

import (
	"errors"
	"strings"
)

func ValidateAddress(s string) error {
	if !strings.HasPrefix(s, "0x") || len(s) != 42 {
		return errors.New("wallet address malform")
	}
	return nil
}
