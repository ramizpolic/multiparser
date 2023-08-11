package multiparser

type Marshaller interface {
	Marshal(object interface{}) ([]byte, error)
}

type Unmarshaller interface {
	Unmarshal(from []byte, to interface{}) error
}

// Parser implements bidirectional parsing (from raw to object and vice versa).
// TODO: consider using different name
type Parser interface {
	Marshaller
	Unmarshaller
}
