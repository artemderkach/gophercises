package urlshort

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := pathsToUrls[r.URL.Path]; ok {
			fmt.Println(pathsToUrls[r.URL.Path])
			http.Redirect(w, r, pathsToUrls[r.URL.Path], http.StatusSeeOther)
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
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYaml(yml)
	if err != nil {
		return nil, err
	}

	fmt.Println(parsedYaml)
	return MapHandler(parsedYaml, fallback), nil
}

type Map struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func parseYaml(data []byte) (map[string]string, error) {
	var linkMap []Map
	err := yaml.Unmarshal(data, &linkMap)
	if err != nil {
		return nil, err
	}

	var m = make(map[string]string, len(linkMap))
	for _, elem := range linkMap {
		m[elem.Path] = elem.URL
	}

	return m, nil
}
