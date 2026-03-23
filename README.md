# cbconv
CLI tool to convert comic book archive files into PDFs. It was created so that I can convert my digital comic book collection to a format readable on any device that can read PDFs. You can read more about the comic book archive file format on its [Wikipedia entry](https://en.wikipedia.org/wiki/Comic_book_archive). This application can convert .cb7, .cbr, .cbt, and .cbz files to PDF.

## How to build this application
- Clone this repo and run `go build ./...`

## How to use this application

Flags
- `-i`: Path to a comic book archive file or directory containing comic book archive files. Accepts: .cbz, .cbr
- `-o`: Specify an output file or directory for the converted PDF(s). Defaults to a sibling directory beside input path named `cbconv_output`
- `-r`: Recursively convert subdirectories while in directory mode
- `-help`: Display help with examples on how to use the application
