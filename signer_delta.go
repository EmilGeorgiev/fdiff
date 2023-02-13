package fdiff

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

// SignerDelta contains methods for sign a file (Sign) anf
// find difference (FindDelta) between two version of files.
type SignerDelta interface {
	Sign(file, signatureFile string) error
	FindDelta(signatureFile, newFile string) (Delta, error)
}

// Delta contains the difference between two data bytes.
type Delta struct {
	// NewChunks contains new chunks that are missing in the old data bytes
	NewChunks []Chunk

	// OldChunks contains all chunks from the old data
	// bytes that are removed or updated and are not up-to-date.
	OldChunks []Chunk
}

// Chunk represent one chunk of the data bytes.
type Chunk struct {
	// Offset show from which byte the Chunk is started.
	Offset uint64

	// Data contains all bytes in the chunk.
	Data []byte

	// Length is the number of bytes in the chunk.
	Length uint64

	// Signature is the unique signature/hash of the data.
	// Two chunks with equal Data will have the same Signature.
	Signature string
}

// String return string representation of the chunk in the format <offset>-<length>-<signature>.
func (ch Chunk) String() string {
	return fmt.Sprintf("%d-%d-%s", ch.Offset, len(ch.Data), ch.Signature)
}

// createChunkFromString create a new chunk from a string. The parameter
// 'str' MUST contain a value in format <offset>-<length>-<signature>.
func createChunkFromString(str string) Chunk {
	p := strings.Split(str, "-")
	if len(p) != 3 {
		return Chunk{}
	}

	offset, err := strconv.ParseUint(p[0], 10, 64)
	if err != nil {

	}
	length, err := strconv.ParseUint(p[1], 10, 64)
	if err != nil {

	}

	return Chunk{
		Offset:    offset,
		Length:    length,
		Signature: p[2],
	}
}

type fileSignerDelta struct {
	chunks <-chan Chunk
	data   chan<- byte
}

// NewFileSignerDelta initialize and return a new SignerDelta.
func NewFileSignerDelta(d chan<- byte, ch <-chan Chunk) SignerDelta {
	return fileSignerDelta{
		chunks: ch,
		data:   d,
	}
}

// Sign create a new file that contains chunk's signatures of a file. The method
// read all data from a file and send bytes to the chunker worker. Then ged created
// chunks and store them to signatureFile.
func (fs fileSignerDelta) Sign(file, signatureFile string) error {
	f, err := os.Create(signatureFile)
	if err != nil {
		return err
	}

	if err = fs.sendFileDataToChunkerWorker(file); err != nil {
		return err
	}

	defer f.Close()
	for ch := range fs.chunks {
		_, _ = f.Write([]byte(ch.String() + "\n"))
	}
	return err
}

// FindDelta find the difference between old and new version of a file. The method accept two parameters,
// the first one, fileSignature, is the file that contains all chunks' signatures that are used to find
// difference in the new version of the file 'newFile'.
func (fio fileSignerDelta) FindDelta(fileSignature, newFile string) (Delta, error) {
	chunks := decodeChunksOfSignatureFile(fileSignature)
	if err := fio.sendFileDataToChunkerWorker(newFile); err != nil {
		return Delta{}, err
	}

	var newChunks []Chunk
	for ch := range fio.chunks {
		if _, ok := chunks[ch.Signature]; ok {
			delete(chunks, ch.Signature)
			continue
		}
		newChunks = append(newChunks, ch)
	}

	var oldChunks []Chunk
	for _, ch := range chunks {
		oldChunks = append(oldChunks, ch)
	}

	sort.Slice(oldChunks, func(i, j int) bool {
		return oldChunks[i].Offset < oldChunks[j].Offset
	})
	return Delta{NewChunks: newChunks, OldChunks: oldChunks}, nil
}

// sendFileDataToChunkerWorker sends the data in the file
// through a channel to worker that will split data to chunks.
//
// The parameter 'file' is the file of which should be split to chunks.
func (fs fileSignerDelta) sendFileDataToChunkerWorker(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	go func() {
		for {
			data := make([]byte, 48)
			n, errr := f.Read(data)
			if errr != nil {
				close(fs.data)
				if err != io.EOF {
					return
				}
				return
			}
			for _, b := range data[:n] {
				fs.data <- b
			}
		}
	}()
	return nil
}

func decodeChunksOfSignatureFile(f string) map[string]Chunk {
	chunks := map[string]Chunk{}
	file, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		ch := createChunkFromString(scanner.Text())
		chunks[ch.Signature] = ch
	}
	return chunks
}
