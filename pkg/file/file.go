package file

import (
	"github.com/satori/go.uuid"
	"io"
)

// File is a object stored in the cloud
type File struct {
	ID   uuid.UUID
	Key  string
	Body io.ReadCloser
}
