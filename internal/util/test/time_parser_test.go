package test

import (
	"testing"

	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestParsePostBackTime(t *testing.T) {
	asserts := assert.New(t)
	s := "2023-08-24 14:53:02.000000"

	res, err := util.ParsePostBackTime(s)
	asserts.NoError(err)
	asserts.Equal(2023, res.Year())
}

func TestGetSinceUntilTime(t *testing.T) {
	asserts := assert.New(t)
	s := "2023-08-24 14:53:02.000000"

	res, err := util.ParsePostBackTime(s)
	asserts.NoError(err)
	asserts.Equal(2023, res.Year())

	since, until := util.GetSinceUntilTime(res, 1)
	asserts.Equal(until.Day(), since.Day()+2)
}
