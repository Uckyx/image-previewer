package imagepreviewer

import (
	"context"
	"reflect"
	"testing"

	"github.com/Uckyx/image-previewer/pkg/cache"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestDefaultService_Fill_positive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	imageOrigin := loadImage(OriginalImgName)
	imageResized := loadImage(ResizedImgName)
	downloadedImage := &DownloadResponse{img: imageOrigin}

	logger := log.With().Logger()
	c := cache.NewCache(2)
	resizer := NewImageResizer(logger)
	mockDownloader := NewMockImageDownloader(ctrl)

	type fields struct {
		l          zerolog.Logger
		cache      cache.Cache
		downloader ImageDownloader
		resizer    ImageResizer
	}
	tests := []struct {
		name    string
		fields  fields
		params  *ResizeRequest
		want    *ResizeResponse
		wantErr bool
	}{
		{
			name: "success_resized",
			fields: fields{
				l:          logger,
				cache:      c,
				downloader: mockDownloader,
				resizer:    resizer,
			},
			params: NewResizeRequest(
				context.Background(),
				1000,
				500,
				ImageURL+OriginalImgName,
				nil,
			),
			want:    &ResizeResponse{imageResized, nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewApp(
				tt.fields.l,
				tt.fields.cache,
				tt.fields.downloader,
				tt.fields.resizer,
			)

			mockDownloader.EXPECT().Download(
				tt.params.ctx,
				tt.params.url,
				tt.params.headers,
			).Return(downloadedImage, nil)

			_, err := svc.Resize(tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fill() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_service_Resize_OriginalImageInCache(t *testing.T) {
	logger := log.With().Logger()
	ctx := context.Background()

	c := cache.NewCache(10)
	svc := NewApp(
		logger,
		c,
		NewImageDownloader(logger),
		NewImageResizer(logger),
	)

	url := ImageURL + OriginalImgName
	key := c.GenerateOriginalImgKey(url)
	img := loadImage(OriginalImgName)

	c.Set(key, img)

	t.Run("get_original_img_from_cache_case", func(t *testing.T) {
		request := NewResizeRequest(
			ctx,
			256,
			126,
			ImageURL+OriginalImgName,
			nil,
		)

		gotImg, err := svc.Resize(request)
		if err != nil {
			t.Errorf("Resize() error = %v", err)
			return
		}

		wantImg := loadImage(ResizedImgName)
		if !reflect.DeepEqual(gotImg.Img, wantImg) {
			t.Errorf("Resize() gotImg = %v, want %v", gotImg.Img, wantImg)
		}
	})
}

func Test_service_Resize_ResizedImageInCache(t *testing.T) {
	logger := log.With().Logger()
	ctx := context.Background()

	c := cache.NewCache(10)
	svc := NewApp(
		logger,
		c,
		NewImageDownloader(logger),
		NewImageResizer(logger),
	)

	url := ImageURL + OriginalImgName
	key := c.GenerateResizedImgKey(url, 256, 126)
	img := loadImage(ResizedImgName)

	c.Set(key, img)

	t.Run("get_resized_img_from_cache_case", func(t *testing.T) {
		request := NewResizeRequest(ctx, 256, 126, url, nil)

		gotImg, err := svc.Resize(request)
		if err != nil {
			t.Errorf("Resize() error = %v", err)
			return
		}

		wantImg := loadImage(ResizedImgName)
		if !reflect.DeepEqual(gotImg.Img, wantImg) {
			t.Errorf("Resize() gotImg = %v, want %v", gotImg.Img, wantImg)
		}
	})
}

func Test_service_Resize_RemoveImageInCache(t *testing.T) {
	logger := log.With().Logger()
	ctx := context.Background()

	c := cache.NewCache(2)
	svc := NewApp(
		logger,
		c,
		NewImageDownloader(logger),
		NewImageResizer(logger),
	)

	url := ImageURL + OriginalImgName

	originalImgKey := c.GenerateOriginalImgKey(url)
	originalImg := loadImage(OriginalImgName)
	c.Set(originalImgKey, originalImg)

	resizedImgKey := c.GenerateResizedImgKey(url, 256, 126)
	resizedImg := loadImage(ResizedImgName)
	c.Set(resizedImgKey, resizedImg)

	t.Run("remove_old_img_from_cache_case", func(t *testing.T) {
		request := NewResizeRequest(ctx, 333, 666, url, nil)

		gotImg, err := svc.Resize(request)
		if err != nil {
			t.Errorf("Resize() error = %v", err)
			return
		}

		wantImg := loadImage("gopher_333x666_resized.jpg")
		if !reflect.DeepEqual(gotImg.Img, wantImg) {
			t.Errorf("Resize() gotImg = %v, want %v", gotImg.Img, wantImg)
		}

		cachedImg, ok := c.Get(originalImgKey)
		require.True(t, ok)
		require.Equal(t, originalImg, cachedImg)

		cachedImg, ok = c.Get(resizedImgKey)
		require.False(t, ok)
		require.Equal(t, []byte(nil), cachedImg)

		resizedImgKey = c.GenerateResizedImgKey(url, 333, 666)
		cachedImg, ok = c.Get(resizedImgKey)
		require.True(t, ok)
		require.Equal(t, wantImg, cachedImg)
	})
}
