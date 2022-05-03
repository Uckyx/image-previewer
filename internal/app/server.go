package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"image-previewer/internal/handler"
	"image-previewer/pkg/cache"
	"image-previewer/pkg/image_previewer"
)

const CacheCapacity = 100

type Server struct {
	svc    image_previewer.Service
	logger zerolog.Logger
	router *mux.Router
}

func NewServer(logger zerolog.Logger) (*Server, error) {
	svc := image_previewer.NewApp(
		logger,
		cache.NewCache(CacheCapacity),
		image_previewer.NewImageDownloader(logger),
		image_previewer.NewImageResizer(logger),
	)

	srv := &Server{
		svc:    svc,
		logger: logger,
	}

	srv.createRoute()

	return srv, nil
}

func (s *Server) WithLogger(logger zerolog.Logger) {
	s.logger = logger
}

func (s *Server) Listen(ctx context.Context, port int) error {
	httpSrv := http.Server{
		Addr:         fmt.Sprintf(":%d", port),
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

	r.HandleFunc("/fill/{width:[0-9]+}/{height:[0-9]+}/{imageUrl:.*}", handlers.FillHandler)
	http.Handle("/", r)

	s.router = r
}
