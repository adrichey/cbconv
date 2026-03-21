package main

import (
	"flag"
	"fmt"
)

var help bool
var inputFile string
var outputFile string

func init() {
	flag.BoolVar(&help, "help", false, "Help")
	flag.StringVar(&inputFile, "i", "", "Path to a comic book archive file (e.g. *.cbz, *.cbr)")
	flag.StringVar(&outputFile, "o", "", "Specify an output file for the converted PDF")

	flag.Parse()
}

func main() {
	if help {
		displayHelp()
		return
	}

	// TODO:
	// Unzip CBZ file using the archive/zip go package
	// Recursively traverse the unzipped directory and find the first instance that contains image files (can be jpg, png, gif, etc.)
	// Once we find that directory, create a slice to hold the image file names in alphabetical order
	// Traverse that slice and print the image name, the image height, and the image width to the console using the fmt package
}

func displayHelp() {
	fmt.Println("How to use this script:")
	fmt.Println("-i: Path to a comic book archive file (e.g. *.cbz, *.cbr)")
	fmt.Println("-o: Specify an output file for the converted PDF")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("./cbconv -i ./path/to/file/example.cbz")
	fmt.Println("This will save the converted file to ./path/to/file/example.pdf")
	fmt.Println()
	fmt.Println("Example with optional args:")
	fmt.Println("./cbconv -i ./path/to/file/example.cbz -o ./converted/comic.pdf")
	fmt.Println("This will save the converted file to ./converted/comic.pdf")
	fmt.Println()
}
