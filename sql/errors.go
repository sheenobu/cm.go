package sql

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// CodedError is an error with an attached code
type CodedError interface {
	Code() int
}

// GetErrorCode extracts out the error code for the given object
func GetErrorCode(err error) int {
	type causer interface {
		Cause() error
	}

	for err != nil {
		coded, ok := err.(CodedError)
		if ok {
			return coded.Code()
		}

		cause, ok := err.(causer)
		if !ok {
			// done with stack
			break
		}
		err = cause.Cause()
	}

	return CodeNotAvailable
}

// Error code constants
const (
	CodeNotAvailable   int = 0
	TableAlreadyExists int = 1
)

func mapError(err error) error {
	if err == nil {
		return err
	}

	//TODO: need to enumerate all the different 'already exist' codes for each driver
	if strings.Contains(errors.Cause(err).Error(), "already exists") {
		return &codedError{
			C:           TableAlreadyExists,
			Description: "Table already exists",
			Err:         err,
		}
	}

	//TODO: enumerate other errors with their enumerated types above

	return err
}

type codedError struct {
	C           int
	Description string
	Err         error
}

func (err *codedError) Code() int {
	return err.C
}

func (err *codedError) Error() string {
	return fmt.Sprintf("%s", err.Description)
}

func (err *codedError) Cause() error {
	return err.Err
}
