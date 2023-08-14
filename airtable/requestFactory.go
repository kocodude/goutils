package airtable

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kocodude/goutils/foundation"
	"net/http"
)

func (rf requestFactory[T]) CreateRecords(collection *RecordCollection[T]) (*http.Request, error) {

	body, err := json.Marshal(collection)
	if err != nil {
		return nil, foundation.RavelError{
			Err: err,
			Message: fmt.Sprintf(
				"failed to marshal into json, collection: %v",
				collection),
		}
	}

	request, err := http.NewRequest("POST", rf.tableClient.createRecordsURLString, bytes.NewBuffer(body))
	if err != nil {
		return nil, foundation.RavelError{
			Err: err,
			Message: fmt.Sprintf(
				"failed to create request\nstring: %v\nbody: %v",
				rf.tableClient.createRecordsURLString,
				body),
		}
	}

	request.Header.Set("Content-Type", "application/json")

	return rf.requestWithAuthorization(request)
}

func (rf requestFactory[T]) ListFilteredRecords(filterByFormula string) (*http.Request, error) {

	request, err := rf.createListRequest()
	if err != nil {
		return nil, err
	}

	queryParams := request.URL.Query()
	queryParams.Add("filterByFormula", filterByFormula)
	request.URL.RawQuery = queryParams.Encode()

	return rf.requestWithAuthorization(request)
}

func (rf requestFactory[T]) ListRecords() (*http.Request, error) {

	request, err := rf.createListRequest()
	if err != nil {
		return nil, err
	}

	return rf.requestWithAuthorization(request)
}

func (rf requestFactory[T]) ListSortedRecords(sortSpecifier string) (*http.Request, error) {

	request, err := rf.createListRequest()
	if err != nil {
		return nil, err
	}

	queryParams := request.URL.Query()
	queryParams.Add("sort", sortSpecifier)
	request.URL.RawQuery = queryParams.Encode()

	return rf.requestWithAuthorization(request)
}

type requestFactory[T any] struct {
	tableClient *TableClient[T]
}

func (rf requestFactory[T]) createListRequest() (*http.Request, error) {

	request, err := http.NewRequest("GET", rf.tableClient.listRecordsURLString, nil)
	if err != nil {
		return nil, foundation.RavelError{
			Err: err,
			Message: fmt.Sprintf(
				"failed to create new request using string: %v",
				rf.tableClient.listRecordsURLString),
		}
	}

	return request, nil
}

func (rf requestFactory[T]) requestWithAuthorization(request *http.Request) (*http.Request, error) {

	request.Header.Set("Authorization", rf.tableClient.apiKeyString)

	return request, nil
}
