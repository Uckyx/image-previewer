package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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
		h.logger.Err(err).Msg("не удалось сформировать реквест")

		return
	}

	imgResponse, err := h.svc.Resize(ctx, request.width, request.height, request.url)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		h.logger.Err(err).Msg(err.Error())

		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(imgResponse)))
	if _, err := w.Write(imgResponse); err != nil {
		w.WriteHeader(http.StatusBadGateway)
		h.logger.Err(err).Msg(err.Error())
	}
}

func (h *Handlers) createRequest(vars map[string]string, r *ResizeRequest) (err error) {
	if r.width, err = strconv.Atoi(vars["width"]); err != nil {
		return errors.New("поле width должно быть целочисленным")
	}
	if r.height, err = strconv.Atoi(vars["height"]); err != nil {
		return errors.New("поле width должно быть целочисленным")
	}

	parsedUrl, err := url.Parse(vars["imageUrl"])
	if err != nil {
		return errors.New("не корректный формат ссылки для скачивания картинки")
	}

	r.url = parsedUrl.String()
	if parsedUrl.Scheme == "" {
		r.url = fmt.Sprintf("https://%s", parsedUrl.String())
	}

	return nil
}
