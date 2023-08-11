package parser

import (
	"github.com/ramizpolic/multiparser/parser/json"
	"github.com/ramizpolic/multiparser/parser/yaml"
)

var (
	JSON = json.Parser
	YAML = yaml.Parser
)
