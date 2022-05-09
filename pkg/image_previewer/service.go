package image_previewer

import (
	"context"

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

type Service interface {
	Resize(ctx context.Context, width int, height int, url string) ([]byte, error)
}

type service struct {
	logger          zerolog.Logger
	cache           cache.Cache
	imageDownloader ImageDownloader
	imageResizer    ImageResizer
}

func (s *service) Resize(ctx context.Context, width int, height int, url string) ([]byte, error) {
	rCacheKey := s.cache.GenerateResizedImgKey(url, width, height)
	img, ok := s.cache.Get(rCacheKey)

	if ok {
		return img, nil
	}

	oCacheKey := s.cache.GenerateOriginalImgKey(url)
	oImg, ok := s.cache.Get(oCacheKey)

	if !ok {
		dImg, err := s.imageDownloader.Download(ctx, url)
		if err != nil {
			s.logger.Err(err).Msg(err.Error())
			return nil, err
		}

		s.cache.Set(oCacheKey, dImg)

		oImg = dImg
	}

	rImg, err := s.imageResizer.Resize(ctx, oImg, width, height)
	if err != nil {
		s.logger.Err(err).Msg(err.Error())
		return nil, err
	}

	s.cache.Set(rCacheKey, rImg)

	return rImg, nil
}
