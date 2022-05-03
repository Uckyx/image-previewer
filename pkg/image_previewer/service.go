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
	Fill(ctx context.Context, width int, height int, url string) ([]byte, error)
}

type Request struct {
	height int
	width  int
	url    string
}

type service struct {
	logger          zerolog.Logger
	cache           cache.Cache
	imageDownloader ImageDownloader
	imageResizer    ImageResizer
}

// Fill todo переписать на dto если получится
func (s *service) Fill(ctx context.Context, width int, height int, url string) ([]byte, error) {
	//todo Проверить кеш На наличие оригинальной картинки

	img, err := s.imageDownloader.Download(ctx, url)
	if err != nil {
		s.logger.Err(err).Msg("не удалось скачать файл с удаленного сервера")
		return nil, err
	}

	//todo проверить кеш на наличие обрезанной картинки

	rImg, err := s.imageResizer.Resize(ctx, img, width, height)
	if err != nil {
		s.logger.Err(err).Msg("не удалось изменить размер изображения")
		return nil, err
	}

	//todo Положить в кеш
	return rImg, nil
}
