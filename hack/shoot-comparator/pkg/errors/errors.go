package errors

import "fmt"

var (
	ErrInvalidType  = fmt.Errorf("invalid type")
	ErrInvalidValue = fmt.Errorf("invalid value")
	ErrNilValue     = fmt.Errorf("%w: nil", ErrInvalidValue)
)
