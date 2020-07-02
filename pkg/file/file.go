package file

import (
	"github.com/DiscoreMe/SecureCloud/filebuffer"
	"github.com/satori/go.uuid"
)

// File is a object stored in the cloud
type File struct {
	ID  uuid.UUID
	Key string
	filebuffer.FileBuffer
}
