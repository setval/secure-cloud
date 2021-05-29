package file_test

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/DiscoreMe/SecureCloud/pkg/file"
	"github.com/magiconair/properties/assert"
	uuid "github.com/satori/go.uuid"
	"github.com/xxtea/xxtea-go/xxtea"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestChunks(t *testing.T) {
	defer os.Remove("testdata/out.chunks")
	defer os.Remove("testdata/out2.jpg")

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
}

func TestXXTEA(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	fill := func(size int) (b []byte) {
		for i := 0; i < size; i++ {
			b = append(b, byte(rand.Intn(127)))
		}
		return
	}

	key := []byte(uuid.NewV4().String() + uuid.NewV4().String())

	for i := 0; i < 1_000; i++ {
		b := fill(int(32 + rand.Int31n(50)))
		eb := xxtea.Encrypt(b, key)
		b2 := xxtea.Decrypt(eb, key)

		if !bytes.Equal(b, b2) {
			t.Fatalf("b != b2; idx: %d", i)
		}
	}
}

func BenchmarkEncryptChunks(b *testing.B) {
	key := "123456789"

	testFile, err := os.Open("../../testdata/1.jpg")
	if err != nil {
		b.Fatal(err)
	}
	defer testFile.Close()

	buffer := bytes.Buffer{}

	for i := 0; i < b.N; i++ {
		if err := file.ReaderToChunks(testFile, &buffer, func(i []byte) []byte {
			return file.Encrypt(key, i)
		}); err != nil {
			b.Fatal(err)
		}

		buffer.Reset()
	}
}

func BenchmarkEncryptFullFile(b *testing.B) {
	key := "123456789"

	bb, err := ioutil.ReadFile("../../testdata/1.jpg")
	if err != nil {
		b.Fatal(err)
	}

	buffer := bytes.Buffer{}

	for i := 0; i < b.N; i++ {
		eb := file.Encrypt(key, bb)
		_, err := buffer.Write(eb)
		if err != nil {
			b.Fatal(err)
		}
		buffer.Reset()
	}
}
