package image_previewer

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
	currentTimeStamp, err := fmt.Println(time.Now().Format(time.RFC850))
	if err != nil {
		return nil, err
	}

	rImgName := fmt.Sprintf("image_%d_resized.jpg", currentTimeStamp)

	file, err := os.Create(rImgName)
	if err != nil {
		ir.logger.Err(err).Msg(err.Error())
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
	}(rImgName)

	_, err = io.Copy(file, bytes.NewReader(originalImg))
	if err != nil {
		ir.logger.Err(err).Msg(err.Error())
	}

	src, err := imaging.Open(rImgName)
	if err != nil {
		return nil, err
	}

	rImg := imaging.Resize(src, width, height, imaging.Lanczos)
	if err != nil {
		return nil, err
	}

	imgBuffer := new(bytes.Buffer)
	err = jpeg.Encode(imgBuffer, rImg, nil)

	return imgBuffer.Bytes(), nil
}
