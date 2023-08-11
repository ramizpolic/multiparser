package multiparser

import (
	"errors"
	"reflect"
	"sync"
)

var (
	ErrEmptyParsers  = errors.New("no Parser passed, at least one required")
	ErrMarshal       = errors.New("no Marshaller could convert data")
	ErrUnmarshal     = errors.New("no Unmarshaller could convert data")
	ErrInvalidObject = errors.New("object must be non-nil pointer")
)

type multiParser struct {
	parsers []Parser
}

// New creates a new Parser which can be used to serialize and deserialize
// data with different formats.
func New(parsers ...Parser) (Parser, error) {
	if len(parsers) == 0 {
		return nil, ErrEmptyParsers
	}
	return &multiParser{
		parsers: parsers,
	}, nil
}

func (m *multiParser) Marshal(object interface{}) ([]byte, error) {
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}

	// Try every Marshaller concurrently, if it succeeds, that's our result
	var result []byte
	for _, parser := range m.parsers {
		wg.Add(1)
		go func(parser Parser) {
			defer wg.Done()

			// Marshal
			data, err := parser.Marshal(object)
			if err == nil {
				mu.Lock()
				if result == nil {
					result = data
				}
				mu.Unlock()
			}
		}(parser)
	}

	// Wait
	wg.Wait()

	// Check
	if result == nil {
		return nil, ErrMarshal
	}
	return result, nil
}

func (m *multiParser) Unmarshal(from []byte, to interface{}) error {
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}

	// Validate
	if rv := reflect.ValueOf(to); to == nil || rv.Kind() != reflect.Pointer || rv.IsNil() {
		return ErrInvalidObject
	}

	// Try every Unmarshaller concurrently, if it succeeds, that's our result
	var resultPtr interface{}
	for _, parser := range m.parsers {
		wg.Add(1)
		go func(parser Parser) {
			defer wg.Done()

			// Unmarshal
			fromPtr := reflect.New(reflect.TypeOf(to).Elem()).Interface()
			if parser.Unmarshal(from, fromPtr) == nil {
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
		return ErrUnmarshal
	}

	// Update
	reflect.ValueOf(to).Elem().Set(reflect.ValueOf(resultPtr).Elem())

	return nil
}
