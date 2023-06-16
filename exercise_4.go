package main

import (
	"exercise_4"
	"flag"

	"github.com/spf13/afero"
)

func main() {
	source_url_ptr := flag.String("source-url", "index.html", "Specifies the source URL.")
	target_file_ptr := flag.String("target-file", "output.txt", "Specifies the target file.")
	flag.Parse()

	fs := afero.NewOsFs()

	exercise_4.ExtractLinks(*source_url_ptr, *target_file_ptr, fs)
}
