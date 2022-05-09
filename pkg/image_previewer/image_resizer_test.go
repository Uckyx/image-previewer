package image_previewer

import "testing"

import (
	"context"
	"github.com/rs/zerolog/log"
	"reflect"
)

func Test_Positive(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		ctx         context.Context
		originalImg []byte
		resizedImg  []byte
		width       int
		height      int
	}{
		{
			name:        "success resize 256x126",
			ctx:         ctx,
			width:       256,
			height:      126,
			originalImg: loadImage("_gopher_original_1024x504.jpg"),
			resizedImg:  loadImage("gopher_256x126_resized.jpg"),
		},
	}

	logger := log.With().Logger()
	id := NewImageResizer(logger)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotImg, err := id.Resize(tt.ctx, tt.originalImg, tt.width, tt.height)
			if err != nil {
				t.Errorf("Resize() error = %v", err)
				return
			}

			if !reflect.DeepEqual(gotImg, tt.resizedImg) {
				t.Errorf("Resize()  gotImg = %v, want %v", gotImg, tt.resizedImg)
			}
		})
	}
}
