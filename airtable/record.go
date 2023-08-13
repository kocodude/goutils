package airtable

type Record[T any] struct {
	Identifier string `json:"id,omitempty"`
	Fields     T      `json:"fields"`
}

type RecordFactory[T any] struct{}

func (rf RecordFactory[T]) CreteFromFields(fields T) *Record[T] {
	return &Record[T]{
		Fields: fields,
	}
}
