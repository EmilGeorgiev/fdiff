package fdiff_test

import (
	"crypto/sha1"
	"fmt"
	"testing"

	"github.com/EmilGeorgiev/fdiff"
	"github.com/EmilGeorgiev/fdiff/rollinghash"
	"github.com/stretchr/testify/assert"
)

func TestNewChunker(t *testing.T) {
	// SetUp
	data := []byte("If you abcd want to draw abcd readers to a story, you need abcd" +
		" to make them want to choose abc ABCD it. Hello World!!! abc d")
	b := make(chan byte)
	ch := make(chan fdiff.Chunk, 1)
	cfg := fdiff.ChunkConfig{
		WindowSize:            4,
		MinSizeChunk:          20,
		MaxSizeChunk:          50,
		FingerprintBreakPoint: 3194, // this is the hash fingerprint of "abcd"
	}
	c := fdiff.NewChunker(rollinghash.NewRabinFingerprint, cfg, b, ch)

	// Action
	c.Start()
	go func() {
		for _, d := range data {
			b <- d
		}
		close(b)
	}()

	// Assert
	var actual []fdiff.Chunk
	for chunk := range ch {
		actual = append(actual, chunk)
	}
	expected := []fdiff.Chunk{
		{
			Offset:    0,
			Data:      []byte("If you abcd want to draw abcd"),
			Length:    29,
			Signature: fmt.Sprintf("%x", sha1.Sum([]byte("If you abcd want to draw abcd"))),
		},
		{
			Offset:    29,
			Data:      []byte(" readers to a story, you need abcd"),
			Length:    34,
			Signature: fmt.Sprintf("%x", sha1.Sum([]byte(" readers to a story, you need abcd"))),
		},
		{
			Offset:    63,
			Data:      []byte(" to make them want to choose abc ABCD it. Hello Wo"),
			Length:    50,
			Signature: fmt.Sprintf("%x", sha1.Sum([]byte(" to make them want to choose abc ABCD it. Hello Wo"))),
		},
		{
			Offset:    113,
			Data:      []byte("rld!!! abc d"),
			Length:    12,
			Signature: fmt.Sprintf("%x", sha1.Sum([]byte("rld!!! abc d"))),
		},
	}

	assert.Equal(t, expected, actual)
}

func TestNewChunker_WhenTheFirstWindowIsTheFirstChunk(t *testing.T) {
	// SetUp
	data := []byte("If you want to draw ")
	b := make(chan byte)
	ch := make(chan fdiff.Chunk, 1)
	cfg := fdiff.ChunkConfig{
		WindowSize:            20,
		MinSizeChunk:          20,
		MaxSizeChunk:          50,
		FingerprintBreakPoint: 2245, // this is the hash fingerprint of "If you want to draw "
	}
	c := fdiff.NewChunker(rollinghash.NewRabinFingerprint, cfg, b, ch)

	// Action
	c.Start()
	go func() {
		for _, d := range data {
			b <- d
		}
		close(b)
	}()

	// Assert
	var actual []fdiff.Chunk
	for chunk := range ch {
		actual = append(actual, chunk)
	}
	expected := []fdiff.Chunk{
		{
			Offset:    0,
			Data:      []byte("If you want to draw "),
			Length:    20,
			Signature: fmt.Sprintf("%x", sha1.Sum([]byte("If you want to draw "))),
		},
	}

	assert.Equal(t, expected, actual)
}
