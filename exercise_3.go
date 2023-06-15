package main

import (
	"exercise_3"
	"flag"
	"fmt"
	"net/http"
)

func main() {
	story_file_ptr := flag.String("story-file", "story.json", "Specifies the story file.")
	initial_arc_ptr := flag.String("initial-arc", "intro", "Specifies the initial arc.")
	arc_template_ptr := flag.String("arc-template", "arc_template.html", "Specifies the arc HTML template.")
	flag.Parse()

	content_err := exercise_3.InitializeTemplate(*arc_template_ptr)

	if content_err != nil {
		panic(content_err)
	}

	mux, mux_err := exercise_3.GetStoryBookHandler(*story_file_ptr, *initial_arc_ptr)

	if mux_err != nil {
		panic(mux_err)
	}
	fmt.Print("Starting the server on :8080\n")
	http.ListenAndServe(":8080", mux)
}
