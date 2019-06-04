package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/boltdb/bolt"
	urlshort "github.com/mind-rot/gophercises/2_urlshort"
)

var defaultMap = map[string]string{
	"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
	"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
}

func main() {

	mux := defaultMux()

	initDB()
	// Build the MapHandler using the mux as the fallback
	db, err := bolt.Open("db.bolt", 0600, nil)
	if err != nil {
		panic(err)
	}

	dbHandler := urlshort.DBHandler(db, mux)
	// mapHandler := urlshort.MapHandler(defaultMap, mux)

	// Build the YAMLHandler using the mapHandler as the fallback
	yamlHandler, err := urlshort.YAMLHandler(parseYaml(), dbHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func initDB() {
	db, err := bolt.Open("db.bolt", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("urls"))
		if err != nil {
			return err
		}
		for k, v := range defaultMap {
			b.Put([]byte(k), []byte(v))
			fmt.Println("Writing key value pair:", k, v)
		}
		return nil
	})
	defer db.Close()
}

func parseYaml() []byte {
	var defaultYamlMap = []byte(`
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

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
