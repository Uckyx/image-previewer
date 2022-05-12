package handler

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/Uckyx/image-previewer/pkg/imagepreviewer"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

var defaultImgURL = "https://raw.githubusercontent.com/Uckyx/image-previewer/master/img_example/"

func TestHandlers_ResizeHandler_Positive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := imagepreviewer.NewMockService(ctrl)
	logger := log.With().Logger()

	image := loadImage("_gopher_original_1024x504.jpg")

	tests := []struct {
		name           string
		width          int
		height         int
		url            string
		response       string
		resizeResponse *imagepreviewer.ResizeResponse
		responseCode   int
		img            []byte
	}{
		{
			name:           "ok_case",
			width:          500,
			height:         600,
			url:            defaultImgURL + "_gopher_original_1024x504.jpg",
			response:       string(image),
			resizeResponse: &imagepreviewer.ResizeResponse{Img: image},
			responseCode:   http.StatusOK,
			img:            image,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			req := httptest.NewRequest(http.MethodGet, "https://example.com", nil)
			req = mux.SetURLVars(req, map[string]string{
				"width":    strconv.Itoa(tt.width),
				"height":   strconv.Itoa(tt.height),
				"imageURL": tt.url,
			})

			request := imagepreviewer.NewResizeRequest(req.Context(), tt.width, tt.height, tt.url, req.Header)
			mockService.EXPECT().Resize(request).Return(tt.resizeResponse, nil)
			h := &Handlers{
				logger: logger,
				svc:    mockService,
			}

			w := httptest.NewRecorder()

			h.ResizeHandler(w, req)

			if status := w.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}

			if w.Body.String() != tt.response {
				t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), tt.response)
			}
		})
	}
}

func TestHandlers_ResizeHandler_Negative(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := imagepreviewer.NewMockService(ctrl)
	l := log.With().Logger()

	tests := []struct {
		name           string
		width          string
		height         string
		url            string
		response       string
		resizeResponse *imagepreviewer.ResizeResponse
		err            error
		httpStatus     int
	}{
		{
			name:       "bad_request_case",
			width:      "foo",
			height:     "bar",
			url:        defaultImgURL + "_gopher_original_1024x504.jpg",
			response:   "validation error",
			httpStatus: http.StatusBadRequest,
		},
		{
			name:           "bad_gateway_case",
			width:          "300",
			height:         "400",
			url:            defaultImgURL + "_gopher_original_1024x504.jpg",
			response:       "resize image failed",
			resizeResponse: nil,
			httpStatus:     http.StatusBadGateway,
			err:            errors.New("error"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "https://example.com", nil)
			req = mux.SetURLVars(req, map[string]string{
				"width":    tt.width,
				"height":   tt.height,
				"imageURL": tt.url,
			})

			if tt.resizeResponse != nil || tt.err != nil {
				width, err := strconv.Atoi(tt.width)
				if err != nil {
					t.Errorf("error converted width to int")
				}

				height, err := strconv.Atoi(tt.height)
				if err != nil {
					t.Errorf("error converted height to int")
				}

				request := imagepreviewer.NewResizeRequest(req.Context(), width, height, tt.url, req.Header)
				mockService.EXPECT().Resize(request).Return(tt.resizeResponse, tt.err)
			}

			h := &Handlers{
				logger: l,
				svc:    mockService,
			}

			w := httptest.NewRecorder()

			h.ResizeHandler(w, req)

			if status := w.Code; status == http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}

			if w.Body.String() != tt.response {
				t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), tt.response)
			}
		})
	}
}

func TestHandlers_ResizeHandler_ProxyHeaders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := imagepreviewer.NewMockService(ctrl)
	l := log.With().Logger()

	image1 := loadImage("gopher_256x126_resized.jpg")

	headers := map[string][]string{
		"Content-Length": {0: "6495"},
		"Content-Type":   {0: "image/jpeg"},
	}

	tests := []struct {
		name         string
		width        int64
		height       int64
		url          string
		response     string
		fillResponse *imagepreviewer.ResizeResponse
		err          error
		httpStatus   int64
	}{
		{
			name:         "good headers",
			width:        200,
			height:       300,
			url:          defaultImgURL + "_gopher_original_1024x504.jpg",
			response:     string(image1),
			fillResponse: &imagepreviewer.ResizeResponse{Img: image1, Headers: headers},
			httpStatus:   http.StatusOK,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "https://example.com", nil)
			req = mux.SetURLVars(req, map[string]string{
				"width":    strconv.Itoa(int(tt.width)),
				"height":   strconv.Itoa(int(tt.height)),
				"imageURL": tt.url,
			})

			fillParams := imagepreviewer.NewResizeRequest(
				req.Context(),
				int(tt.width),
				int(tt.height),
				tt.url, req.Header,
			)
			mockService.EXPECT().Resize(fillParams).Return(tt.fillResponse, tt.err)
			h := &Handlers{
				logger: l,
				svc:    mockService,
			}

			w := httptest.NewRecorder()
			h.ResizeHandler(w, req)

			for name, values := range tt.fillResponse.Headers {
				for _, value := range values {
					require.Equal(t, value, w.Header().Get(name))
				}
			}
		})
	}
}

func loadImage(imgName string) []byte {
	fileToBeUploaded := "../../img_example/" + imgName
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
