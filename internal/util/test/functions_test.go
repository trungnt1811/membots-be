package test

import (
	"github.com/astraprotocol/affiliate-system/internal/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RequestIdToImportId(t *testing.T) {
	asserts := assert.New(t)

	importId := 1
	requestId := util.ImportIdToRequestId(importId)
	asserts.Equal("import_id:1", requestId)

	parsed := util.RequestIdToImportId(requestId)
	asserts.Equal(importId, parsed)

	parsed = util.RequestIdToImportId("import_id:")
	asserts.Equal(0, parsed)
}
