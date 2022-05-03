package image_previewer

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/rs/zerolog"
)

var (
	ErrTimeout = fmt.Errorf("не удалось скачать за отведенный таймаут")
	ErrRequest = fmt.Errorf("не удалось сформировать реквест")
)

type ImageDownloader interface {
	Download(ctx context.Context, imageUrl string) (img []byte, err error)
}

type imageDownloader struct {
	logger zerolog.Logger
}

func NewImageDownloader(logger zerolog.Logger) ImageDownloader {
	return &imageDownloader{
		logger: logger,
	}
}

func (i *imageDownloader) Download(ctx context.Context, imageUrl string) (img []byte, err error) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, imageUrl, nil)
	if err != nil {
		i.logger.Error().Msg(ErrRequest.Error())
		return nil, err
	}

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		if networkErr, ok := err.(net.Error); ok && networkErr.Timeout() {
			i.logger.Error().Msg(ErrTimeout.Error())
			return nil, ErrTimeout
		}

		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
