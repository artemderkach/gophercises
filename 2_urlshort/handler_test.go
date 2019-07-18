package urlshort

import (
	"testing"

	"gopkg.in/yaml.v2"
)

type s struct {
	S string `yaml:"s"`
}

func TestYAMLHandler(t *testing.T) {
	yamlData := `s: aaaa`
	var parsedYaml s
	err := yaml.Unmarshal([]byte(yamlData), &parsedYaml)
	if err != nil {
		t.Fatal(err)
	}
	t.Error("error:", parsedYaml)
}
