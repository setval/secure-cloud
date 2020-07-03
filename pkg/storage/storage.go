package storage

import (
	"github.com/DiscoreMe/SecureCloud/pkg/file"
)

// Name of the folder where encrypted files will be uploaded.
const Folder = "secure-cloud"

// Storage is an interface that allows you
// to create a logical structure to support new cloud services.
// Upload is used for uploading files in the cloud storage. The name of the file that
// you upload must be in the form of ID is uuid.
// Download is used for downloading files. Data about the file must be saved
// to the same structure using the Write or WriteFromReader method.
//
// A pointer to the filebuffer.Filebuffer structure is passed to both functions
type Storage interface {
	Upload(*file.File) error
	Download(*file.File) error
}
