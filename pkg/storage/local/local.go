package local

import (
	"github.com/DiscoreMe/SecureCloud/pkg/file"
	"github.com/DiscoreMe/SecureCloud/pkg/storage"
	"github.com/labstack/echo"
	"net/http"
	"os"
	"path"
)

type Local struct{}

func New() *Local {
	os.MkdirAll(storage.Folder, os.ModePerm)
	return &Local{}
}

func (l *Local) Upload(f *file.File) error {
	fcr, err := os.Create(l.filepath(f))
	if err != nil {
		return err
	}
	defer fcr.Close()

	_, err = fcr.Write(f.Bytes())
	if err != nil {
		return err
	}
	return fcr.Close()
}

func (l *Local) Download(f *file.File) error {
	fh, err := os.Open(l.filepath(f))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "")
	}
	defer fh.Close()
	_, err = f.WriteFromReader(fh)
	return err
}

func (l *Local) filepath(f *file.File) string {
	return path.Join(storage.Folder, f.ID.String())
}
