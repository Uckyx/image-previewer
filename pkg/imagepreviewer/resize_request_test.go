package imagepreviewer

import (
	"context"
	"reflect"
	"testing"
)

func TestNewResizeRequest(t *testing.T) {
	type args struct {
		ctx     context.Context
		width   int
		height  int
		url     string
		headers map[string][]string
	}
	tests := []struct {
		name string
		args args
		want *ResizeRequest
	}{
		{
			name: "success_generate",
			want: &ResizeRequest{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResizeRequest(
				tt.args.ctx,
				tt.args.width,
				tt.args.height,
				tt.args.url,
				tt.args.headers,
			); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResizeRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
