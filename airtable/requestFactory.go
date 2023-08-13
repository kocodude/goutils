package airtable

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type requestFactory[T any] struct {
	tableClient *TableClient[T]
}

func (r requestFactory[T]) createRecords(collection *RecordCollection[T]) (*http.Request, error) {

	body, err := json.Marshal(collection)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", r.tableClient.createRecordsURLString, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", r.tableClient.apiKeyString)
	request.Header.Set("Content-Type", "application/json")

	return request, nil
}

func (r requestFactory[T]) listRecords(filterByFormula string) (*http.Request, error) {

	request, err := http.NewRequest("GET", r.tableClient.listRecordsURLString, nil)
	if err != nil {
		return nil, err
	}

	queryParams := request.URL.Query()
	queryParams.Add("filterByFormula", filterByFormula)
	request.URL.RawQuery = queryParams.Encode()

	request.Header.Set("Authorization", r.tableClient.apiKeyString)

	return request, nil
}
