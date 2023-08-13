package foundation

import "fmt"

type RavelError struct {
	Err     error
	Message string
}

func (r RavelError) Error() string {

	if r.Err != nil {
		return fmt.Sprintf("RavelError: %v ...wrapping:\n\t%v", r.Message, r.Err)
	}

	return fmt.Sprintf("RavelError: %v", r.Message)
}
