package image_previewer

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"image-previewer/pkg/cache"
	"reflect"
	"testing"
)

func TestOriginalImageInCache(t *testing.T) {
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
		gotImg, err := svc.Resize(ctx, 256, 126, ImageURL+OriginalImgName)
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

func TestResizedImageInCache(t *testing.T) {
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
		gotImg, err := svc.Resize(ctx, 256, 126, ImageURL+OriginalImgName)
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

func TestRemoveImageInCache(t *testing.T) {
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
		gotImg, err := svc.Resize(ctx, 333, 666, ImageURL+OriginalImgName)
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
