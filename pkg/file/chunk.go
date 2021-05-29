package file

import (
	"io"
)

const (
	ChunkSize      = 32768 // 32KB
	TrashChunkSize = 4     // increase 4B after encrypt
)

type cryptFn func([]byte) []byte

func ReaderToChunks(r io.Reader, w io.Writer, encryptFn cryptFn) error {
	for {
		b := make([]byte, ChunkSize)
		_, err := r.Read(b)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		eb := encryptFn(b)
		if _, err := w.Write(eb); err != nil {
			return err
		}
	}

	return nil
}

func ChunksToWriter(r io.Reader, w io.WriteCloser, decryptFn cryptFn) error {
	defer w.Close()

	for {
		b := make([]byte, ChunkSize+TrashChunkSize)
		_, err := r.Read(b)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if _, err := w.Write(decryptFn(b)); err != nil {
			return err
		}
	}

	return nil
}
