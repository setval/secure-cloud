package file

import (
	"io"
)

const (
	ChunkSize      = 32768    // 32KB
	TrashChunkSize = 4 * 1024 // increase 4KB after encrypt
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

		if _, err := w.Write(encryptFn(b)); err != nil {
			return err
		}
	}

	return nil
}

func ChunksToWriter(r io.Reader, w io.Writer, decryptFn cryptFn) error {
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
