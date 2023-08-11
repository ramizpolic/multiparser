package multiparser_test

import (
	"github.com/ramizpolic/multiparser"
	"github.com/ramizpolic/multiparser/parser"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type objType struct {
	Data string `json:"data" yaml:"data"`
}

var (
	parsers = []multiparser.Parser{
		parser.JSON,
		parser.YAML,
	}
	defaultObj    = objType{Data: "data"}
	defaultObjPtr = &defaultObj
)

func TestUnmarshal(t *testing.T) {
	parser, _ := multiparser.New(parsers...)
	for _, tt := range []struct {
		name     string
		inputRaw string
		expected objType
	}{
		{
			name:     "json",
			inputRaw: `{"data": "data"}`,
			expected: objType{Data: "data"},
		},
		{
			name:     "yaml",
			inputRaw: `data: data`,
			expected: objType{Data: "data"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var got objType
			_ = parser.Unmarshal([]byte(tt.inputRaw), &got)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestAll(t *testing.T) {
	for _, tt := range []struct {
		name      string
		parsers   []multiparser.Parser
		parserErr error
		input     interface{}
		inputErr  error
		outputPtr interface{}
		outputErr error
	}{
		{
			name:      "nil-prs",
			parserErr: multiparser.ErrEmptyParsers,
		},
		{
			name:      "all-prs",
			parsers:   parsers,
			input:     defaultObj,
			outputPtr: defaultObjPtr,
		},
		{
			name:      "all-prs-duplicate",
			parsers:   []multiparser.Parser{parser.JSON, parser.JSON, parser.JSON},
			input:     defaultObj,
			outputPtr: defaultObjPtr,
		},
		{
			name:      "all-prs-input-nil",
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

			// Marshal
			input, err := parser.Marshal(tt.input)
			assert.Equalf(t, tt.inputErr, err, "marshal error expected %v, got %v for %s", tt.inputErr, err, tt.name)
			if tt.inputErr != nil {
				return
			}

			// Unmarshal
			var output interface{}
			if tt.outputPtr != nil {
				output = reflect.New(reflect.TypeOf(tt.outputPtr).Elem()).Interface()
			}

			err = parser.Unmarshal(input, output)
			assert.Equalf(t, tt.outputErr, err, "unmarshal error expected %v, got %v for %s", tt.outputErr, err, tt.name)
			if tt.outputErr != nil {
				return
			}

			assert.Equalf(t, tt.outputPtr, output, "unmarshal result expected %v, got %v for %s", tt.outputPtr, output, tt.name)
		})
	}
}
