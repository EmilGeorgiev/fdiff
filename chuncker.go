package fdiff

import (
	"crypto/sha1"
	"fmt"

	"github.com/EmilGeorgiev/fdiff/rollinghash"
)

// ChunkConfig contains config information of the Chunker.
type ChunkConfig struct {
	// WindowSize is the number of bytes that are included in
	// the window that going to be rolling/shifted through the data.
	WindowSize uint64 `yaml:"window_size"`

	// MinSizeChunk point how much must be the minimum size of a Chunk.
	MinSizeChunk int `yaml:"min_size_chunk"`

	// MaxSizeChunk point how much must be the maximum size of a Chunk.
	MaxSizeChunk int `yaml:"max_size_chunk"`

	// FingerprintBreakPoint point when boundary of the chunks. When the
	// hash value of the bytes in window are equal to FingerprintBreakPoint
	// this means that the Chuncker should create a new chunk
	FingerprintBreakPoint uint64 `yaml:"fingerprint_break_point"`
}

// Chunker split data to chunks. It read data from a channel and
// split data to chunks based on the information in ChunkConfig.
// After a chunk is created the Chunker sends it through a channel.
// It is not responsible for storing and processing created chunks.
type Chunker struct {
	config ChunkConfig

	// newRollingHash is creating a new rolling hash
	newRollingHash func([]byte) rollinghash.Hash

	// bytesOfTheChunk contains current bytes that will be included in the next Chunk.
	bytesOfTheChunk []byte

	// bytes is a channel from which Chunker read data and split it to chunks.
	bytes <-chan byte

	// chunks is the channel through Chunker sends chunks when they are created.
	chunks chan<- Chunk

	// offset points from where a Chunk started.
	offset uint64
}

// NewChunker initialize and return *Chunker.
func NewChunker(new func([]byte) rollinghash.Hash, cfg ChunkConfig, b chan byte, ch chan Chunk) *Chunker {
	return &Chunker{
		config:         cfg,
		newRollingHash: new,
		bytes:          b,
		chunks:         ch,
	}
}

// Start a goroutine that listen for a new bytes that should be split in chunks.
// and send
func (ch *Chunker) Start() {
	go func() {
		var h rollinghash.Hash
		for b := range ch.bytes {
			ch.bytesOfTheChunk = append(ch.bytesOfTheChunk, b)
			if h == nil {
				if uint64(len(ch.bytesOfTheChunk)) < ch.config.WindowSize {
					// the number of bytes should be equal to the windows, then we
					// can calculate the hash/sign of the first window
					continue
				}
				h = ch.newRollingHash(ch.bytesOfTheChunk)
				if ch.shouldCreateAChunk(h) {
					ch.createChunk()
					continue
				}

			}
			h.Next(b)
			if ch.shouldCreateAChunk(h) {
				ch.createChunk()
			}
		}

		if len(ch.bytesOfTheChunk) > 0 {
			// if there are bytes that are still not send, a chunk
			// is created and send it. This will be the last chunk.
			ch.createChunk()
		}
		close(ch.chunks)
	}()
}

func (ch *Chunker) shouldCreateAChunk(h rollinghash.Hash) bool {
	if (h.Value() == ch.config.FingerprintBreakPoint) && (len(ch.bytesOfTheChunk) >= ch.config.MinSizeChunk) {
		return true
	}
	return len(ch.bytesOfTheChunk) >= ch.config.MaxSizeChunk
}

func (ch *Chunker) createChunk() {
	sum := fmt.Sprintf("%s", sha1.Sum(ch.bytesOfTheChunk))
	sign := fmt.Sprintf("%x", sum)
	c := Chunk{
		Offset:    ch.offset,
		Length:    uint64(len(ch.bytesOfTheChunk)),
		Data:      ch.bytesOfTheChunk,
		Signature: sign,
	}
	ch.chunks <- c
	ch.offset += uint64(len(ch.bytesOfTheChunk))
	// reset bytes of the chunk because next byte will be part of the next chunk
	ch.bytesOfTheChunk = []byte{}
}
