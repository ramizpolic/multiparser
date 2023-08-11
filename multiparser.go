package multiparser

import (
	"errors"
	"reflect"
	"sync"
)

var (
	ErrMarshal       = errors.New("no Marshaller could convert data")
	ErrUnmarshal     = errors.New("no Unmarshaller could convert data")
	ErrInvalidObject = errors.New("object must be an initialized pointer")
)

type multiconverter struct {
	converters []Converter
}

func New(converters ...Converter) Converter {
	return &multiconverter{
		converters: converters,
	}
}

func (m *multiconverter) Marshal(object interface{}) ([]byte, error) {
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}

	// Try every Marshaller concurrently, if it succeeds, that's our result
	var result []byte
	for _, conv := range m.converters {
		wg.Add(1)
		go func(conv Converter) {
			defer wg.Done()

			// Marshal
			data, err := conv.Marshal(object)
			if err == nil {
				mu.Lock()
				if result == nil {
					result = data
				}
				mu.Unlock()
			}
		}(conv)
	}

	// Wait
	wg.Wait()

	// Check
	if result == nil {
		return nil, ErrMarshal
	}
	return result, nil
}

func (m *multiconverter) Unmarshal(from []byte, to interface{}) error {
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}

	// Validate
	if rv := reflect.ValueOf(to); to == nil || rv.Kind() != reflect.Pointer || rv.IsNil() {
		return ErrInvalidObject
	}

	// Try every Unmarshaller concurrently, if it succeeds, that's our result
	var resultPtr interface{}
	for _, conv := range m.converters {
		wg.Add(1)
		go func(conv Converter) {
			defer wg.Done()

			// Unmarshal
			fromPtr := reflect.New(reflect.TypeOf(to).Elem()).Interface()
			if conv.Unmarshal(from, fromPtr) == nil {
				mu.Lock()
				if resultPtr == nil {
					resultPtr = fromPtr
				}
				mu.Unlock()
			}
		}(conv)
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
