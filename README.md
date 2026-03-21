# cbconv
CLI tool to convert comic book archive files into PDFs.

## How to build this application
- Clone this repo and run `go build ./...`

## How to use this application

Flags
- `-i`: Path to a comic book archive file or directory containing comic book archive files. Accepts: .cbz, .cbr
- `-o`: Specify an output file or directory for the converted PDF(s). Defaults to a sibling directory beside input path named `cbconv_output`
- `-help`: Display help with examples on how to use the application

## References
- [Wikipedia entry for the comic book archive format](https://en.wikipedia.org/wiki/Comic_book_archive)
