package image_previewer

import (
	"bufio"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"os"
	"reflect"
	"testing"
	"time"
)

const ImageURL = "https://raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/"

func TestDefaultImageDownloader_DownloadByUrl_Positive(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		ctx     context.Context
		imgName string
	}{
		{
			ctx:     ctx,
			imgName: "gopher_256x126.jpg",
		},
		{
			ctx:     ctx,
			imgName: "gopher_1024x252.jpg",
		},
	}

	logger := log.With().Logger()
	id := NewImageDownloader(logger)

	for _, tt := range tests {
		t.Run(tt.imgName, func(t *testing.T) {
			gotImg, err := id.Download(tt.ctx, ImageURL+tt.imgName)
			if err != nil {
				t.Errorf("DownloadByUrl() error = %v", err)
				return
			}

			wantImg := loadImage(tt.imgName)
			if !reflect.DeepEqual(gotImg, wantImg) {
				t.Errorf("DownloadByUrl() gotImg = %v, want %v", gotImg, wantImg)
			}
		})
	}
}

func TestDefaultImageDownloader_DownloadByUrl_Negative(t *testing.T) {
	ctx := context.Background()
	ctxWithTimeOut, closefn := context.WithTimeout(ctx, time.Microsecond*1)
	defer closefn()

	tests := []struct {
		ctx     context.Context
		imgName string
		url     string
		err     error
	}{
		{
			ctx:     ctxWithTimeOut,
			imgName: "gopher_200x700.jpg",
			url:     ImageURL,
			err:     ErrTimeout,
		},
		{
			ctx:     ctxWithTimeOut,
			imgName: "gopher_1024x252.jpg",
			url:     ImageURL,
			err:     ErrTimeout,
		},
		{
			ctx:     ctx,
			imgName: "gopher_200x700.png",
			url:     ImageURL,
			err:     ErrUnknownImgType,
		},
		{
			ctx:     ctx,
			imgName: "",
			url:     ImageURL,
			err:     ErrResponseStatus,
		},
	}

	logger := log.With().Logger()
	id := NewImageDownloader(logger)

	for _, tt := range tests {
		t.Run(tt.imgName, func(t *testing.T) {
			_, err := id.Download(tt.ctx, tt.url+tt.imgName)
			require.Errorf(t, err, tt.err.Error())
		})
	}
}

func loadImage(imgName string) []byte {
	fileToBeUploaded := "./image_test/" + imgName
	file, err := os.Open(fileToBeUploaded)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	bytes := make([]byte, fileInfo.Size())

	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return bytes
}
