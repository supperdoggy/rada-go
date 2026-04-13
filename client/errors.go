package client

import (
	"errors"

	"github.com/supperdoggy/vr_api/internal/apperr"
)

type (
	ParseError = apperr.Error
	ErrorCode  = apperr.Code
)

const (
	ErrorCodeLayoutChanged      = apperr.CodeLayoutChanged
	ErrorCodeMissingField       = apperr.CodeMissingField
	ErrorCodeUnsupportedProfile = apperr.CodeUnsupportedProfile
)

var (
	ErrLayoutChanged      = apperr.ErrLayoutChanged
	ErrMissingField       = apperr.ErrMissingField
	ErrUnsupportedProfile = apperr.ErrUnsupportedProfile

	ErrNotImplemented = errors.New("not implemented")
)
