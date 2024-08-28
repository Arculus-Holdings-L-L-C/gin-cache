package startup

import (
	"testing"

	"github.com/Arculus-Holdings-L-L-C/gin-cache/pkg/define"
	"github.com/stretchr/testify/assert"
)

func TestMemCache(t *testing.T) {
	type args struct {
		onCacheHit []define.OnCacheHit
	}
	tests := []struct {
		name    string
		args    args
		success bool
	}{
		{name: "init success", args: args{}, success: true},
		{name: "init error", args: args{}, success: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MemCache(tt.args.onCacheHit...)
			if err != nil {
				assert.Error(t, err)
				return
			}
			assert.True(t, got != nil, tt.success)
		})
	}
}
