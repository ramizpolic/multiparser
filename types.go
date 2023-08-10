package multiparser

type Marshaller interface {
	Marshal(object interface{}) ([]byte, error)
}

type Unmarshaller interface {
	Unmarshal(from []byte, to interface{}) error
}

type Converter interface {
	Marshaller
	Unmarshaller
}
