package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseStory(t *testing.T) {
	expectedStory := map[string]Arc{
		"intro": Arc{
			Title: "title 1",
			Story: []string{
				"story a",
				"story b",
				"story c",
			},
			Options: []Option{
				{Text: "text1", Arc: "arc1"},
				{Text: "text2", Arc: "arc2"},
			},
		},
		"arc2": Arc{
			Title: "title 2",
			Story: []string{
				"story aa",
				"story bb",
				"story cc",
				"story dd",
				"story ee",
			},
			Options: []Option{
				{Text: "text11", Arc: "arc11"},
				{Text: "text22", Arc: "arc22"},
			},
		},
	}

	f, err := os.Open("./gopher_test.json")
	require.Nil(t, err)

	story, err := parseStory(f)
	assert.Equal(t, expectedStory, story)
}

func TestHome(t *testing.T) {
	rest := &Rest{}
	rest.Story = map[string]Arc{
		"arc": Arc{
			Title: "title",
			Story: []string{"story1", "story2"},
			Options: []Option{
				{Text: "text", Arc: "arc"},
			},
		},
	}

	req, err := http.NewRequest("GET", "/", nil)
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(rest.home)

	handler.ServeHTTP(rr, req)

	body, err := ioutil.ReadAll(rr.Body)
	require.Nil(t, err)

	expected := "<h1></h1><br>\n\n<br></br>\n\n\n"
	assert.Equal(t, expected, string(body))

	req, err = http.NewRequest("GET", "/arc", nil)
	require.Nil(t, err)

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(rest.home)

	handler.ServeHTTP(rr, req)

	body, err = ioutil.ReadAll(rr.Body)
	require.Nil(t, err)

	expected = "<h1>title</h1><br>\n\n\tstory1<br><br>\n\n\tstory2<br><br>\n\n<br></br>\n\n\n\t<a href=\"http://localhost:8080/arc\">text<br></br></a>\n\n"
	assert.Equal(t, expected, string(body))
}
