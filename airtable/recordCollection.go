package airtable

type RecordCollection[T any] struct {
	Records []Record[T] `json:"records"`
}

func (rc RecordCollection[T]) GetFirstOrNil() *Record[T] {

	if len(rc.Records) > 0 {
		return &rc.Records[0]
	}

	return nil
}

type RecordCollectionFactory[T any] struct {
	RecordFactory RecordFactory[T]
}

func (rcf RecordCollectionFactory[T]) CreateCollectionFromFields(fields T) *RecordCollection[T] {
	return &RecordCollection[T]{
		Records: []Record[T]{*rcf.RecordFactory.CreateFromFields(fields)},
	}
}

func (rcf RecordCollectionFactory[T]) CreateCollectionFromRecord(theRecord *Record[T]) *RecordCollection[T] {
	return &RecordCollection[T]{
		Records: []Record[T]{*theRecord},
	}
}
