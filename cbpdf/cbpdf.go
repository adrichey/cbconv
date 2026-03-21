package cbpdf

import (
	"archive/zip"
	"errors"
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
	"github.com/nwaples/rardecode"
)

var validExts = map[string]bool{
	".cb7": true,
	".cba": true,
	".cbr": true,
	".cbt": true,
	".cbz": true,
}

var imageExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".bmp":  true,
	".tiff": true,
	".webp": true,
}

// Convert reads a comic book archive at inputPath and writes a PDF to
// outputPath. If outputPath is empty, the PDF is written alongside the input
// file with the same base name and a .pdf extension.
//
// Returns an error if the file extension is not a recognized comic book format,
// if the file cannot be opened as a zip archive, or if any step of the
// conversion fails.
func Convert(inputPath, outputPath string) error {
	ext := strings.ToLower(filepath.Ext(inputPath))
	if !validExts[ext] {
		return errors.New("unsupported file extension: must be one of .cb7, .cba, .cbr, .cbt, .cbz")
	}

	destDir, err := os.MkdirTemp("", "cbpdf-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(destDir)

	if ext == ".cbr" {
		if err := extractRAR(inputPath, destDir); err != nil {
			return err
		}
	} else {
		if err := extractZip(inputPath, destDir); err != nil {
			return err
		}
	}

	// Recursively find the first directory containing image files
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
		return errors.New("no image files found in archive")
	}

	// Collect image files in alphabetical order
	entries, err := os.ReadDir(imageDir)
	if err != nil {
		return err
	}

	var imageFiles []string
	for _, e := range entries {
		if !e.IsDir() && imageExts[strings.ToLower(filepath.Ext(e.Name()))] {
			imageFiles = append(imageFiles, e.Name())
		}
	}
	sort.Strings(imageFiles)

	// Resolve output path
	if outputPath == "" {
		outputPath = strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + ".pdf"
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
			return err
		}
		cfg, imgType, err := image.DecodeConfig(f)
		f.Close()
		if err != nil {
			return err
		}

		w, h := float64(cfg.Width), float64(cfg.Height)
		pdf.AddPageFormat("P", fpdf.SizeType{Wd: w, Ht: h})

		switch strings.ToLower(imgType) {
		case "jpeg":
			imgType = "JPG"
		default:
			imgType = strings.ToUpper(imgType)
		}

		pdf.Image(imgPath, 0, 0, w, h, false, imgType, 0, "")
	}

	return pdf.OutputFileAndClose(outputPath)
}

func extractZip(inputPath, destDir string) error {
	r, err := zip.OpenReader(inputPath)
	if err != nil {
		return errors.New("file is not a valid zip archive: " + err.Error())
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
			return err
		}
		out, err := os.Create(outPath)
		if err != nil {
			rc.Close()
			return err
		}
		io.Copy(out, rc)
		out.Close()
		rc.Close()
	}
	return nil
}

func extractRAR(inputPath, destDir string) error {
	r, err := rardecode.OpenReader(inputPath, "")
	if err != nil {
		return errors.New("file is not a valid RAR archive: " + err.Error())
	}
	defer r.Close()

	for {
		header, err := r.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		outPath := filepath.Join(destDir, header.Name)
		if header.IsDir {
			os.MkdirAll(outPath, 0755)
			continue
		}
		os.MkdirAll(filepath.Dir(outPath), 0755)
		out, err := os.Create(outPath)
		if err != nil {
			return err
		}
		io.Copy(out, r)
		out.Close()
	}
	return nil
}
