package handler

import (
	"image-previewer/pkg/imagepreviewer"
	"reflect"
	"testing"

	"github.com/rs/zerolog"
)

func TestNewHandlers(t *testing.T) {
	type args struct {
		logger zerolog.Logger
		svc    imagepreviewer.Service
	}
	tests := []struct {
		name string
		args args
		want *Handlers
	}{
		{
			name: "success_create_new_handlers",
			want: &Handlers{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHandlers(tt.args.logger, tt.args.svc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHandlers() = %v, want %v", got, tt.want)
			}
		})
	}
}
