package multiparser

import (
	"errors"
	"sync"
)

var _ Converter = &multiconverter{}

var (
	ErrMarshal   = errors.New("no Marshaller could convert data")
	ErrUnmarshal = errors.New("no Unmarshaller could convert data")
)

type multiconverter struct {
	converters []Converter
}

func (m *multiconverter) Marshal(object interface{}) ([]byte, error) {
	resultCh := make(chan []byte, len(m.converters))

	// try every Marshaller concurrently, if it succeeds, that's our result
	for _, conv := range m.converters {
		go func(conv Converter) {
			data, err := conv.Marshal(object)
			if err != nil {
				resultCh <- nil
				return
			}
			resultCh <- data
		}(conv)
	}

	// get the result
	var result []byte
	for i := 0; i < len(m.converters); i++ {
		data := <-resultCh
		if data != nil {
			result = data
			break
		}
	}

	// check
	if result == nil {
		return nil, ErrMarshal
	}
	return result, nil
}

func (m *multiconverter) Unmarshal(from []byte, to interface{}) error {
	var result interface{}
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}

	// try every Unmarshaller concurrently, if it succeeds, that's our result
	for _, conv := range m.converters {
		wg.Add(1)
		go func(conv Converter) {
			defer wg.Done()

			var copyObj interface{}
			if conv.Unmarshal(from, &copyObj) == nil {
				mu.Lock()
				result = copyObj
				mu.Unlock()
			}
		}(conv)
	}

	wg.Wait()
	if result == nil {
		return ErrUnmarshal
	}
	to = result
	return nil
}

func New(converters ...Converter) Converter {
	return &multiconverter{
		converters: converters,
	}
}
