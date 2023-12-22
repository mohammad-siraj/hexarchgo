package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Returns a valid Iconfigs object with default values when no options are passed
func TestDefaultValues(t *testing.T) {
	config := NewlogConfigOptions(false)
	assert.Equal(t, false, config.IsProduction())
	assert.Equal(t, "", config.FileName())
	assert.Equal(t, 1024, config.MaxSize())
	assert.Equal(t, 30, config.MaxBackups())
	assert.Equal(t, 90, config.MaxAge())
	assert.Equal(t, true, config.IslocalTime())
	assert.Equal(t, true, config.IsCompressed())
}
