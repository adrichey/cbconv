package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrichey/cbconv/cbpdf"
)

var help bool
var input string
var output string

func init() {
	flag.BoolVar(&help, "help", false, "Help")
	flag.StringVar(&input, "i", "", "Path to a comic book archive file (.cb7, .cba, .cbr, .cbt, .cbz)")
	flag.StringVar(&output, "o", "", "Specify an output file for the converted PDF")

	flag.Parse()
}

func main() {
	if help {
		displayHelp()
		return
	}

	info, err := os.Stat(input)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if !info.IsDir() {
		if err := cbpdf.Convert(input, output); err != nil {
			fmt.Println("Error:", err)
		}
		return
	}

	// input is a directory — collect all comic archive files
	validExts := map[string]bool{
		".cb7": true, ".cba": true, ".cbr": true, ".cbt": true, ".cbz": true,
	}

	var files []string
	entries, err := os.ReadDir(input)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}
	for _, e := range entries {
		if !e.IsDir() && validExts[strings.ToLower(filepath.Ext(e.Name()))] {
			files = append(files, filepath.Join(input, e.Name()))
		}
	}

	if len(files) == 0 {
		fmt.Println("No comic archive files found in directory.")
		return
	}

	// Use -o as output directory if provided, otherwise default to sibling cbconv_output
	outDir := output
	if outDir == "" {
		outDir = filepath.Join(filepath.Dir(filepath.Clean(input)), "cbconv_output")
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Println("Error creating output directory:", err)
		return
	}

	fullOutputPath, err := filepath.Abs(outDir)
	if err != nil {
		fullOutputPath = outDir
	}

	fmt.Printf("Saving converted PDFs to: %s\n", fullOutputPath)

	for _, f := range files {
		base := strings.TrimSuffix(filepath.Base(f), filepath.Ext(f))
		outPath := filepath.Join(outDir, base+".pdf")
		fmt.Printf("Converting %s...\n", filepath.Base(f))
		if err := cbpdf.Convert(f, outPath); err != nil {
			fmt.Printf("Error converting %s: %v\n", filepath.Base(f), err)
		}
	}
}

func displayHelp() {
	fmt.Println("How to use this script:")
	fmt.Println("-i: Path to a comic book archive file (.cb7, .cba, .cbr, .cbt, .cbz) or a directory containing comic archive files")
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
