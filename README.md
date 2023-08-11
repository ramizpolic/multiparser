# multiparser
Serialize and deserialize different data types easily. 

### Usage
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
    _ = parser.Unmarshal([]byte(`{"data": "data"}`), &jsonObj)


    // Parse YAML
    var yamlObj object
    _ = parser.Unmarshal([]byte(`data: data`), &yamlObj)
}
```

### Supported parsers
- JSON - `encoding.json`
- YAML - `gopkg.in/yaml.v3`

### Bring your own parser
All you have to do is implement `multiparser.Parser` interface, e.g.
```golang
type parser struct {}

// Marshal converts object to raw
func (p *parser) Marshal(object interface{}) ([]byte, error) {
    panic("implement me")
}

// Unmarshal converts raw to object
func (p *parser) Unmarshal(from []byte, to interface{}) error {
    panic("implement me")
}
```
