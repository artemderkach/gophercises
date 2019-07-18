package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type Arc struct {
	Title   string
	Story   []string
	Options []Option
}

type Option struct {
	Text string
	Arc  string
}

type Rest struct {
	Story map[string]Arc
}

const tpl = `<h1>{{.Title}}</h1><br>
{{range $story := .Story}}
	{{$story}}<br><br>
{{end}}
<br></br>

{{range $option := .Options}}
	<a href="http://localhost:8080/{{$option.Arc}}">{{$option.Text}}<br></br></a>
{{end}}
`

func main() {
	rest := &Rest{}

	f, err := os.Open("./gopher.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	story, err := parseStory(f)
	if err != nil {
		fmt.Println(err)
		return
	}

	rest.Story = story

	http.HandleFunc("/", rest.home)

	fmt.Println("listening on localhost:8080")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func (rest *Rest) home(w http.ResponseWriter, r *http.Request) {
	path := strings.ReplaceAll(r.URL.Path, "/", "")

	t, err := template.New("webpage").Parse(tpl)
	if err != nil {
		err = errors.Wrap(err, "error parsing template")
		fmt.Println(err)
		w.Write([]byte(err.Error()))
		return
	}

	if path == "" || rest.Story[path].Title == "" {
		err = t.Execute(w, rest.Story["intro"])
		if err != nil {
			err = errors.Wrap(err, "error templating story")
			fmt.Println(err)
			w.Write([]byte(err.Error()))
			return
		}
		return
	}

	err = t.Execute(w, rest.Story[path])
	if err != nil {
		err = errors.Wrap(err, "error tempalting story")
		fmt.Println(err)
		w.Write([]byte(err.Error()))
		return
	}
}

func parseStory(r io.Reader) (map[string]Arc, error) {
	var story map[string]Arc

	body, err := ioutil.ReadAll(r)
	if err != nil {
		return story, errors.Wrap(err, "cannot read from reader")
	}
	err = json.Unmarshal(body, &story)
	if err != nil {
		return story, errors.Wrap(err, "erron unmarshalling json body")
	}
	return story, nil
}
