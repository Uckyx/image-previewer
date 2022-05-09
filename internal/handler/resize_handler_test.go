package handler

import (
	"bufio"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	mock_image_previewer "image-previewer/pkg/image_previewer/mock"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const ImageURL = "raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer"

func TestHandlers_ResizeHandler(t *testing.T) {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()
	mockService := mock_image_previewer.NewMockService(ctrl)
	logger := log.With().Logger()

	image1 := loadImage("./image_test/_gopher_original_1024x504.jpg")
	image2 := loadImage("./image_test/gopher_1024x252.jpg")

	tests := []struct {
		name     string
		width    string
		height   string
		url      string
		response string
		img      []byte
	}{
		{
			name:   "good",
			width:  "500",
			height: "600",
			url:    generateUrl("_gopher_original_1024x504.jpg"),
			//response: string(image1),
			img: image1,
		},
		{
			name:   "good",
			width:  "300",
			height: "200",
			url:    generateUrl("gopher_1024x252.jpg"),
			//response: string(image2),
			img: image2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			req = mux.SetURLVars(req, map[string]string{
				"width":    tt.width,
				"height":   tt.height,
				"imageUrl": tt.url,
			})

			mockService.EXPECT().Resize(req.Context(), tt.width, tt.height, tt.url).Return(tt.img, nil)
			h := &Handlers{
				logger: logger,
				svc:    mockService,
			}

			w := httptest.NewRecorder()

			h.ResizeHandler(w, req)
			require.Equal(t, http.StatusOK, w.Result().StatusCode)
			//require.Equal(t, strings.TrimSpace(w.Body.String()), tt.response)
			require.Equal(t, w.Header().Get("Content-Type"), "image/jpeg")
		})
	}
}

func loadImage(filePath string) []byte {
	fileToBeUploaded := filePath
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

func generateUrl(uri string) string {
	return fmt.Sprintf("%s/%s", ImageURL, uri)
}
