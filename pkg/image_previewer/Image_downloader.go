package image_previewer

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	ErrTimeout        = fmt.Errorf("timeout on download img")
	ErrRequest        = fmt.Errorf("generate request error")
	ErrUnknownImgType = fmt.Errorf("unknown file type uploaded")
	ErrResponseStatus = fmt.Errorf("response status code not 200")
)

type ImageDownloader interface {
	Download(ctx context.Context, imageUrl string) (imgResponse *DownloadResponse, err error)
}

type DownloadResponse struct {
	img     []byte
	headers map[string][]string
}

type imageDownloader struct {
	logger zerolog.Logger
}

func NewImageDownloader(logger zerolog.Logger) ImageDownloader {
	return &imageDownloader{
		logger: logger,
	}
}

func (i *imageDownloader) Download(ctx context.Context, imageUrl string) (imgResponse *DownloadResponse, err error) {
	req, err := http.NewRequest(http.MethodGet, imageUrl, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrResponseStatus
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			i.logger.Err(err).Msg(err.Error())
		}
	}(resp.Body)

	responseImg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = validateImageType(responseImg)
	if err != nil {
		return nil, err
	}

	return &DownloadResponse{responseImg, resp.Header}, nil
}

func validateImageType(img []byte) error {
	imgType := http.DetectContentType(img)
	switch imgType {
	case "image/jpeg", "image/jpg":
		return nil
	default:
		return ErrUnknownImgType
	}
}
