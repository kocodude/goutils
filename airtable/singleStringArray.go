package airtable

import (
	"fmt"

	"github.com/kocodude/goutils/foundation"
)

type singleStringArray []string

type singleStringArrayFactory struct{}

var sharedSingleStringArrayFactory = singleStringArrayFactory{}

func (ssa singleStringArray) string() (string, error) {

	if len(ssa) != 1 {
		return "", foundation.RavelError{
			Message: fmt.Sprintf(
				"cannot get item from single item array %v",
				ssa),
		}
	}

	return ssa[0], nil
}

func (ssa singleStringArrayFactory) createFromString(stringValue string) singleStringArray {
	return []string{stringValue}
}
