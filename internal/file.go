package internal

import (
	"fmt"
	"github.com/DiscoreMe/SecureCloud/pkg/file"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type UploadFileResp struct {
	Key    string `json:"key"`
	FileID string `json:"file_id"`
	URL    string `json:"url"`
}

// UploadFile is handler for upload file
func (s *Server) UploadFile(c echo.Context) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return err
	}

	storParam := c.QueryParam("storage")
	if storParam == "" {
		storParam = "local"
	}
	stor, ok := s.stors[storParam]
	if !ok {
		return fmt.Errorf("storage not found")
	}

	f, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	var fl file.File
	if _, err := fl.WriteFromReader(f); err != nil {
		return err
	}
	if err := fl.Encrypt(); err != nil {
		return err
	}

	if err := stor.Upload(&fl); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, UploadFileResp{
		Key:    fl.Key,
		FileID: fl.ID.String(),
		URL:    fmt.Sprintf("%s/api/file/%s/%s?storage=%s", c.Request().Host, fl.ID.String(), fl.Key, storParam),
	})
}

func (s *Server) File(c echo.Context) error {
	key := c.Param("key")
	fileID, err := uuid.FromString(c.Param("fileID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var f = &file.File{
		ID:  fileID,
		Key: key,
	}

	storParam := c.QueryParam("storage")
	if storParam == "" {
		storParam = "local"
	}
	stor, ok := s.stors[storParam]
	if !ok {
		return fmt.Errorf("storage not found")
	}

	if err := stor.Download(f); err != nil {
		return err
	}

	if err := f.Decrypt(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "key is invalid")
	}

	return c.Stream(http.StatusOK, "", f)
}
