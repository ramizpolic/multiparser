package multiparser_test

import (
	"github.com/ramizpolic/multiparser"
	"github.com/ramizpolic/multiparser/parser/json"
	"github.com/ramizpolic/multiparser/parser/yaml"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

// register all converters
var converters = []multiparser.Converter{
	json.Converter,
	yaml.Converter,
}

// TestMultiparser will test against all converters
func TestMultiparser(t *testing.T) {
	t.Parallel()

	type objType struct {
		Data string `json:"data" yaml:"data"`
	}

	for _, tt := range []struct {
		name            string
		convs           []multiparser.Converter
		mInput          interface{}
		mResult         []byte
		mErr            string
		umOverrideInput []byte
		umResultPtr     interface{}
		umErr           string
	}{
		{
			name:        "multiparser-nil",
			convs:       converters,
			mInput:      nil,
			mResult:     []byte(`null`),
			umResultPtr: &struct{}{},
		},
		{
			name:        "multiparser-json",
			convs:       converters,
			mInput:      objType{Data: "data"},
			mResult:     []byte(`{"data":"data"}`),
			umResultPtr: &objType{Data: "data"},
		},
		{
			name:   "multiparser-yaml",
			convs:  converters,
			mInput: objType{Data: "data"},
			mResult: []byte(`data: data
`),
			umResultPtr: &objType{Data: "data"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Multiparser
			jsonParser := multiparser.New(tt.convs...)

			// Marshal
			mResult, err := jsonParser.Marshal(tt.mInput)
			if tt.mErr != "" {
				assert.ErrorContainsf(t, err, tt.mErr, "marshal error expected %v, got %v for %s", tt.mErr, err, tt.name)
			} else {
				assert.Nilf(t, err, "marshal error must be nil for %s", tt.name)
			}
			assert.Equalf(t, tt.mResult, mResult, "marshal result expected %v, got %v for %s", tt.mResult, mResult, tt.name)

			// Unmarshal
			umInput := tt.umOverrideInput
			if umInput == nil {
				umInput = mResult
			}
			umResult := reflect.New(reflect.TypeOf(tt.umResultPtr).Elem()).Interface()
			err = jsonParser.Unmarshal(umInput, umResult)
			if tt.umErr != "" {
				assert.ErrorContainsf(t, err, tt.umErr, "unmarshal error expected %v, got %v for %s", tt.umErr, err, tt.name)
			} else {
				assert.Nilf(t, err, "unmarshal error must be nil for %s", tt.name)
			}
			assert.Equalf(t, tt.umResultPtr, umResult, "unmarshal result expected %v, got %v for %s", tt.umResultPtr, umResult, tt.name)

		})
	}
}
