package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrichey/cbconv/cbpdf"
	"github.com/adrichey/cbconv/project"
)

const OUTPUT_DIR = "cbconv_output"
const inputHelpText = "Path to a comic book archive file or directory containing comic book archive files. Accepts: .cb7, .cbr"

var outputHelpText = fmt.Sprintf("Specify an output file or directory for the converted PDF(s). Defaults to a sibling directory beside input path named %s", OUTPUT_DIR)

var help bool
var input string
var output string

func init() {
	flag.BoolVar(&help, "help", false, "Help")
	flag.StringVar(&input, "i", "", inputHelpText)
	flag.StringVar(&output, "o", "", outputHelpText)

	flag.Parse()
}

func main() {
	if input == "" && output == "" {
		displayHelp()
		return
	}

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

	// Use -o as output directory if provided, otherwise default to sibling directory with name OUTPUT_DIR
	outDir := output
	if outDir == "" {
		outDir = filepath.Join(filepath.Dir(filepath.Clean(input)), OUTPUT_DIR)
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
	fmt.Printf("-i: %s\n", inputHelpText)
	fmt.Printf("-o: %s\n", outputHelpText)

	ex, err := os.Executable()
	if err != nil {
		return
	}

	ex = filepath.Base(ex)

	fmt.Println()

	fmt.Println("Example - Single File:")
	formattedString := fmt.Sprintf("%s -i .|path|to|file|example.cbz", ex)
	fmt.Println(strings.ReplaceAll(formattedString, "|", project.PATH_SEPARATOR))
	fmt.Println(strings.ReplaceAll("This will save the converted file to .|path|to|file|example.pdf", "|", project.PATH_SEPARATOR))
	fmt.Println()

	fmt.Println("Example - Directory:")
	formattedString = fmt.Sprintf("%s -i .|one|two|comics", ex)
	fmt.Println(strings.ReplaceAll(formattedString, "|", project.PATH_SEPARATOR))
	formattedString = fmt.Sprintf("This will save the converted files to .|one|two|%s", OUTPUT_DIR)
	fmt.Println(strings.ReplaceAll(formattedString, "|", project.PATH_SEPARATOR))
	fmt.Println()

	fmt.Println("Example - Single File with Specified Output File:")
	formattedString = fmt.Sprintf("%s -i .|path|to|file|example.cbz -o .|converted|comic.pdf", ex)
	fmt.Println(strings.ReplaceAll(formattedString, "|", project.PATH_SEPARATOR))
	fmt.Println(strings.ReplaceAll("This will save the converted file to .|converted|comic.pdf", "|", project.PATH_SEPARATOR))
	fmt.Println()

	fmt.Println("Example - Directory with Specified Output Directory:")
	formattedString = fmt.Sprintf("%s -i .|one|two|comics -o .|one|converted_comics", ex)
	fmt.Println(strings.ReplaceAll(formattedString, "|", project.PATH_SEPARATOR))
	fmt.Println(strings.ReplaceAll("This will save the converted files to .|one|converted_comics", "|", project.PATH_SEPARATOR))
	fmt.Println()
}
