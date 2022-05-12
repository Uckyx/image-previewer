package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Uckyx/image-previewer/pkg/imagepreviewer"
	"github.com/gorilla/mux"
)

var (
	ErrWidthIsNumeric  = fmt.Errorf("field width must be number")
	ErrHeightIsNumeric = fmt.Errorf("field height must be number")
)

func (h *Handlers) ResizeHandler(w http.ResponseWriter, r *http.Request) {
	request, err := h.createRequest(r.Context(), mux.Vars(r), r.Header)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("validation error"))
		h.logger.Err(err).Msg(err.Error())

		return
	}

	resizeResponse, err := h.svc.Resize(request)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("resize image failed"))
		h.logger.Err(err).Msg(err.Error())

		return
	}

	for name, values := range resizeResponse.Headers {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(resizeResponse.Img)))
	if _, err := w.Write(resizeResponse.Img); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Err(err).Msg(err.Error())
	}
}

func (h *Handlers) createRequest(
	ctx context.Context,
	vars map[string]string,
	headers map[string][]string,
) (r *imagepreviewer.ResizeRequest, err error) {
	width, err := strconv.Atoi(vars["width"])
	if err != nil {
		return nil, ErrWidthIsNumeric
	}

	height, err := strconv.Atoi(vars["height"])
	if err != nil {
		return nil, ErrHeightIsNumeric
	}

	imageURL, err := url.Parse(vars["imageURL"])
	if err != nil {
		return nil, err
	}

	imageURL.Scheme = "https"

	return imagepreviewer.NewResizeRequest(ctx, width, height, imageURL.String(), headers), nil
}
