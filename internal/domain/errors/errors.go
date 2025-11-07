package errors

import "fmt"

var (
	ErrDepNotFound = fmt.Errorf("job dependency not found")
)
