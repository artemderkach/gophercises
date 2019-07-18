package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Arc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Arc    `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

// func main() {
// 	story, err:= parseStory()
// 	if err != nil {
// 		log.Fatal(err)
// 		return
// 	}
// }

// parse Arc format
func parseStory() (map[string]Arc, error) {
	f, err := os.Open("./gopher.json")
	if err != nil {
		return make(map[string]Arc, 0), err
	}

	body, err := ioutil.ReadAll(f)
	if err != nil {
		return make(map[string]Arc, 0), err
	}

	var story map[string]Arc
	err = json.Unmarshal(body, &story)
	if err != nil {
		return make(map[string]Arc, 0), err
	}

	return story, nil
}
