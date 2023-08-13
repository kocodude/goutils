package airtable

import (
	"github.com/kocodude/goutils/foundation"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TableClient[T any] struct {
	apiKeyString                   string
	createRecordsURLString         string
	listRecordsURLString           string
	retrieveOneRecordURLStringBase string
	requestFactory                 requestFactory[T]
}

type TableClientFactory[T any] struct{}

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

	request, err := tc.requestFactory.createRecords(records)
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

func (tc TableClient[T]) ListRecords(filterByFormula string) (*RecordCollection[T], error) {

	request, err := tc.requestFactory.listRecords(filterByFormula)
	if err != nil {
		return nil, foundation.RavelError{
			Err: err,
			Message: fmt.Sprintf(
				"failed to build list records request using formula %v",
				filterByFormula),
		}
	}

	return tc.executeRequestForRecords(request)
}

func (tc TableClient[T]) RetrieveOneRecord(reference []string) (*Record[T], error) {

	request, err := tc.requestFactory.retrieveOneRecord(reference)
	if err != nil {
		return nil, foundation.RavelError{
			Err: err,
			Message: fmt.Sprintf(
				"failed to build retrieve one Record request using %v",
				reference),
		}
	}

	var record Record[T]
	err = tc.executeRequestAndUnmarshalBody(request, &record)
	if err != nil {
		return nil, foundation.RavelError{
			Err: err,
			Message: fmt.Sprintf(
				"failed to execute retrieve one Record request %v and unmarshal body",
				request),
		}
	}

	return &record, nil
}

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

func (tc TableClient[T]) executeRequestAndUnmarshalBody(request *http.Request, unmarshalledBody interface{}) error {

	client := &http.Client{}

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

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, unmarshalledBody)
}

func (tc TableClient[T]) executeRequestForOneRecord(request *http.Request) (*Record[T], error) {

	var createdRecord Record[T]
	err := tc.executeRequestAndUnmarshalBody(request, &createdRecord)
	if err != nil {
		return nil, foundation.RavelError{
			Err: err,
			Message: fmt.Sprintf(
				"failed to execute retrieve one Record request %v and unmarshal body",
				request),
		}
	}

	return &createdRecord, nil
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
