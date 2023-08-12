package multiparser_test

import (
	"github.com/ramizpolic/multiparser"
	"github.com/ramizpolic/multiparser/parser"
	"github.com/stretchr/testify/assert"
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
		expected map[string]string
		err      string
	}{
		{
			name:     "json",
			inputRaw: `{"data": "data"}`,
			expected: map[string]string{"data": "data"},
		},
		{
			name:     "yaml",
			inputRaw: `data: data`,
			expected: map[string]string{"data": "data"},
		},
		{
			name:     "invalid",
			inputRaw: `INVALID RAW DATA`,
			err:      "parsing failed for *multiparser.multiParser",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var got map[string]string
			err := parser.Parse([]byte(tt.inputRaw), &got)
			if tt.err != "" {
				assert.ErrorContainsf(t, err, tt.err, "parse error expected %v, got %v for %s", tt.err, err, tt.name)
			} else {
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestAll(t *testing.T) {
	for _, tt := range []struct {
		name        string
		parsers     []multiparser.Parser
		newErr      string
		input       string
		parseErr    string
		parseResult map[string]interface{}
	}{
		{
			name:   "nil-prs",
			newErr: multiparser.ErrEmptyParsers.Error(),
		},
		{
			name:        "all-prs-json-map",
			parsers:     parsers,
			input:       `{"data": "data"}`,
			parseResult: map[string]interface{}{"data": "data"},
		},
		{
			name:        "all-prs-json-slice",
			parsers:     parsers,
			input:       `{"data": ["data"]}`,
			parseResult: map[string]interface{}{"data": []interface{}{"data"}},
		},
		{
			name:        "all-prs-yaml-map",
			parsers:     parsers,
			input:       `data: data`,
			parseResult: map[string]interface{}{"data": "data"},
		},
		{
			name:        "all-prs-yaml-slice",
			parsers:     parsers,
			input:       `data: [data]`,
			parseResult: map[string]interface{}{"data": []interface{}{"data"}},
		},
		{
			name:        "all-prs-duplicate",
			parsers:     []multiparser.Parser{parser.JSON, parser.JSON, parser.JSON},
			input:       `{"data": "data"}`,
			parseResult: map[string]interface{}{"data": "data"},
		},
		{
			name:    "all-prs-input-empty",
			parsers: parsers,
		},
		{
			name:    "all-prs-input-invalid",
			parsers: parsers,
			input:   `------------------`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// Create
			parser, err := multiparser.New(tt.parsers...)
			if tt.newErr != "" {
				assert.ErrorContainsf(t, err, tt.newErr, "new error expected %v, got %v for %s", tt.newErr, err, tt.name)
				return
			}

			// Parse
			var output map[string]interface{}
			err = parser.Parse([]byte(tt.input), &output)
			if tt.parseErr != "" {
				assert.ErrorContainsf(t, err, tt.parseErr, "parse error expected %v, got %v for %s", tt.parseErr, err, tt.name)
				return
			}

			assert.Equalf(t, tt.parseResult, output, "parse result expected %v, got %v for %s", tt.parseResult, output, tt.name)
		})
	}

}
