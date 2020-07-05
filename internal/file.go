package internal

import (
	"fmt"
	"github.com/DiscoreMe/SecureCloud/pkg/file"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"net/url"
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

	var params url.URL
	params.Scheme = "http"
	params.Host = c.Request().Host
	params.Path = fmt.Sprintf("api/file/%s/%s", fl.ID.String(), fl.Key)
	values := params.Query()
	values.Add("storage", storParam)
	values.Add("name", fileHeader.Filename)
	values.Add("type", fileHeader.Header.Get("Content-Type"))
	params.RawQuery = values.Encode()

	return c.JSON(http.StatusCreated, UploadFileResp{
		Key:    fl.Key,
		FileID: fl.ID.String(),
		URL:    params.String(),
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

	params := c.Request().URL.Query()

	if params.Get("storage") == "" {
		params.Set("local", "local")
	}
	stor, ok := s.stors[params.Get("storage")]
	if !ok {
		return fmt.Errorf("storage not found")
	}

	if err := stor.Download(f); err != nil {
		return err
	}

	if err := f.Decrypt(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "key is invalid")
	}

	c.Response().Header().Set("Cache-Control", "immutable")
	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("%s; filename=%q", "inline", params.Get("name")))
	return c.Stream(http.StatusOK, params.Get("type"), f)
}
