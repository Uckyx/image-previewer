package server

import (
	"Go/image-previewer/pkg/imagePreviewer"
	"net/http"
)

type Server struct {
	app imagePreviewer.Service
}

func NewServer(app imagePreviewer.Service) *Server {
	return &Server{app}
}

func (s *Server) GetResizedImage(response http.ResponseWriter, request *http.Request) {
}
