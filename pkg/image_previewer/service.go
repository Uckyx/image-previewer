package image_previewer

import (
	"context"
	"sync"

	"image-previewer/pkg/cache"

	"github.com/rs/zerolog"
)

func NewApp(
	logger zerolog.Logger,
	cache cache.Cache,
	imageDownloader ImageDownloader,
	imageResizer ImageResizer,
) Service {
	return &service{
		logger:          logger,
		cache:           cache,
		imageDownloader: imageDownloader,
		imageResizer:    imageResizer,
	}
}

type ResizeResponse struct {
	Img     []byte
	Headers map[string][]string
}

type Service interface {
	Resize(ctx context.Context, width int, height int, url string) (*ResizeResponse, error)
}

type service struct {
	logger          zerolog.Logger
	cache           cache.Cache
	imageDownloader ImageDownloader
	imageResizer    ImageResizer
	w               sync.WaitGroup
}

func (s *service) Resize(ctx context.Context, width int, height int, url string) (*ResizeResponse, error) {
	resizedImgKey := s.cache.GenerateResizedImgKey(url, width, height)
	resizedImg, ok := s.cache.Get(resizedImgKey)

	if ok {
		return &ResizeResponse{resizedImg, nil}, nil
	}

	originalImgKey := s.cache.GenerateOriginalImgKey(url)
	originalImg, ok := s.cache.Get(originalImgKey)

	if ok {
		resizedImg, err := s.imageResizer.Resize(ctx, originalImg, width, height)
		if err != nil {
			s.logger.Err(err).Msg(err.Error())

			return nil, err
		}

		s.asyncCacheWrite(resizedImgKey, resizedImg)

		return &ResizeResponse{resizedImg, nil}, nil
	}

	downloadResponse, err := s.imageDownloader.Download(ctx, url)
	if err != nil {
		s.logger.Err(err).Msg(err.Error())

		return nil, err
	}

	s.asyncCacheWrite(originalImgKey, downloadResponse.img)

	resizedImg, err = s.imageResizer.Resize(ctx, downloadResponse.img, width, height)
	if err != nil {
		s.logger.Err(err).Msg(err.Error())

		return nil, err
	}

	s.asyncCacheWrite(resizedImgKey, resizedImg)

	return &ResizeResponse{resizedImg, downloadResponse.headers}, nil
}

func (s *service) asyncCacheWrite(key string, img []byte) {
	s.w.Add(1)
	go func() {
		s.cache.Set(key, img)

		s.w.Done()
	}()

	s.w.Wait()
}
