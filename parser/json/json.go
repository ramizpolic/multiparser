package json

import (
	"encoding/json"
	"github.com/ramizpolic/multiparser"
)

var Parser multiparser.Parser = &jsonParser{}

type jsonParser struct{}

func (p *jsonParser) Parse(from []byte, to interface{}) error {
	return json.Unmarshal(from, to)
}
