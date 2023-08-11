package multiparser

import (
	"errors"
	"reflect"
	"sync"
)

// Parser implements raw data to object deserialization.
type Parser interface {
	Parse(from []byte, to interface{}) error
}

var (
	ErrEmptyParsers  = errors.New("no Parser passed, at least one required")
	ErrParse         = errors.New("no Parser could convert data")
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

func (m *multiParser) Parse(from []byte, to interface{}) error {
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}

	// Validate
	if rv := reflect.ValueOf(to); to == nil || rv.Kind() != reflect.Pointer || rv.IsNil() {
		return ErrInvalidObject
	}

	// Try every Parser concurrently, if it succeeds, that's our result
	var resultPtr interface{}
	for _, parser := range m.parsers {
		wg.Add(1)
		go func(parser Parser) {
			defer wg.Done()

			// Parse
			fromPtr := reflect.New(reflect.TypeOf(to).Elem()).Interface()
			if parser.Parse(from, fromPtr) == nil {
				mu.Lock()
				if resultPtr == nil {
					resultPtr = fromPtr
				}
				mu.Unlock()
			}
		}(parser)
	}

	// Wait
	wg.Wait()

	// Check
	if resultPtr == nil {
		return ErrParse
	}

	// Update
	reflect.ValueOf(to).Elem().Set(reflect.ValueOf(resultPtr).Elem())

	return nil
}
