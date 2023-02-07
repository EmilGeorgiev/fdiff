package fdiff

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Signer interface {
	Sign(file, signatureFile string) error
}

type Chunk struct {
	Offset        uint64
	Length        uint64
	HashSignature string
}

func (ch Chunk) Bytes() []byte {
	str := fmt.Sprintf("%d-%d-%s\n", ch.Offset, ch.Length, ch.HashSignature)
	return []byte(str)
}

type fileSignature struct {
	chunks <-chan Chunk
	bytes  chan<- []byte
}

func NewFileSignature(ch <-chan Chunk, b chan<- []byte) Signer {
	return fileSignature{
		chunks: ch,
		bytes:  b,
	}
}

func (fs fileSignature) Sign(file, signatureFile string) error {
	fs.storeSignaturesOfTheChunksToFile(signatureFile)

	f, err := os.Open(file)
	if err != nil {
		fmt.Println("Error occurred during opening a file:", err)
		return err
	}

	for {
		b1 := make([]byte, 48)
		if _, err = f.Read(b1); err != nil {
			close(fs.bytes)
			if err != io.EOF {
				fmt.Println("Error occurred during reading:", err)
				return err
			}
			fmt.Println("All bytes from file have read.")
		}
		fs.bytes <- b1
	}

}

func (fs fileSignature) storeSignaturesOfTheChunksToFile(file string) {
	f, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	for ch := range fs.chunks {
		if _, err = f.Write(ch.Bytes()); err != nil {
			fmt.Println("err:", err)
		}
	}
}
