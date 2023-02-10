package main

import (
	"flag"
	"fmt"

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
	chuncker := fdiff.NewChunker(newHash, 48, 48, 1024, 0, b, ch)
	chuncker.Start()

	fmt.Println("Signature: ", *signature)
	fmt.Println("OldFile: ", *oldFile)
	fmt.Println("SignatureFile: ", *signatureFile)
	fs := fdiff.NewFileIO(b, ch)
	if *signature {
		_ = fs.Sign(*oldFile, *signatureFile)
		fmt.Println("Signature file is created")
		return
	} else if *delta {

		d := fs.FindDelta(*signatureFile, *newFile)
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
