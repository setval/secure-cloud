package local

import (
	"github.com/DiscoreMe/SecureCloud/filebuffer"
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

func (l *Local) Upload(f *filebuffer.FileBuffer) error {
	fcr, err := os.Create(l.filepath(f))
	if err != nil {
		return err
	}
	defer fcr.Close()

	return f.Write(fcr)
}

func (l *Local) Download(f *filebuffer.FileBuffer) error {
	fh, err := os.Open(l.filepath(f))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "")
	}
	defer fh.Close()

	return f.Read(fh)
}

func (l *Local) filepath(f *filebuffer.FileBuffer) string {
	return path.Join(storage.Folder, f.ID.String())
}
