package multiparser_test

import (
	"github.com/ramizpolic/multiparser"
	"github.com/ramizpolic/multiparser/parser"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var parsers = []multiparser.Parser{
	parser.JSON,
	parser.YAML,
}

func TestParse(t *testing.T) {
	parser, _ := multiparser.New(parsers...)
	for _, tt := range []struct {
		name     string
		inputRaw string
		expected *map[string]string
	}{
		{
			name:     "json",
			inputRaw: `{"data": "data"}`,
			expected: &map[string]string{"data": "data"},
		},
		{
			name:     "yaml",
			inputRaw: `data: data`,
			expected: &map[string]string{"data": "data"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var got map[string]string
			_ = parser.Parse([]byte(tt.inputRaw), &got)
			assert.Equal(t, *tt.expected, got)
		})
	}
}

func TestAll(t *testing.T) {
	for _, tt := range []struct {
		name      string
		parsers   []multiparser.Parser
		parserErr error
		input     string
		inputErr  error
		outputPtr interface{}
		outputErr error
	}{
		{
			name:      "nil-prs",
			parserErr: multiparser.ErrEmptyParsers,
		},
		{
			name:      "all-prs-json",
			parsers:   parsers,
			input:     `{"data": "data"}`,
			outputPtr: &map[string]string{"data": "data"},
		},
		{
			name:      "all-prs-yaml",
			parsers:   parsers,
			input:     `data: data`,
			outputPtr: &map[string]string{"data": "data"},
		},
		{
			name:      "all-prs-duplicate",
			parsers:   []multiparser.Parser{parser.JSON, parser.JSON, parser.JSON},
			input:     `{"data": "data"}`,
			outputPtr: &map[string]string{"data": "data"},
		},
		{
			name:      "all-prs-input-empty",
			parsers:   parsers,
			outputPtr: &struct{}{},
		},
		{
			name:      "all-prs-output-nil",
			parsers:   parsers,
			outputErr: multiparser.ErrInvalidObject,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// Create
			parser, err := multiparser.New(tt.parsers...)
			assert.Equalf(t, tt.parserErr, err, "parser new error expected %v, got %v for %s", tt.parserErr, err, tt.name)
			if tt.parserErr != nil {
				return
			}

			// Parse
			var output interface{}
			if tt.outputPtr != nil {
				output = reflect.New(reflect.TypeOf(tt.outputPtr).Elem()).Interface()
			}

			err = parser.Parse([]byte(tt.input), output)
			assert.Equalf(t, tt.outputErr, err, "parse error expected %v, got %v for %s", tt.outputErr, err, tt.name)
			if tt.outputErr != nil {
				return
			}

			assert.Equalf(t, tt.outputPtr, output, "parse result expected %v, got %v for %s", tt.outputPtr, output, tt.name)
		})
	}

}
