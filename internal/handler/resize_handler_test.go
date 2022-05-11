package handler

import (
	"bufio"
	"fmt"
	imagepreviewermock "image-previewer/mocks/pkg/imagepreviewer"
	"image-previewer/pkg/imagepreviewer"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestHandlers_ResizeHandler_Positive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := imagepreviewermock.NewMockService(ctrl)
	logger := log.With().Logger()

	image := loadImage("_gopher_original_1024x504.jpg")

	tests := []struct {
		name           string
		width          int
		height         int
		url            string
		uri            string
		response       string
		resizeResponse *imagepreviewer.ResizeResponse
		responseCode   int
		img            []byte
	}{
		{
			name:           "status_ok",
			width:          500,
			height:         600,
			url:            "https://raw.githubusercontent.com/",
			uri:            "OtusGolang/final_project/master/examples/image-previewer_gopher_original_1024x504.jpg",
			response:       string(image),
			resizeResponse: &imagepreviewer.ResizeResponse{Img: image},
			responseCode:   http.StatusOK,
			img:            image,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "https://example.com", nil)
			req = mux.SetURLVars(req, map[string]string{
				"width":    strconv.Itoa(tt.width),
				"height":   strconv.Itoa(tt.height),
				"imageUrl": tt.url + tt.uri,
			})

			mockService.EXPECT().Resize(req.Context(), tt.width, tt.height, tt.url+tt.uri).Return(tt.resizeResponse, nil)
			h := &Handlers{
				logger: logger,
				svc:    mockService,
			}

			w := httptest.NewRecorder()
			resp := w.Result()
			defer resp.Body.Close()

			h.ResizeHandler(w, req)
			require.Equal(t, http.StatusOK, tt.responseCode)
			require.Equal(t, strings.TrimSpace(w.Body.String()), tt.response)
		})
	}
}

func TestHandlers_ResizeHandler_Negative(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := imagepreviewermock.NewMockService(ctrl)
	logger := log.With().Logger()

	image1 := loadImage("_gopher_original_1024x504.jpg")

	tests := []struct {
		name           string
		width          int
		height         int
		url            string
		uri            string
		response       string
		resizeResponse *imagepreviewer.ResizeResponse
		responseCode   int
		img            []byte
	}{
		{
			name:           "bad",
			width:          500,
			height:         600,
			url:            "https://raw.githubusercontent.com/",
			uri:            "OtusGolang/final_project/master/examples/image-previewer_gopher_original_1024x504.jpg",
			response:       string(image1),
			resizeResponse: &imagepreviewer.ResizeResponse{Img: image1},
			responseCode:   http.StatusOK,
			img:            image1,
		},
		{
			name:           "success",
			width:          500,
			height:         600,
			url:            "https://raw.githubusercontent.com/",
			uri:            "OtusGolang/final_project/master/examples/image-previewer_gopher_original_1024x504.jpg",
			response:       string(image1),
			resizeResponse: &imagepreviewer.ResizeResponse{Img: image1},
			responseCode:   http.StatusOK,
			img:            image1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "https://example.com", nil)
			req = mux.SetURLVars(req, map[string]string{
				"width":    strconv.Itoa(tt.width),
				"height":   strconv.Itoa(tt.height),
				"imageUrl": tt.url + tt.uri,
			})

			mockService.EXPECT().Resize(req.Context(), tt.width, tt.height, tt.url+tt.uri).Return(tt.resizeResponse, nil)
			h := &Handlers{
				logger: logger,
				svc:    mockService,
			}

			w := httptest.NewRecorder()
			resp := w.Result()
			defer resp.Body.Close()

			h.ResizeHandler(w, req)
			require.Equal(t, http.StatusOK, tt.responseCode)
			require.Equal(t, strings.TrimSpace(w.Body.String()), string(image1))
		})
	}
}

func loadImage(imgName string) []byte {
	fileToBeUploaded := "./img_example/" + imgName
	file, err := os.Open(fileToBeUploaded)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fileInfo, _ := file.Stat()
	bytes := make([]byte, fileInfo.Size())

	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	return bytes
}
