package imagepreviewer

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func Test_Download_Positive(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		ctx     context.Context
		imgName string
	}{
		{
			name:    "success_download_img",
			ctx:     ctx,
			imgName: OriginalImgName,
		},
	}

	logger := log.With().Logger()
	id := NewImageDownloader(logger)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotImg, err := id.Download(tt.ctx, ImageURL+tt.imgName)
			if err != nil {
				t.Errorf("Download() error = %v", err)
				return
			}

			wantImg := loadImage(tt.imgName)
			if !reflect.DeepEqual(gotImg.img, wantImg) {
				t.Errorf("Download() gotImg = %v, want %v", gotImg.img, wantImg)
			}
		})
	}
}

func Test_Download_Negative(t *testing.T) {
	ctx := context.Background()
	ctxWithTimeOut, closefn := context.WithTimeout(ctx, time.Microsecond*1)
	defer closefn()

	tests := []struct {
		name    string
		ctx     context.Context
		imgName string
		url     string
		err     error
	}{
		{
			name:    "timeout_case",
			ctx:     ctxWithTimeOut,
			imgName: OriginalImgName,
			url:     ImageURL,
			err:     ErrTimeout,
		},
		{
			name:    "not_allowed_type_img_case",
			ctx:     ctx,
			imgName: "gopher_256x126.png",
			url:     ImageURL,
			err:     ErrUnknownImgType,
		},
		{
			name:    "bad_response_case",
			ctx:     ctx,
			imgName: "",
			url:     ImageURL,
			err:     ErrResponseStatus,
		},
	}

	logger := log.With().Logger()
	id := NewImageDownloader(logger)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := id.Download(tt.ctx, tt.url+tt.imgName)
			require.Errorf(t, err, tt.err.Error())
		})
	}
}
