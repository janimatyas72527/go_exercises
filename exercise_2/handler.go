package exercise_2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/yaml.v3"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, found := pathsToUrls[r.URL.Path]
		if found {
			http.Redirect(w, r, url, http.StatusMovedPermanently)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// TODO: Write doc comment
type ParsedEntry struct {
	Path string `yaml:"path" json:"path"`
	Url  string `yaml:"url" json:"url"`
}

// TODO: Write doc comment
func (pe ParsedEntry) String() string {
	return fmt.Sprintf("[Path: %s, Url: %s]", pe.Path, pe.Url)
}

// TODO: Write doc comment
type ParsedJsonType struct {
	Mappings []ParsedEntry `json:"mappings"`
}

// TODO: Write doc comment
func (pjt ParsedJsonType) String() string {
	var items []string
	for _, item := range pjt.Mappings {
		items = append(items, item.String())
	}
	return fmt.Sprintf("[Mappings: [%s]]", strings.Join(items, ","))
}

// TODO: Write doc comment
func processParsedEntries(parsedEntries []ParsedEntry, fallback http.Handler, err error) (http.HandlerFunc, error) {
	if err != nil {
		return nil, err
	} else {
		convertedMap := make(map[string]string)
		for _, parsedEntry := range parsedEntries {
			convertedMap[parsedEntry.Path] = parsedEntry.Url
		}
		return MapHandler(convertedMap, fallback), nil
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var parsedEntries []ParsedEntry

	err := yaml.Unmarshal(yml, &parsedEntries)

	var items []string

	for _, item := range parsedEntries {
		items = append(items, item.String())
	}
	fmt.Printf("[%s]\n", strings.Join(items, ","))

	return processParsedEntries(parsedEntries, fallback, err)
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
//		[
//			{
//				"path": "/some-path"
//	    		"url": "https://www.some-url.com/demo"
//			}
//		]
//
// The only errors that can be returned all related to having
// invalid JSON data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var parsedJson ParsedJsonType

	err := json.Unmarshal(jsn, &parsedJson)
	fmt.Println(parsedJson.String())
	return processParsedEntries(parsedJson.Mappings, fallback, err)
}
