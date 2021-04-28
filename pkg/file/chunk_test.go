package file_test

import (
	"github.com/DiscoreMe/SecureCloud/pkg/file"
	"io"
	"os"
	"testing"
)

func TestChunks(t *testing.T) {
	os.MkdirAll("testdata/outfile", os.ModePerm)
	defer os.RemoveAll("testdata/outfile")

	ri, wi := io.Pipe()
	ro, wo := io.Pipe()

	go func() {
		// encrypt simulation
		for {
			b := make([]byte, file.ChunkSize+file.TrashChunkSize)
			if _, err := ri.Read(b); err != nil {
				return
			}
			wo.Write(b[:file.ChunkSize])
		}
	}()

	testfile, err := os.Open("testdata/out.jpg")
	if err != nil {
		t.Fatal(err)
	}

	if err := file.ReaderToChunks(testfile, "testdata/outfile", wi, ro); err != nil {
		t.Fatal(err)
	}

	output, err := os.Create("testdata/out2.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer output.Close()

	if err := file.ChunksToWriter("testdata/outfile", wi, ro, output); err != nil {
		t.Fatal(err)
	}
}
