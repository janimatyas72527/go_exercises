package main

import (
	"exercise_2"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func main() {
	file_ptr := flag.String("file", "", "Specifies the source file.")
	flag.Parse()

	_, err_file_info := os.Stat(*file_ptr)

	if err_file_info != nil {
		panic(err_file_info)
	}

	mux := defaultMux()

	file_ext := strings.ToLower(path.Ext(*file_ptr))

	f, err_f := os.Open(*file_ptr)

	if err_f != nil {
		log.Fatal(err_f)
	}

	var content []byte

	_, err_content := f.Read(content)

	if err_content != nil {
		log.Fatal(err_content)
	}

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v3",
	}
	mapHandler := exercise_2.MapHandler(pathsToUrls, mux)
	var handler http.Handler
	var err error

	switch file_ext {
	case "yaml":
		// Build the YAMLHandler using the mapHandler as the
		// fallback
		handler, err = exercise_2.YAMLHandler(content, mapHandler)
		if err != nil {
			panic(err)
		}
	case "json:":
		// Build the JSONHandler using the mapHandler as the
		// fallback
		handler, err = exercise_2.JSONHandler(content, mapHandler)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Starting the", strings.ToUpper(file_ext), "-configured server on :8080")
	http.ListenAndServe(":8080", handler)
}

// TODO: Write doc comment
func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

// TODO: Write doc comment
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
