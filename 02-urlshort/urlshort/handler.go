package urlshort

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"
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
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
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
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yamlFile string, fallback http.Handler) (http.HandlerFunc, error) {
	pathURLs, err := parseYAML(yamlFile)
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap(pathURLs)
	return MapHandler(pathsToUrls, fallback), nil
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
func JSONHandler(jsonFile string, fallback http.Handler) (http.HandlerFunc, error) {
	pathURLs, err := parseJSON(jsonFile)
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap(pathURLs)
	return MapHandler(pathsToUrls, fallback), nil
}

type pathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func parseYAML(yamlFile string) ([]pathURL, error) {
	yamlIO, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return nil, err
	}

	var pathURLs []pathURL
	err = yaml.Unmarshal(yamlIO, &pathURLs)
	if err != nil {
		return nil, err
	}

	return pathURLs, nil
}

func parseJSON(jsonFile string) ([]pathURL, error) {
	jsonIO, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return nil, err
	}

	var pathURLs []pathURL
	err = json.Unmarshal(jsonIO, &pathURLs)
	if err != nil {
		return nil, err
	}

	return pathURLs, nil
}

func buildMap(pathURLs []pathURL) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, pu := range pathURLs {
		pathsToUrls[pu.Path] = pu.URL
	}
	return pathsToUrls
}
