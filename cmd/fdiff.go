package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/EmilGeorgiev/fdiff"
	"github.com/EmilGeorgiev/fdiff/rollinghash"
	"gopkg.in/yaml.v3"
)

// signature indicate that the program will
var signature = flag.Bool("signature", false, "create a signature file of a file.")
var delta = flag.Bool("delta", false, "find the difference between two files or two versions of the file.")
var oldFile = flag.String("old-file", "", "show for which file the signature will be created.")
var signatureFile = flag.String("signature-file", "", "show what will be the name of the signature file.")
var newFile = flag.String("new-file", "", "show the version of the file or the new file for which the command will find the delta.")
var showDelta = flag.Bool("show-data", false, "print the data in the new chunks")
var help = flag.Bool("help", false, "describe how to use the tool")

func main() {
	flag.Parse()

	if *help {
		printHelp()
		return
	}

	b := make(chan byte, 1000)
	ch := make(chan fdiff.Chunk, 1000)

	cfg := getConfig()
	newHash := rollinghash.NewRabinFingerprint

	chuncker := fdiff.NewChunker(newHash, cfg, b, ch)
	chuncker.Start()

	fs := fdiff.NewFileSignerDelta(b, ch)
	if *signature {
		fmt.Println("Creating a signature of the file: ", *signatureFile)
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
			if *showDelta {
				fmt.Printf("	- %s\n", c.Data)
			}
		}
	}
}

func getConfig() fdiff.ChunkConfig {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var config fdiff.ChunkConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("	fdiff -signature=true -old-file <name-of-file> -signature-file <name-of-sign-file>")
	fmt.Println("	fdiff -delta=true -signature-file <name-of-sign-file> -new-file <name-of_new-file>")

	fmt.Println("Flags:")
	fmt.Println("	- signature - create a signature file of a file.")
	fmt.Println("	- delta - find the difference between two files or two versions of the file.")
	fmt.Println("	- old-file - show for which file the signature will be created.")
	fmt.Println("	- signature-file - show what will be the name of the signature file.")
	fmt.Println("	- new-file - show the version of the file or the new file for which the command will find the delta.")
	fmt.Println("	- show-data - print the data in the new chunks.")
	fmt.Println("	- help - describe how to use the tool.")
}
