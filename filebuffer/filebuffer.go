package filebuffer

import (
	"github.com/DiscoreMe/SecureCloud/pkg/file"
	uuid "github.com/satori/go.uuid"
	"io"
)

// FileBuffer implements io.ReadWriteCloser interface.
// The feature of this structure, in contrast to the simple use of a buffer or passing os.
// File as the default parameter , is the ability to operate with old and new bytes of files
// and support for additional things.
type FileBuffer struct {
	ID  uuid.UUID
	Key string
	r   io.Reader
	w   io.WriteCloser
}

func New(r io.Reader) (*FileBuffer, error) {
	id := uuid.NewV4()
	key, err := file.GenKey(id)
	return &FileBuffer{
		r:   r,
		ID:  id,
		Key: key,
	}, err
}

func (f *FileBuffer) SetPipeWriter(w *io.PipeWriter) {
	f.w = w
}

func (f *FileBuffer) Write(w io.Writer) error {
	return file.ReaderToChunks(f.r, w, func(bytes []byte) []byte {
		return file.Encrypt(f.Key, bytes)
	})
}

func (f *FileBuffer) Read(r io.Reader) error {
	return file.ChunksToWriter(r, f.w, func(bytes []byte) []byte {
		return file.Decrypt(f.Key, bytes)
	})
}
