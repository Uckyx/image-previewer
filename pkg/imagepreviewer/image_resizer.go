package imagepreviewer

import (
	"bytes"
	"context"
	"fmt"
	"image/jpeg"
	"io"
	"os"
	"time"

	"github.com/disintegration/imaging"
	"github.com/rs/zerolog"
)

func NewImageResizer(logger zerolog.Logger) ImageResizer {
	return &imageResizer{
		logger: logger,
	}
}

type ImageResizer interface {
	Resize(ctx context.Context, img []byte, width int, height int) (resizedImg []byte, err error)
}

type imageResizer struct {
	logger zerolog.Logger
}

func (ir *imageResizer) Resize(
	ctx context.Context,
	originalImg []byte,
	width int,
	height int,
) (resizedImg []byte, err error) {
	currentTimeStamp, err := fmt.Println(time.Now().Unix())
	if err != nil {
		return nil, err
	}

	imgName := fmt.Sprintf("image_%d_resized.jpg", currentTimeStamp)

	file, err := os.Create(imgName)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			ir.logger.Err(err).Msg(err.Error())
		}
	}(file)

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			ir.logger.Err(err).Msg(err.Error())
		}
	}(imgName)

	_, err = io.Copy(file, bytes.NewReader(originalImg))
	if err != nil {
		return nil, err
	}

	src, err := imaging.Open(imgName)
	if err != nil {
		return nil, err
	}

	img := imaging.Resize(src, width, height, imaging.Lanczos)

	imgBuffer := new(bytes.Buffer)
	err = jpeg.Encode(imgBuffer, img, nil)
	if err != nil {
		return nil, err
	}

	return imgBuffer.Bytes(), nil
}
