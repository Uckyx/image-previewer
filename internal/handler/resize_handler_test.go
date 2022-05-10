package handler

import (
	"bufio"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"image-previewer/pkg/image_previewer"
	mock_image_previewer "image-previewer/pkg/image_previewer/mock"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestHandlers_ResizeHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_image_previewer.NewMockService(ctrl)
	logger := log.With().Logger()

	image1 := loadImage("_gopher_original_1024x504.jpg")

	tests := []struct {
		name           string
		width          int
		height         int
		url            string
		response       string
		resizeResponse *image_previewer.ResizeResponse
		img            []byte
	}{
		{
			name:           "success",
			width:          500,
			height:         600,
			url:            "https://raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer_gopher_original_1024x504.jpg",
			response:       string(image1),
			resizeResponse: &image_previewer.ResizeResponse{Img: image1},
			img:            image1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "https://example.com", nil)
			req = mux.SetURLVars(req, map[string]string{
				"width":    strconv.Itoa(tt.width),
				"height":   strconv.Itoa(tt.height),
				"imageUrl": tt.url,
			})

			mockService.EXPECT().Resize(req.Context(), tt.width, tt.height, tt.url).Return(tt.resizeResponse, nil)
			h := &Handlers{
				logger: logger,
				svc:    mockService,
			}

			w := httptest.NewRecorder()

			h.ResizeHandler(w, req)
			require.Equal(t, http.StatusOK, w.Result().StatusCode)
			require.Equal(t, strings.TrimSpace(w.Body.String()), string(image1))
			require.Equal(t, w.Header().Get("Content-Type"), "image/jpeg")
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
