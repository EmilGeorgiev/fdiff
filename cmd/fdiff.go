package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/EmilGeorgiev/fdiff"
	"github.com/EmilGeorgiev/fdiff/rollinghash"
)

// signature indicate that the program will
var signature = flag.Bool("signature", false, "create a file with signatures")
var delta = flag.Bool("delta", false, "create a file with signatures")
var oldFile = flag.String("old-file", "", "")
var signatureFile = flag.String("signature-file", "", "")
var newFile = flag.String("new-file", "", "")

func main() {
	flag.Parse()

	b := make(chan byte, 1000)
	ch := make(chan fdiff.Chunk, 1000)

	newHash := rollinghash.NewRabinFingerprint
	cfg := fdiff.ChunkConfig{
		WindowSize:            48,
		MinSizeChunk:          2048,
		MaxSizeChunk:          65536,
		FingerprintBreakPoint: 0,
	}

	chuncker := fdiff.NewChunker(newHash, cfg, b, ch)
	chuncker.Start()

	fmt.Println("Signature: ", *signature)
	fmt.Println("OldFile: ", *oldFile)
	fmt.Println("SignatureFile: ", *signatureFile)
	fs := fdiff.NewFileSignerDelta(b, ch)
	if *signature {
		if err := fs.Sign(*oldFile, *signatureFile); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Signature file is created")
		return
	} else if *delta {

		d, err := fs.FindDelta(*signatureFile, *newFile)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Old chunks that are updated or removed:")
		for _, c := range d.OldChunks {
			fmt.Printf("	- offset: %d, length: %d, hash: %s\n", c.Offset, c.Length, c.Signature)
		}

		fmt.Println("New chunks that replace the old ones:")
		for _, c := range d.NewChunks {
			fmt.Printf("	- offset: %d, length: %d, hash: %s\n", c.Offset, c.Length, c.Signature)
			fmt.Printf("	- %s\n", c.Data)
		}
	}
}
