package file

import (
	"bufio"
	uuid "github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"os"
	"path"
)

const (
	ChunkSize      = 32768    // 32KB
	TrashChunkSize = 4 * 1024 // increase 4KB after encrypt
)

func ReaderToChunks(reader io.Reader, dstFile string, w *io.PipeWriter, r *io.PipeReader) error {
	if w == nil || r == nil {
		panic("pipe is nil")
	}
	if err := os.MkdirAll(dstFile, os.ModePerm); err != nil {
		return err
	}

	out, err := os.Create(path.Join(dstFile, "meta"))
	if err != nil {
		return err
	}
	defer out.Close()

	for {
		chunkID := uuid.NewV4().String()
		b := make([]byte, ChunkSize)
		_, err := reader.Read(b)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if _, err := w.Write(b); err != nil {
			return err
		}

		b = make([]byte, ChunkSize+TrashChunkSize)
		if _, err := r.Read(b); err != nil {
			return err
		}

		if _, err := out.WriteString(chunkID + "\n"); err != nil {
			return err
		}

		if err := os.WriteFile(path.Join(dstFile, chunkID), b, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func ChunksToWriter(srcDir string, w *io.PipeWriter, r *io.PipeReader, output io.Writer) error {
	metaFile, err := os.Open(path.Join(srcDir, "meta"))
	if err != nil {
		return err
	}
	defer metaFile.Close()

	scan := bufio.NewScanner(metaFile)
	for scan.Scan() {
		b, err := ioutil.ReadFile(path.Join(srcDir, scan.Text()))
		if err != nil {
			return err
		}
		if _, err := w.Write(b); err != nil {
			return err
		}

		if _, err := r.Read(b); err != nil {
			return err
		}

		if _, err := output.Write(b[:ChunkSize]); err != nil {
			return err
		}
	}

	if scan.Err() != nil {
		return scan.Err()
	}
	return nil
}
