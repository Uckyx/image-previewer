package app

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"image-previewer/pkg/cache"
	"net"
	"net/http"
	"time"

	"image-previewer/pkg/imagePreviewer"
)

const CacheCapacity = 100

type Server struct {
	app    imagePreviewer.Service
	router *mux.Router
}

func NewServer() (*Server, error) {
	svc := imagePreviewer.NewApp(
		cache.NewCache(CacheCapacity),
	)

	srv := &Server{
		app: svc,
	}

	srv.createRoute()

	return srv, nil
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

func (s *Server) Fill(w http.ResponseWriter, r *http.Request) {
	err := s.app.Fill(w, r)
	if err != nil {
		return
	}
}

func (s *Server) createRoute() {
	r := mux.NewRouter()
	r.HandleFunc("/fill/{width:[0-9]+}/{height:[0-9]+}/{imageUrl:.*}", s.Fill)
	http.Handle("/", r)

	s.router = r
}
