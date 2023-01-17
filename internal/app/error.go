package app

import (
	"errors"
	"fmt"
)

// InvalidCommandError should be returned by the implementations of the interface when the handler does not receive the needed command.
type InvalidCommandError struct {
	expected string
	had      string
}

// NewInvalidCommandError is a constructor
func NewInvalidCommandError(expected string, had string) InvalidCommandError {
	return InvalidCommandError{expected: expected, had: had}
}

const errMsgInvalidCommand = "invalid command, expected '%s' but found '%s'"

func (e InvalidCommandError) Error() string {
	return fmt.Sprintf(errMsgInvalidCommand, e.expected, e.had)
}

// InvalidQueryError is self described
type InvalidQueryError struct {
	Expected string
	Had      string
}

const errMsgInvalidQuery = "invalid query, expected '%s' but found '%s'"

// Error implements the error.Error interface
func (ewq InvalidQueryError) Error() string {
	return fmt.Sprintf(errMsgInvalidQuery, ewq.Expected, ewq.Had)
}

// NewInvalidQueryError is a constructor
func NewInvalidQueryError(expected, had string) InvalidQueryError {
	return InvalidQueryError{Expected: expected, Had: had}
}

// ErrNotFound is an entity not found error
var ErrNotFound = errors.New("not found")
