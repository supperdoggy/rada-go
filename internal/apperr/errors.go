package apperr

import "fmt"

// Code describes a typed parser/client error code.
type Code string

const (
	CodeLayoutChanged      Code = "layout_changed"
	CodeMissingField       Code = "missing_field"
	CodeUnsupportedProfile Code = "unsupported_profile"
)

// Error is a typed error that supports errors.Is checks by code.
type Error struct {
	Code   Code
	Op     string
	Detail string
	Cause  error
}

func (e *Error) Error() string {
	base := string(e.Code)
	if e.Op != "" {
		base = fmt.Sprintf("%s [%s]", base, e.Op)
	}
	if e.Detail != "" {
		base = fmt.Sprintf("%s: %s", base, e.Detail)
	}
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", base, e.Cause)
	}
	return base
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func (e *Error) Is(target error) bool {
	typed, ok := target.(*Error)
	if !ok {
		return false
	}
	return e.Code == typed.Code
}

var (
	ErrLayoutChanged      = &Error{Code: CodeLayoutChanged}
	ErrMissingField       = &Error{Code: CodeMissingField}
	ErrUnsupportedProfile = &Error{Code: CodeUnsupportedProfile}
)

func NewLayoutChangedError(op, detail string, cause error) error {
	return &Error{
		Code:   CodeLayoutChanged,
		Op:     op,
		Detail: detail,
		Cause:  cause,
	}
}

func NewMissingFieldError(op, detail string, cause error) error {
	return &Error{
		Code:   CodeMissingField,
		Op:     op,
		Detail: detail,
		Cause:  cause,
	}
}

func NewUnsupportedProfileError(op, detail string, cause error) error {
	return &Error{
		Code:   CodeUnsupportedProfile,
		Op:     op,
		Detail: detail,
		Cause:  cause,
	}
}
