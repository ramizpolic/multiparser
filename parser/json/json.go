package json

import (
	"encoding/json"
	"github.com/ramizpolic/multiparser"
)

var Converter multiparser.Converter = &jsonParser{}

type jsonParser struct{}

func (p *jsonParser) Marshal(object interface{}) ([]byte, error) {
	return json.Marshal(object)
}

func (p *jsonParser) Unmarshal(from []byte, to interface{}) error {
	return json.Unmarshal(from, to)
}
