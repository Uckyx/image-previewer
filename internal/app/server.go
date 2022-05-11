package app

import (
	"context"
	"image-previewer/internal/handler"
	"image-previewer/pkg/cache"
	"image-previewer/pkg/imagepreviewer"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

type Server struct {
	svc    imagepreviewer.Service
	logger zerolog.Logger
	router *mux.Router
}

func NewServer(logger zerolog.Logger, cacheCapacity int) (*Server, error) {
	svc := imagepreviewer.NewApp(
		logger,
		cache.NewCache(cacheCapacity),
		imagepreviewer.NewImageDownloader(logger),
		imagepreviewer.NewImageResizer(logger),
	)

	srv := &Server{
		svc:    svc,
		logger: logger,
	}

	srv.createRoute()

	return srv, nil
}

func (s *Server) Listen(ctx context.Context) error {
	httpSrv := &http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      s.router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	return httpSrv.ListenAndServe()
}

func (s *Server) createRoute() {
	r := mux.NewRouter()
	handlers := handler.NewHandlers(s.logger, s.svc)

	r.HandleFunc("/resize/{width:[0-9]+}/{height:[0-9]+}/{imageURL:.*}", handlers.ResizeHandler)

	s.router = r
}
