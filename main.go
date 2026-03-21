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

	"github.com/go-pdf/fpdf"
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

	// Determine output path
	outPath := outputFile
	if outPath == "" {
		ext := filepath.Ext(inputFile)
		outPath = strings.TrimSuffix(inputFile, ext) + ".pdf"
	}

	// Build PDF, one page per image sized to match the image dimensions
	pdf := fpdf.NewCustom(&fpdf.InitType{
		UnitStr: "pt",
		Size:    fpdf.SizeType{Wd: 2480, Ht: 3508}, // Default to A4 page size, overridden per page
	})
	pdf.SetAutoPageBreak(false, 0)

	for _, name := range imageFiles {
		imgPath := filepath.Join(imageDir, name)

		f, err := os.Open(imgPath)
		if err != nil {
			fmt.Println("Error opening image:", err)
			continue
		}
		cfg, imgType, err := image.DecodeConfig(f)
		f.Close()
		if err != nil {
			fmt.Println("Error decoding image:", err)
			continue
		}

		w, h := float64(cfg.Width), float64(cfg.Height)
		pdf.AddPageFormat("P", fpdf.SizeType{Wd: w, Ht: h})

		// fpdf uses the extension to detect type; normalize to a supported name
		switch strings.ToLower(imgType) {
		case "jpeg":
			imgType = "JPG"
		default:
			imgType = strings.ToUpper(imgType)
		}

		pdf.Image(imgPath, 0, 0, w, h, false, imgType, 0, "")
	}

	if err := pdf.OutputFileAndClose(outPath); err != nil {
		fmt.Println("Error writing PDF:", err)
		return
	}
	fmt.Println("PDF written to", outPath)
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
