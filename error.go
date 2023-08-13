package foundation

import "fmt"

type RavelError struct {
	Err     error
	Message string
}

func (r RavelError) Error() string {
	return fmt.Sprintf("RavelError: %v ...wrapping:\n\t%v", r.Message, r.Err)
}
