package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
)

var ErrIsNumeric = fmt.Errorf("field must be number")

type ResizeRequest struct {
	width  int
	height int
	url    string
}

func (h *Handlers) ResizeHandler(w http.ResponseWriter, r *http.Request) {
	request, err := h.createRequest(mux.Vars(r))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("validation error"))
		h.logger.Err(err).Msg(err.Error())

		return
	}

	resizeResponse, err := h.svc.Resize(r.Context(), request.width, request.height, request.url, r.Header)
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

func (h *Handlers) createRequest(vars map[string]string) (r *ResizeRequest, err error) {
	width, err := strconv.Atoi(vars["width"])
	if err != nil {
		return nil, ErrIsNumeric
	}

	height, err := strconv.Atoi(vars["height"])
	if err != nil {
		return nil, ErrIsNumeric
	}

	imageURL, err := url.Parse(vars["imageURL"])
	if err != nil {
		return nil, err
	}

	imageURL.Scheme = "http"

	return &ResizeRequest{width: width, height: height, url: imageURL.String()}, nil
}
