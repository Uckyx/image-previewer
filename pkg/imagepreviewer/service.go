package imagepreviewer

import (
	"sync"

	"github.com/Uckyx/image-previewer/pkg/cache"
	"github.com/rs/zerolog"
)

type Service interface {
	Resize(request *ResizeRequest) (*ResizeResponse, error)
}

type ResizeResponse struct {
	Img     []byte
	Headers map[string][]string
}

type service struct {
	logger          zerolog.Logger
	cache           cache.Cache
	imageDownloader ImageDownloader
	imageResizer    ImageResizer
	w               sync.WaitGroup
}

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

func (s *service) Resize(request *ResizeRequest) (*ResizeResponse, error) {
	resizedImgKey := s.cache.GenerateResizedImgKey(request.url, request.width, request.height)
	resizedImg, ok := s.cache.Get(resizedImgKey)

	if ok {
		return &ResizeResponse{resizedImg, nil}, nil
	}

	originalImgKey := s.cache.GenerateOriginalImgKey(request.url)
	originalImg, ok := s.cache.Get(originalImgKey)

	if ok {
		resizedImg, err := s.imageResizer.Resize(request.ctx, originalImg, request.width, request.height)
		if err != nil {
			return nil, err
		}

		s.cache.Set(resizedImgKey, resizedImg)

		return &ResizeResponse{resizedImg, nil}, nil
	}

	downloadResponse, err := s.imageDownloader.Download(request.ctx, request.url, request.headers)
	if err != nil {
		return nil, err
	}

	s.cache.Set(originalImgKey, downloadResponse.img)

	resizedImg, err = s.imageResizer.Resize(request.ctx, downloadResponse.img, request.width, request.height)
	if err != nil {
		return nil, err
	}

	s.cache.Set(resizedImgKey, resizedImg)

	return &ResizeResponse{resizedImg, downloadResponse.headers}, nil
}
