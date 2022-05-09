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

func TestDownload_Positive(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		ctx     context.Context
		imgName string
	}{
		{
			name:    "success_download_img",
			ctx:     ctx,
			imgName: "_gopher_original_1024x504.jpg",
		},
	}

	logger := log.With().Logger()
	id := NewImageDownloader(logger)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotImg, err := id.Download(tt.ctx, ImageURL+tt.imgName)
			if err != nil {
				t.Errorf("Download() error = %v", err)
				return
			}

			wantImg := loadImage(tt.imgName)
			if !reflect.DeepEqual(gotImg, wantImg) {
				t.Errorf("Download() gotImg = %v, want %v", gotImg, wantImg)
			}
		})
	}
}

func TestDownload_Negative(t *testing.T) {
	ctx := context.Background()
	ctxWithTimeOut, closefn := context.WithTimeout(ctx, time.Microsecond*1)
	defer closefn()

	tests := []struct {
		name    string
		ctx     context.Context
		imgName string
		url     string
		err     error
	}{
		{
			name:    "timeout_case",
			ctx:     ctxWithTimeOut,
			imgName: "_gopher_original_1024x504.jpg",
			url:     ImageURL,
			err:     ErrTimeout,
		},
		{
			name:    "not_allowed_type_img_case",
			ctx:     ctx,
			imgName: "_gopher_original_1024x504.png",
			url:     ImageURL,
			err:     ErrUnknownImgType,
		},
		{
			name:    "bad_response_case",
			ctx:     ctx,
			imgName: "",
			url:     ImageURL,
			err:     ErrResponseStatus,
		},
	}

	logger := log.With().Logger()
	id := NewImageDownloader(logger)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
