package exercise_3

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"strings"
)

type OptionType struct {
	Text  string `json:"text"`
	ArcId string `json:"arc"`
}

type ArcType struct {
	Title   string       `json:"title"`
	Story   []string     `json:"story"`
	Options []OptionType `json:"options"`
}

type StoryBookType map[string]ArcType

var story_book StoryBookType
var arc_template *template.Template
var default_arc string

func LoadFile(file_name string, expected_ext string) ([]byte, error) {
	file_ext := strings.ToLower(strings.TrimPrefix(path.Ext(file_name), "."))

	if file_ext != expected_ext {
		return nil, errors.New(fmt.Sprint("story file needs to be a valid .", file_ext, " file"))
	}

	// Get file info
	file_info, err_file_info := os.Stat(file_name)

	if err_file_info != nil {
		return nil, err_file_info
	}

	// Open file
	f, err_f := os.Open(file_name)

	if err_f != nil {
		return nil, err_f
	}

	// Prepare buffer for content
	content := make([]byte, file_info.Size())

	// Read content
	_, content_err := f.Read(content)
	f.Close()

	if content_err != nil {
		return nil, content_err
	}
	return content, nil
}

func InitializeTemplate(template_file string) error {
	// Create template
	var arc_template_err error

	arc_template, arc_template_err = template.ParseFiles(template_file)

	if arc_template_err != nil {
		return arc_template_err
	}
	return nil
}

func GetStoryHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			var key string

			if r.URL.Path == "/" {
				key = default_arc
			} else {
				key = strings.TrimPrefix(r.URL.Path, "/")
			}
			template := story_book[key]
			arc_template.Execute(rw, template)
		}
	}
}

func PrepareKey(key string) string {
	var sb strings.Builder

	sb.WriteString("/")
	sb.WriteString(key)
	return sb.String()
}

func GetStoryBookHandler(story_file string, initial_arc string) (*http.ServeMux, error) {
	default_arc = initial_arc
	// Load and verify story file
	content, content_err := LoadFile(story_file, "json")

	if content_err != nil {
		return nil, content_err
	}

	// Parse storybook file
	err := json.Unmarshal(content, &story_book)

	if err != nil {
		return nil, err
	}

	// Create server
	mux := http.NewServeMux()

	for key := range story_book {
		mux.Handle(PrepareKey(key), GetStoryHandler())
	}
	return mux, nil
}
