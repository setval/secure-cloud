package file_test

import (
	"crypto/md5"
	"fmt"
	"github.com/DiscoreMe/SecureCloud/pkg/file"
	"github.com/magiconair/properties/assert"
	"os"
	"testing"
)

func TestChunks(t *testing.T) {
	os.MkdirAll("testdata/outfile", os.ModePerm)
	defer os.RemoveAll("testdata/outfile")

	testfile, err := os.Open("testdata/out.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer testfile.Close()

	chunksFile, err := os.Create("testdata/out.chunks")
	if err != nil {
		t.Fatal(err)
	}
	defer chunksFile.Close()

	output, err := os.Create("testdata/out2.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer output.Close()

	if err := file.ReaderToChunks(testfile, chunksFile, func(bytes []byte) []byte {
		b := make([]byte, file.TrashChunkSize)
		bytes = append(bytes, b...)
		return bytes
	}); err != nil {
		t.Fatal(err)
	}

	chunksFile.Close()
	chunksFile, _ = os.Open("testdata/out.chunks")

	if err := file.ChunksToWriter(chunksFile, output, func(bytes []byte) []byte {
		return bytes[:file.ChunkSize]
	}); err != nil {
		t.Fatal(err)
	}

	bytes1, err := os.ReadFile("testdata/out.jpg")
	if err != nil {
		t.Fatal(err)
	}
	bytes2, err := os.ReadFile("testdata/out2.jpg")
	if err != nil {
		t.Fatal(err)
	}

	hash1 := fmt.Sprintf("%x", md5.Sum(bytes1))
	hash2 := fmt.Sprintf("%x", md5.Sum(bytes2))
	assert.Equal(t, hash1, hash2)
	//
	//	if err := file.ChunksToWriter("testdata/outfile", wi, ro, output); err != nil {
	//		t.Fatal(err)
	//	}
}
