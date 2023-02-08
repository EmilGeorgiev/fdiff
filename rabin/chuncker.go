package rabin

import (
	"github.com/EmilGeorgiev/fdiff"
)

type Chunker struct {
	fingerprintHash       FingerprintHash
	windowSize            uint64
	minSizeChunk          uint64
	maxSizeChunk          uint64
	fingerprintBreakPoint uint64
	bytes                 <-chan []byte
	chunks                chan<- fdiff.Chunk
	numberOfReadBytes     uint64
	offset                uint64
}

func NewChunker(ws, minSch, maxSch, breckpoint uint64, b chan []byte, ch chan fdiff.Chunk) Chunker {
	return Chunker{
		windowSize:            ws,
		minSizeChunk:          minSch,
		maxSizeChunk:          maxSch,
		fingerprintBreakPoint: breckpoint,
		bytes:                 b,
		chunks:                ch,
	}
}

func (ch Chunker) SplitDataToChunks() {
	go func() {
		var h FingerprintHash
		for b := range ch.bytes {
			if h == nil {
				ch.numberOfReadBytes = ch.windowSize
				h = NewHash(b[:ch.windowSize])
				if (h.Value() == ch.fingerprintBreakPoint) && (ch.windowSize >= ch.minSizeChunk) {
					c := fdiff.Chunk{
						Offset: ch.offset,
						Data:   b[:ch.windowSize],
					}
					ch.chunks <- c
					ch.offset = ch.windowSize
				}
				ch.tryToCreateChunks(h, b[ch.windowSize:])
				continue
			}

			ch.tryToCreateChunks(h, b)
		}
	}()
}

func (ch Chunker) tryToCreateChunks(h FingerprintHash, bytes []byte) {
	for _, b := range bytes {
		ch.offset++
		v := h.Next(b)
		if v == ch.fingerprintBreakPoint {
			// create a chunk
		}
	}
}
