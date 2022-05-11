package handler

import (
	"github.com/rs/zerolog"
	"image-previewer/pkg/imagepreviewer"
)

type Handlers struct {
	logger zerolog.Logger
	svc    imagepreviewer.Service
}

func NewHandlers(
	logger zerolog.Logger,
	svc imagepreviewer.Service,
) *Handlers {
	return &Handlers{logger: logger, svc: svc}
}
