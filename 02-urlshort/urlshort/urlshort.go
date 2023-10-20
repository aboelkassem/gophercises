package urlshort

import (
	"encoding/json"
	"fmt"
	"net/http"

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
		path := r.URL.Path
		fmt.Printf("Path: %v\n", path)
		fmt.Println(pathsToUrls)
		longURL, ok := pathsToUrls[path]
		if !ok {
			fmt.Println("Cannot file path, redirecting to default error page")
			// couldn't find the request's path in the map
			// pass it to next middleware pipe
			fallback.ServeHTTP(w, r)
			return
		}

		// otherwise, redirect to longURL
		// fmt.Fprintln(w, "redirecting to %v", longURL)
		http.Redirect(w, r, longURL, http.StatusMovedPermanently)
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

type LinkData struct {
	Path, URL string
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// parse yaml data into map of (maps of key and value)
	// parsedYamlData := []map[string]string{}
	var parsedYamlData = []LinkData{}
	err := yaml.Unmarshal([]byte(yml), &parsedYamlData)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", parsedYamlData)

	// build pathsToUrls from parsed yaml data to be passed into map handler
	pathsToUrls := map[string]string{}
	for _, yamlEntry := range parsedYamlData {
		pathsToUrls[yamlEntry.Path] = yamlEntry.URL
	}

	return MapHandler(pathsToUrls, fallback), nil
}

func JSONHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// parse json data into map of (maps of key and value)
	// parsedJsonData := []map[string]string{}
	var parsedJsonData = []LinkData{}
	err := json.Unmarshal([]byte(yml), &parsedJsonData)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", parsedJsonData)

	// build pathsToUrls from parsed yaml data to be passed into map handler
	pathsToUrls := map[string]string{}
	for _, yamlEntry := range parsedJsonData {
		pathsToUrls[yamlEntry.Path] = yamlEntry.URL
	}

	return MapHandler(pathsToUrls, fallback), nil
}
