package filebuffer

import (
	"bytes"
	"io"
	"io/ioutil"
)

// FileBuffer implements io.ReadWriteCloser interface.
// The feature of this structure, in contrast to the simple use of a buffer or passing os.
// File as the default parameter , is the ability to operate with old and new bytes of files
// and support for additional things.
type FileBuffer struct {
	b bytes.Buffer
}

func (f *FileBuffer) Read(p []byte) (n int, err error) {
	return f.b.Read(p)
}

func (f *FileBuffer) Write(p []byte) (n int, err error) {
	return f.b.Write(p)
}

func (f *FileBuffer) Reset() {
	f.b.Reset()
}

func (f *FileBuffer) Close() error {
	f.b.Reset()
	return nil
}

func (f *FileBuffer) WriteFromReader(r io.Reader) (n int, err error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, err
	}
	return f.Write(b)
}

func (f *FileBuffer) Bytes() []byte {
	return f.b.Bytes()
}
