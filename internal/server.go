package internal

import (
	"github.com/DiscoreMe/SecureCloud/pkg/storage"
	"github.com/DiscoreMe/SecureCloud/pkg/storage/local"
	"github.com/labstack/echo"
)

// Server is a implementation of the http handlers
type Server struct {
	e     *echo.Echo
	stors map[string]storage.Storage
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
		stors:        make(map[string]storage.Storage),
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

func (s *Server) EnableLocalStorage() error {
	s.stors["local"] = local.New()
	return nil
}

func (s *Server) EnableS3Storage(endpoint, bucket, accessKey, secretKey, location string) error {
	//stor, err := s3.New(endpoint, bucket, accessKey, secretKey, location)
	//if err != nil {
	//	return err
	//}
	//s.stors["s3"] = stor
	return nil
}

func (s *Server) EnableDriveStorage() error {
	//stor, err := drive.New()
	//if err != nil {
	//	return err
	//}
	//s.stors["drive"] = stor
	return nil
}

func (s *Server) UpdateDriveToken() error {
	//d := drive.Drive{}
	//return d.UpdateToken()
	return nil
}
