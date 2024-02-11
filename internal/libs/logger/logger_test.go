package logger

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Creates a new logger with the given configuration
func TestNewLoggerWithGivenConfiguration(t *testing.T) {
	config := NewlogConfigOptions(false)
	logger := NewLogger(config)
	assert.NotNil(t, logger)
}

func TestNewLogger(t *testing.T) {
	type args struct {
		config Iconfigs
	}
	tests := []struct {
		name string
		args args
		want ILogger
	}{
		// TODO: Add test cases.
		{
			name: "Test Logger Creation",
			args: args{NewlogConfigOptions(false)},
			want: NewLogger(NewlogConfigOptions(true)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLogger(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				got.Error(context.Background(), "testing")
			}
		})
	}
}
