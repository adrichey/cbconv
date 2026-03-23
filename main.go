package main

import (
	"flag"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/adrichey/cbconv/cbpdf"
	"github.com/adrichey/cbconv/project"
)

const OUTPUT_DIR = "cbconv_output"
const RECURSIVE_FLAG_HELP = "Recursively convert subdirectories while in directory mode"

// Flags
var help bool
var input string
var output string
var recursive bool

var validExts map[string]bool
var inputHelpText string
var outputHelpText string

func init() {
	validExts = map[string]bool{
		".cb7": true,
		".cbr": true,
		".cbt": true,
		".cbz": true,
	}

	keys := maps.Keys(validExts)

	inputHelpText = "Path to a comic book archive file or directory containing comic book archive files. Accepts: " + strings.Join(slices.Collect(keys), ", ")
	outputHelpText = fmt.Sprintf("Specify an output file or directory for the converted PDF(s). Defaults to a sibling directory beside input path named %s", OUTPUT_DIR)

	flag.BoolVar(&help, "help", false, "Help")
	flag.StringVar(&input, "i", "", inputHelpText)
	flag.StringVar(&output, "o", "", outputHelpText)
	flag.BoolVar(&recursive, "r", false, RECURSIVE_FLAG_HELP)

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

	// Use -o as output directory if provided, otherwise default to sibling directory with name OUTPUT_DIR
	outDir := output
	if outDir == "" {
		outDir = filepath.Join(filepath.Dir(filepath.Clean(input)), OUTPUT_DIR)
	}

	// input is a directory — recursively collect all comic archive files
	var files []string
	inputClean := filepath.Clean(input)
	filepath.Walk(inputClean, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		if validExts[strings.ToLower(filepath.Ext(path))] {
			files = append(files, path)
		}
		return nil
	})

	if len(files) == 0 {
		fmt.Println("No comic archive files found in directory.")
		return
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
		// Preserve subdirectory structure relative to input root
		rel, err := filepath.Rel(inputClean, f)
		if err != nil {
			rel = filepath.Base(f)
		}
		outPath := filepath.Join(outDir, strings.TrimSuffix(rel, filepath.Ext(rel))+".pdf")
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			fmt.Printf("Error creating output directory for %s: %v\n", filepath.Base(f), err)
			continue
		}
		fmt.Printf("Converting %s...\n", rel)
		if err := cbpdf.Convert(f, outPath); err != nil {
			fmt.Printf("Error converting %s: %v\n", rel, err)
		}
	}
}

func displayHelp() {
	fmt.Println("How to use this script:")
	fmt.Printf("-i: %s\n", inputHelpText)
	fmt.Printf("-o: %s\n", outputHelpText)
	fmt.Printf("-r: %s\n", RECURSIVE_FLAG_HELP)

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
	formattedString = fmt.Sprintf("This will save the converted files in .|one|two|comics to .|one|two|%s", OUTPUT_DIR)
	fmt.Println(strings.ReplaceAll(formattedString, "|", project.PATH_SEPARATOR))
	fmt.Println("If you include the -r flag with this command, it will recursively convert all subdirectory comic archives; otherwise it just converts top-level directory passed into the application")
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
