package yaml

import (
	"github.com/ramizpolic/multiparser"
	"gopkg.in/yaml.v3"
)

var Parser multiparser.Parser = &yamlParser{}

type yamlParser struct{}

func (p *yamlParser) Parse(from []byte, to interface{}) error {
	return yaml.Unmarshal(from, to)
}
