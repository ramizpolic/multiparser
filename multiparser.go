package multiparser

import (
	"errors"
	"fmt"
	"reflect"
)

// Parser implements raw data to object deserialization.
type Parser interface {
	Parse(from []byte, to interface{}) error
}

var (
	ErrEmptyParsers  = errors.New("no Parser passed, at least one required")
	ErrInvalidObject = errors.New("object must be non-nil pointer")
)

type multiParser struct {
	parsers []Parser
}

// New creates a new Parser which can be used to serialize data from different
// formats.
func New(parsers ...Parser) (Parser, error) {
	if len(parsers) == 0 {
		return nil, ErrEmptyParsers
	}
	return &multiParser{
		parsers: parsers,
	}, nil
}

func (p *multiParser) Parse(from []byte, to interface{}) error {
	// Validate
	if rv := reflect.ValueOf(to); to == nil || rv.Kind() != reflect.Pointer || rv.IsNil() {
		return ErrInvalidObject
	}

	// Try every Parser, if any of them succeeds, that's our result.
	// The assumption is that Parsing will fail fast.
	var err error
	for _, parser := range p.parsers {
		fromPtr := reflect.New(reflect.TypeOf(to).Elem()).Interface()
		parseErr := parser.Parse(from, fromPtr)
		if parseErr == nil {
			// Update result
			reflect.ValueOf(to).Elem().Set(reflect.ValueOf(fromPtr).Elem())
			return nil
		} else {
			err = errors.Join(err, &ParserError{
				Parser: fmt.Sprintf("%T", parser),
				Err:    parseErr,
			})
		}
	}

	return &ParserError{
		Parser: fmt.Sprintf("%T", p),
		Err:    err,
	}
}

type ParserError struct {
	Parser string
	Err    error
}

func (e *ParserError) Error() string {
	return fmt.Sprintf("parsing failed for %s: %v", e.Parser, e.Err)
}

func (e *ParserError) Unwrap() error {
	return e.Err
}
