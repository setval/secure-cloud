package internal

import (
	"fmt"
	"github.com/DiscoreMe/SecureCloud/pkg/file"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"os"
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

	fcr, err := os.Create(fl.ID.String())
	if err != nil {
		return err
	}
	defer fcr.Close()

	_, err = fcr.Write(fl.Bytes())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, UploadFileResp{
		Key:    fl.Key,
		FileID: fl.ID.String(),
		URL:    fmt.Sprintf("%s/api/file/%s/%s", c.Request().Host, fl.ID.String(), fl.Key),
	})
}

func (s *Server) File(c echo.Context) error {
	key := c.Param("key")
	fileID, err := uuid.FromString(c.Param("fileID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	fh, err := os.Open(fileID.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "")
	}
	defer fh.Close()

	var f = file.File{
		ID:  fileID,
		Key: key,
	}

	if _, err := f.WriteFromReader(fh); err != nil {
		return err
	}
	if err := f.Decrypt(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "key is invalid")
	}

	return c.Stream(http.StatusOK, "", &f)
}
