package imagepreviewer

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog"
)

var (
	ErrTimeout        = fmt.Errorf("timeout on download img")
	ErrUnknownImgType = fmt.Errorf("unknown file type uploaded")
	ErrResponseStatus = fmt.Errorf("response status code not 200")
)

type ImageDownloader interface {
	Download(
		ctx context.Context,
		imageURL string,
		headers map[string][]string,
	) (imgResponse *DownloadResponse, err error)
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

func (i *imageDownloader) Download(
	ctx context.Context,
	imageURL string,
	headers map[string][]string,
) (imgResponse *DownloadResponse, err error) {
	req, err := http.NewRequest(http.MethodGet, imageURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header = headers

	client := http.Client{}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrResponseStatus
	}

	responseImg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = validateImageType(responseImg)
	if err != nil {
		return nil, err
	}

	if err := resp.Body.Close(); err != nil {
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
