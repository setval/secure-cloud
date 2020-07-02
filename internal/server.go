package internal

import (
	"github.com/labstack/echo"
)

// Server is a implementation of the http handlers
type Server struct {
	e *echo.Echo
	ServerConfig
}

// ServerConfig is a config of the server
type ServerConfig struct {
	// ValidToken is a confirming all requests
	ValidToken string
}

func NewServer(cfg ServerConfig) *Server {
	return &Server{
		ServerConfig: cfg,
	}
}

func (s *Server) SetupAPI() {
	e := echo.New()
	apiGroup := e.Group("/api")

	apiGroup.POST("/upload", s.UploadFile, s.TokenValidator)
	apiGroup.GET("/file/:fileID/:key", s.File)

	s.e = e
}

func (s *Server) Listen(address string) error {
	return s.e.Start(address)
}
