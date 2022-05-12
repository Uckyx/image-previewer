package handler

import (
	"github.com/Uckyx/image-previewer/pkg/imagepreviewer"
	"github.com/rs/zerolog"
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
