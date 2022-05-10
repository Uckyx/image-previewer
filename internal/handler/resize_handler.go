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
	vars := mux.Vars(r)

	request, err := h.createRequest(vars)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Err(err).Msg(err.Error())

		return
	}

	resizeResponse, err := h.svc.Resize(r.Context(), request.width, request.height, request.url)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
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
	r = &ResizeRequest{}

	if r.width, err = strconv.Atoi(vars["width"]); err != nil {
		return nil, ErrIsNumeric
	}

	if r.height, err = strconv.Atoi(vars["height"]); err != nil {
		return nil, ErrIsNumeric
	}

	parsedURL, err := url.Parse(vars["imageUrl"])
	if err != nil {
		return nil, err
	}

	r.url = parsedURL.String()
	if parsedURL.Scheme == "" {
		r.url = fmt.Sprintf("https://%s", parsedURL.String())
	}

	return r, nil
}
