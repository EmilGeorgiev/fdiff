package fdiff_test

import (
	"crypto/sha1"
	"fmt"
	"os"
	"testing"

	"github.com/EmilGeorgiev/fdiff"
	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	// SetUp
	d := make(chan []byte, 100)
	ch := make(chan fdiff.Chunk, 100)
	fs := fdiff.NewFileSigner(d, ch)
	fch := fakeChunker{data: d, chunks: ch}
	fch.Start("./test/expected_sign_test_data")

	// Action
	err := fs.Sign("./test/test_data", "./test/sign_test_data")

	// Assert
	assert.Nil(t, err)
	assert.True(t, equalFileContent("./test/expected_sign_test_data", "./test/sign_test_data"))
}

func equalFileContent(expected, actual string) bool {
	expectedData, _ := os.ReadFile(expected)
	actualData, _ := os.ReadFile(actual)

	return sha1.Sum(expectedData) == sha1.Sum(actualData)
}

// fakeChunker create a chunks and send it back to the signer
type fakeChunker struct {
	chunks chan<- fdiff.Chunk
	data   <-chan []byte
}

func (fch fakeChunker) Start(file string) {
	f, _ := os.Create(file)

	go func() {
		var offset uint64
		for {
			b, ok := <-fch.data
			if !ok {
				close(fch.chunks)
				return
			}

			fch.chunks <- fdiff.Chunk{
				Offset: offset,
				Data:   b,
			}
			h := fmt.Sprintf("%s", sha1.Sum(b))
			str := fmt.Sprintf("%x\n", h)
			_, _ = f.Write([]byte(str))
			offset += uint64(len(b))
		}
	}()
}
