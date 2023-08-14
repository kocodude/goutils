package airtable

import (
	"github.com/kocodude/goutils/foundation"
	log "github.com/sirupsen/logrus"
	"io"

	"encoding/json"
	"fmt"
	"net/http"
)

// TableClient is the main interface for communicating with an Airtable Table
type TableClient[T any] struct {
	apiKeyString                   string
	createRecordsURLString         string
	listRecordsURLString           string
	retrieveOneRecordURLStringBase string
	requestFactory                 requestFactory[T]
}

func (tc TableClient[T]) CreateOneRecord(record *Record[T]) (*Record[T], error) {

	collection := RecordCollectionFactory[T]{RecordFactory: RecordFactory[T]{}}.CreateCollectionFromRecord(record)

	createdCollection, err := tc.CreateRecords(collection)
	if err != nil {
		return nil, foundation.RavelError{
			Err: err,
			Message: fmt.Sprintf(
				"failed to create records using collection %v",
				collection),
		}
	}

	if len(createdCollection.Records) > 0 {
		return &createdCollection.Records[0], nil
	}

	return nil, foundation.RavelError{
		Message: fmt.Sprintf(
			"failed to create Record in airtable and recieved empty records result trying to create Record %v",
			record),
	}
}

func (tc TableClient[T]) CreateRecords(records *RecordCollection[T]) (*RecordCollection[T], error) {

	request, err := tc.requestFactory.CreateRecords(records)
	if err != nil {
		return nil, foundation.RavelError{
			Err: err,
			Message: fmt.Sprintf(
				"failed to build create records request using records %v",
				records),
		}
	}

	return tc.executeRequestForRecords(request)
}

func (tc TableClient[T]) ListFilteredRecords(filterFormula string) (*RecordCollection[T], error) {

	request, err := tc.requestFactory.ListFilteredRecords(filterFormula)
	if err != nil {
		return nil, foundation.RavelError{
			Err: err,
			Message: fmt.Sprintf(
				"failed to create list records request\nformula: %v",
				filterFormula),
		}
	}

	return tc.executeRequestForRecords(request)
}

func (tc TableClient[T]) ListRecords() (*RecordCollection[T], error) {

	request, err := tc.requestFactory.ListRecords()
	if err != nil {
		return nil, foundation.RavelError{
			Err:     err,
			Message: "failed to create list records request",
		}
	}

	return tc.executeRequestForRecords(request)
}

func (tc TableClient[T]) ListSortedRecords(field string, direction string) (*RecordCollection[T], error) {

	request, err := tc.requestFactory.ListSortedRecords(field, direction)
	if err != nil {
		return nil, foundation.RavelError{
			Err:     err,
			Message: "failed to create list records request",
		}
	}

	return tc.executeRequestForRecords(request)
}

type TableClientFactory[T any] struct{}

func (tcf TableClientFactory[T]) CreateTableClient(apiKeyString string, createRecordsURLString string, listRecordsURLString string, retrieveOneRecordURLStringBase string) *TableClient[T] {

	tableClient := TableClient[T]{
		apiKeyString:                   apiKeyString,
		createRecordsURLString:         createRecordsURLString,
		listRecordsURLString:           listRecordsURLString,
		retrieveOneRecordURLStringBase: retrieveOneRecordURLStringBase,
	}

	tableClient.requestFactory = requestFactory[T]{
		tableClient: &tableClient,
	}

	return &tableClient
}

// PRIVATE INTERFACE

func (tc TableClient[T]) executeRequestAndUnmarshalBody(request *http.Request, unmarshalledBody interface{}) error {

	client := &http.Client{}

	log.Debugf("Table client will execute request\n\t%v", request)

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if http.StatusOK != response.StatusCode {
		return foundation.RavelError{
			Message: fmt.Sprintf(
				"expected status code %v but received %v from response %v",
				http.StatusOK,
				response.StatusCode,
				response),
		}
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, unmarshalledBody)
}

func (tc TableClient[T]) executeRequestForRecords(request *http.Request) (*RecordCollection[T], error) {

	var records RecordCollection[T]
	err := tc.executeRequestAndUnmarshalBody(request, &records)
	if err != nil {
		return nil, foundation.RavelError{
			Err: err,
			Message: fmt.Sprintf(
				"failed to execute create records request %v and unmarshal body",
				request),
		}
	}

	return &records, nil
}
