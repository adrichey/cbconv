package main

import (
	"flag"
	"fmt"

	"github.com/adrichey/cbconv/cbpdf"
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

	if err := cbpdf.Convert(inputFile, outputFile); err != nil {
		fmt.Println("Error:", err)
	}
}

func displayHelp() {
	fmt.Println("How to use this script:")
	fmt.Println("-i: Path to a comic book archive file (e.g. .cb7, .cba, .cbr, .cbt, .cbz)")
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
