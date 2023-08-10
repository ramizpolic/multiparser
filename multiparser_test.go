package multiparser_test

import (
	"github.com/ramizpolic/multiparser"
	"github.com/ramizpolic/multiparser/parser/json"
	"testing"
)

func TestNew(t *testing.T) {
	type jsonObj struct {
		Data string `json:"data"`
	}

	// JSON Parser
	jsonParser := multiparser.New(json.Converter)

	// Marshal
	result, _ := jsonParser.Marshal(jsonObj{Data: "data"})
	if string(result) != "{\"data\":\"data\"}" {
		t.Fatalf("Marshal failed")
	}

	// Unmarshal
	var obj jsonObj
	_ = jsonParser.Unmarshal(result, &obj)
	if obj.Data != "data" {
		t.Fatalf("Unmarshal failed %s", obj.Data)
	}
}
