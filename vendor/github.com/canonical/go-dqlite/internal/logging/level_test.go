package logging_test

import (
	"testing"

	"github.com/canonical/go-dqlite/internal/logging"
	"github.com/stretchr/testify/assert"
)

func TestLevel_String(t *testing.T) {
	assert.Equal(t, "DEBUG", logging.Debug.String())
	assert.Equal(t, "INFO", logging.Info.String())
	assert.Equal(t, "WARN", logging.Warn.String())
	assert.Equal(t, "ERROR", logging.Error.String())

	unknown := logging.Level(666)
	assert.Equal(t, "UNKNOWN", unknown.String())
}
