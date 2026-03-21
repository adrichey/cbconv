package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

	// Unzip CBZ file
	destDir, err := os.MkdirTemp("", "cbconv-*")
	if err != nil {
		fmt.Println("Error creating temp dir:", err)
		return
	}
	defer os.RemoveAll(destDir)

	r, err := zip.OpenReader(inputFile)
	if err != nil {
		fmt.Println("Error opening CBZ file:", err)
		return
	}
	defer r.Close()

	for _, f := range r.File {
		outPath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(outPath, f.Mode())
			continue
		}
		os.MkdirAll(filepath.Dir(outPath), 0755)
		rc, err := f.Open()
		if err != nil {
			fmt.Println("Error opening zip entry:", err)
			return
		}
		out, err := os.Create(outPath)
		if err != nil {
			rc.Close()
			fmt.Println("Error creating file:", err)
			return
		}
		io.Copy(out, rc)
		out.Close()
		rc.Close()
	}

	// Recursively find first directory containing image files
	imageExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".bmp": true, ".tiff": true, ".webp": true,
	}

	imageDir := ""
	filepath.Walk(destDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || imageDir != "" {
			return err
		}
		if !info.IsDir() && imageExts[strings.ToLower(filepath.Ext(path))] {
			imageDir = filepath.Dir(path)
		}
		return nil
	})

	if imageDir == "" {
		fmt.Println("No image files found in archive.")
		return
	}

	// Collect image files in alphabetical order
	entries, err := os.ReadDir(imageDir)
	if err != nil {
		fmt.Println("Error reading image directory:", err)
		return
	}

	var imageFiles []string
	for _, e := range entries {
		if !e.IsDir() && imageExts[strings.ToLower(filepath.Ext(e.Name()))] {
			imageFiles = append(imageFiles, e.Name())
		}
	}
	sort.Strings(imageFiles)

	// Print image name, height, and width
	for _, name := range imageFiles {
		f, err := os.Open(filepath.Join(imageDir, name))
		if err != nil {
			fmt.Println("Error opening image:", err)
			continue
		}
		cfg, _, err := image.DecodeConfig(f)
		f.Close()
		if err != nil {
			fmt.Println("Error decoding image:", err)
			continue
		}
		fmt.Printf("%s: width=%d height=%d\n", name, cfg.Width, cfg.Height)
	}
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
