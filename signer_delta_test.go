package fdiff_test

import (
	"crypto/sha1"
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/EmilGeorgiev/fdiff"
	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	// SetUp
	d := make(chan byte, 100)
	ch := make(chan fdiff.Chunk, 100)
	fs := fdiff.NewFileSignerDelta(d, ch)
	fch := fakeChunker{data: d, chunks: ch, windowsSize: 48}
	fch.Start("./test/expected_sign_test_data")
	defer os.Remove("./test/expected_sign_test_data")
	defer os.Remove("./test/sign_test_data")

	// Action
	err := fs.Sign("./test/test_data", "./test/sign_test_data")

	// Assert
	assert.Nil(t, err)
	assert.True(t, equalFileContent("./test/expected_sign_test_data", "./test/sign_test_data"))
}

func TestFindDelta(t *testing.T) {
	// SetUp
	defer os.Remove("file1")
	defer os.Remove("file2")
	defer os.Remove("sign_file")
	oldFileData := []byte("The Low Bandwidth Network Filesystem (LBFS) from MIT uses Rabin " +
		"fingerprints to implement variable size shift-resistant blocks.")
	newFileData := []byte("The Low Bandwidth Network Filesystem (LBFS) from MIT uses Rabin " +
		"fingerprints to (THIS IS A NEW DATA)implement variable size shift-resistant blocks.")
	writeDataToFile("file1", oldFileData)
	signFile("file1", "sign_file", 30)
	// update file data
	writeDataToFile("file2", newFileData)

	d := make(chan byte, 100)
	ch := make(chan fdiff.Chunk, 100)
	fs := fdiff.NewFileSignerDelta(d, ch)
	fch := fakeChunker{data: d, chunks: ch, windowsSize: 30}
	fch.Start("")

	// Action
	actual, err := fs.FindDelta("sign_file", "file2")

	// Assert
	expected := fdiff.Delta{
		NewChunks: []fdiff.Chunk{
			{
				Offset:    60,
				Data:      []byte("bin fingerprints to (THIS IS A"),
				Length:    30,
				Signature: "b80cb62f9823ff1143099a32a6f46f7798a6b92d",
			},
			{
				Offset:    90,
				Data:      []byte(" NEW DATA)implement variable s"),
				Length:    30,
				Signature: "9ff86b8abc189759d45caab3a1c13704b12f63fe",
			},
			{
				Offset:    120,
				Data:      []byte("ize shift-resistant blocks."),
				Length:    27,
				Signature: "521400e2cf500bb9f745807e8f62e047d566f7d8",
			},
		},
		OldChunks: []fdiff.Chunk{
			{
				Offset:    60,
				Length:    30,
				Signature: "9d23da68e8d2e7b42b1e021b1a4c2912a827f285",
			},
			{
				Offset:    90,
				Length:    30,
				Signature: "301054dadc4095e21ea59492f52ae2518c9d4195",
			},
			{
				Offset:    120,
				Length:    7,
				Signature: "98d34d28921ae7f4c29bcfc1ddd4a87b1dcaf455",
			},
		},
	}
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func signFile(file, filesSign string, windowsSize int) {
	d := make(chan byte, 100)
	ch := make(chan fdiff.Chunk, 100)
	fs := fdiff.NewFileSignerDelta(d, ch)
	defer os.Remove("./test/expected_sign_test_data")
	fch := fakeChunker{data: d, chunks: ch, windowsSize: windowsSize}
	fch.Start("./test/expected_sign_test_data")
	fs.Sign(file, filesSign)
}

func writeDataToFile(file string, data []byte) {
	f, err := os.Open(file)
	if err != nil {
		if _, ok := err.(*fs.PathError); ok {
			f, _ = os.Create(file)
		}
	}
	n, _ := f.Write(data)
	fmt.Println("Number of bytes: ", n)
}

func equalFileContent(expected, actual string) bool {
	expectedData, _ := os.ReadFile(expected)
	actualData, _ := os.ReadFile(actual)

	return sha1.Sum(expectedData) == sha1.Sum(actualData)
}

// fakeChunker create a chunks and send it back to the signer
type fakeChunker struct {
	chunks      chan<- fdiff.Chunk
	data        <-chan byte
	windowsSize int
}

func (fch fakeChunker) Start(file string) {
	var f *os.File
	if file != "" {
		f, _ = os.Create(file)
	}

	go func() {
		var offset uint64
		var bytesOfTheChunk []byte
		for {
			b, ok := <-fch.data
			if !ok {
				if len(bytesOfTheChunk) > 0 {
					h := fmt.Sprintf("%x", sha1.Sum(bytesOfTheChunk))
					chunk := fdiff.Chunk{
						Offset:    offset,
						Data:      bytesOfTheChunk,
						Length:    uint64(len(bytesOfTheChunk)),
						Signature: h,
					}
					fmt.Printf("signature is %s, data is: %v: \n", h, bytesOfTheChunk)
					fch.chunks <- chunk
					if f != nil {
						_, _ = f.Write([]byte(chunk.String() + "\n"))
					}
				}
				close(fch.chunks)
				return
			}
			//fmt.Println("Read byte: ", b)
			bytesOfTheChunk = append(bytesOfTheChunk, b)
			if len(bytesOfTheChunk) < fch.windowsSize {
				continue
			}

			h := fmt.Sprintf("%x", sha1.Sum(bytesOfTheChunk))
			chunk := fdiff.Chunk{
				Offset:    offset,
				Data:      bytesOfTheChunk,
				Length:    uint64(len(bytesOfTheChunk)),
				Signature: h,
			}
			fmt.Printf("signature is %s, data is: %v: \n", h, bytesOfTheChunk)
			fch.chunks <- chunk

			if f != nil {
				_, _ = f.Write([]byte(chunk.String() + "\n"))
			}

			offset += uint64(len(bytesOfTheChunk))
			bytesOfTheChunk = []byte{}
		}
	}()
}
