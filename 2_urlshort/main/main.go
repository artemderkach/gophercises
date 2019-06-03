package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	urlshort "github.com/mind-rot/gophercises/2_urlshort"
)

var defaultMap = map[string]string{
	"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
	"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
}

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	mapHandler := urlshort.MapHandler(defaultMap, mux)

	// Build the YAMLHandler using the mapHandler as the fallback
	yamlHandler, err := urlshort.YAMLHandler(parseYaml(), mapHandler)
	if err != nil {
		panic(err)
	}
	jsonHandler, err := urlshort.JSONHandler(parseJson(), yamlHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func parseYaml() []byte {
	defaultYamlMap := []byte(`
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`)
	yamlFilePath := flag.String("yamlFile", "", "yaml file to read urls and paths from")
	flag.Parse()
	if *yamlFilePath == "" {
		return defaultYamlMap
	}
	yamlFile, err := os.Open(*yamlFilePath)
	if err != nil {
		panic(err)
	}
	yamlFileContent, err := ioutil.ReadAll(yamlFile)
	if err != nil {
		panic(err)
	}
	return yamlFileContent
}

func parseJson() []byte {
	defaultJsonMap := []byte(`
[
    {
		"path": "/urlshort-gopher",
		"url": "https://gophercises.com/exercises/urlshort"
	}
]
`)
	return defaultJsonMap
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
