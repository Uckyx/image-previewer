package handler

import (
	"github.com/rs/zerolog"
	"image-previewer/pkg/image_previewer"
)

type Handlers struct {
	logger zerolog.Logger
	svc    image_previewer.Service
}

func NewHandlers(
	logger zerolog.Logger,
	svc image_previewer.Service,
) *Handlers {
	return &Handlers{logger: logger, svc: svc}
}
