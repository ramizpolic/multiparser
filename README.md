# multiparser
Parse raw data into concrete types unburdened by serialization standards. 

### Example usage
Current multiparser example uses JSON and YAML parsers.
```go
package main

import (
    "github.com/ramizpolic/multiparser"
    "github.com/ramizpolic/multiparser/parser"
)

type object struct {
    Data string `json:"data" yaml:"data"`
}

func main() {
    parser, _ := multiparser.New(parser.JSON, parser.YAML)

    // Parse JSON
    var jsonObj object
    _ = parser.Parse([]byte(`{"data": "data"}`), &jsonObj)
	
    // Parse YAML
    var yamlObj object
    _ = parser.Parse([]byte(`data: data`), &yamlObj)
}
```

### Supported parsers
- JSON - `encoding.json`
- YAML - `gopkg.in/yaml.v3`

### Bring your own parser
All you have to do is implement `multiparser.Parser` interface, e.g.
```golang
type parser struct {}

// Parse converts raw to object
func (p *parser) Parse(from []byte, to interface{}) error {
    panic("implement me")
}
```
