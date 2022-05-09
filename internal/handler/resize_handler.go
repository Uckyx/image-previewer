package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var (
	ErrIsNumeric = fmt.Errorf("field must be number")
)

type ResizeRequest struct {
	width  int
	height int
	url    string
}

func (h *Handlers) ResizeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	request := &ResizeRequest{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	err := h.createRequest(vars, request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Err(err).Msg(err.Error())

		return
	}

	resizeResponse, err := h.svc.Resize(ctx, request.width, request.height, request.url)
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
		w.WriteHeader(http.StatusBadGateway)
		h.logger.Err(err).Msg(err.Error())
	}
}

func (h *Handlers) createRequest(vars map[string]string, r *ResizeRequest) (err error) {
	if r.width, err = strconv.Atoi(vars["width"]); err != nil {
		return ErrIsNumeric
	}
	if r.height, err = strconv.Atoi(vars["height"]); err != nil {
		return ErrIsNumeric
	}

	parsedUrl, err := url.Parse(vars["imageUrl"])
	if err != nil {
		return err
	}

	r.url = parsedUrl.String()
	if parsedUrl.Scheme == "" {
		r.url = fmt.Sprintf("https://%s", parsedUrl.String())
	}

	return nil
}
