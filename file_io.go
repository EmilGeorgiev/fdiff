package fdiff

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type SignerDelta interface {
	Sign(file, signatureFile string) error
	FindDelta(signatureFile, newFile string) Delta
}

type Delta struct {
	NewChunks []Chunk
	OldChunks []Chunk
}

type Chunk struct {
	Offset    uint64
	Data      []byte
	Length    uint64
	Signature string
}

func (ch Chunk) String() string {
	return fmt.Sprintf("%d-%d-%s", ch.Offset, len(ch.Data), ch.Signature)
}

func CreateChunkFromString(str string) Chunk {
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

type FileIO struct {
	chunks <-chan Chunk
	data   chan<- byte
}

func NewFileIO(d chan<- byte, ch <-chan Chunk) SignerDelta {
	return FileIO{
		chunks: ch,
		data:   d,
	}
}

func (fs FileIO) Sign(file, signatureFile string) error {
	fs.sendFileDataToChunkerWorker(file)

	// get the created chunks, sign the data of the chunks and store the signature to the 'signatureFile'.
	// for the signature is used Sha1, because it is fast and it was designed to the check the consistency
	// of large files.
	f, err := os.Create(signatureFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	for ch := range fs.chunks {
		if _, err = f.Write([]byte(ch.String() + "\n")); err != nil {
			fmt.Println("err:", err)
		}
	}
	return err
}

func (fio FileIO) FindDelta(fileSignature, newFile string) Delta {
	chunks := decodeChunksOfSignatureFile(fileSignature)
	fio.sendFileDataToChunkerWorker(newFile)

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
	return Delta{NewChunks: newChunks, OldChunks: oldChunks}
}

// sendFileDataToChunkerWorker sends the data in the file
// through a channel to worker that will split data to chunks.
//
// The parameter 'file' is the file of which should be split to chunks.
func (fs FileIO) sendFileDataToChunkerWorker(file string) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("Error occurred during opening a file:", err)
		//return err
		return
	}

	for {
		data := make([]byte, 48)
		if _, err = f.Read(data); err != nil {
			close(fs.data)
			if err != io.EOF {
				fmt.Println("Error occurred during reading:", err)
				return
			}
			fmt.Println("All data from file have read.")
			return
		}
		for _, b := range data {
			fs.data <- b
		}
	}
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
		ch := CreateChunkFromString(scanner.Text())
		chunks[ch.Signature] = ch
	}
	return chunks
}
