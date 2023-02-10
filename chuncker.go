package fdiff

import (
	"crypto/sha1"
	"fmt"

	"github.com/EmilGeorgiev/fdiff/rollinghash"
)

type Chunker struct {
	newRollingHash        func([]byte) rollinghash.Hash
	windowSize            uint64
	minSizeChunk          uint64
	maxSizeChunk          uint64
	fingerprintBreakPoint uint64
	bytesOfTheChunk       []byte
	bytes                 <-chan byte
	chunks                chan<- Chunk
	offset                uint64
}

func NewChunker(new func([]byte) rollinghash.Hash, ws, minSch, maxSch, breckpoint uint64, b chan byte, ch chan Chunk) *Chunker {
	return &Chunker{
		newRollingHash:        new,
		windowSize:            ws,
		minSizeChunk:          minSch,
		maxSizeChunk:          maxSch,
		fingerprintBreakPoint: breckpoint,
		bytes:                 b,
		chunks:                ch,
	}
}

func (ch *Chunker) Start() {
	go func() {
		var h rollinghash.Hash
		for b := range ch.bytes {
			ch.bytesOfTheChunk = append(ch.bytesOfTheChunk, b)
			if h == nil {
				if uint64(len(ch.bytesOfTheChunk)) < ch.windowSize {
					// the number of bytes should be equal to the windows, then we
					// can calculate the hash/sign of the first window
					continue
				}
				h = ch.newRollingHash(ch.bytesOfTheChunk)
				ch.tryToCreateChunks(h)
				continue
			}
			h.Next(b)
			ch.tryToCreateChunks(h)
		}
		close(ch.chunks)
	}()
}

func (ch *Chunker) tryToCreateChunks(h rollinghash.Hash) {
	if h.Value() == ch.fingerprintBreakPoint {
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
		ch.bytesOfTheChunk = []byte{} // reset bytes of the chunk because next byte will be part of the next chunk
	}
}
