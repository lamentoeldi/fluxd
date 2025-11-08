package errors

import "fmt"

var (
	ErrDepNotFound     = fmt.Errorf("job dependency not found")
	ErrCycleFound      = fmt.Errorf("cyclic dependencies are not allowed")
	ErrUnknownExecutor = fmt.Errorf("unknown executor")
	ErrNodeNotFound    = fmt.Errorf("node not found")
	ErrInvalidStatus   = fmt.Errorf("invalid job status")
	ErrStopTraversal   = fmt.Errorf("stop graph traversal")
)
