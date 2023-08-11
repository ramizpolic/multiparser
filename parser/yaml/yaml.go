package yaml

import (
	"github.com/ramizpolic/multiparser"
	"gopkg.in/yaml.v3"
)

var Converter multiparser.Converter = &yamlParser{}

type yamlParser struct{}

func (p *yamlParser) Marshal(object interface{}) ([]byte, error) {
	return yaml.Marshal(object)
}

func (p *yamlParser) Unmarshal(from []byte, to interface{}) error {
	return yaml.Unmarshal(from, to)
}
