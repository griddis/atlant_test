package logging

import (
	"reflect"
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger("debug", "2006-01-02T15:04:05.999999999Z07:00")
	type args struct {
		level      string
		timeFormat string
	}
	tests := []struct {
		name string
		args args
		want *Logger
	}{
		{
			"new",
			args{
				"debug",
				"2006-01-02T15:04:05.999999999Z07:00",
			},
			logger,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLogger(tt.args.level, tt.args.timeFormat); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}
