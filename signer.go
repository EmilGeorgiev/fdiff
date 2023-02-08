package fdiff

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
)

type Signer interface {
	Sign(file, signatureFile string) error
}

type Chunk struct {
	Offset    uint64
	Data      []byte
	Signature string
}

type fileSignature struct {
	chunks <-chan Chunk
	data   chan<- []byte
}

func NewFileSigner(d chan<- []byte, ch <-chan Chunk) Signer {
	return fileSignature{
		chunks: ch,
		data:   d,
	}
}

func (fs fileSignature) Sign(file, signatureFile string) error {
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
		h := fmt.Sprintf("%s", sha1.Sum(ch.Data))
		str := fmt.Sprintf("%x\n", h)
		if _, err = f.Write([]byte(str)); err != nil {
			fmt.Println("err:", err)
		}
	}
	fmt.Println("Signature file is created.")
	return err
}

// sendFileDataToChunkerWorker sends the data in the file
// through a channel to worker that will split data to chunks.
//
// The parameter 'file' is the file of which should be split to chunks.
func (fs fileSignature) sendFileDataToChunkerWorker(file string) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("Error occurred during opening a file:", err)
		//return err
		return
	}

	for {
		b1 := make([]byte, 48)
		if _, err = f.Read(b1); err != nil {
			close(fs.data)
			if err != io.EOF {
				fmt.Println("Error occurred during reading:", err)
				return
			}
			fmt.Println("All data from file have read.")
			return
		}
		fs.data <- b1
	}
}
