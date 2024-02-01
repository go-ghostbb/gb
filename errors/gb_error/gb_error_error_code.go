package gberror

import gbcode "github.com/Ghostbb-io/gb/errors/gb_code"

// Code returns the error code.
// It returns CodeNil if it has no error code.
func (err *Error) Code() gbcode.Code {
	if err == nil {
		return gbcode.CodeNil
	}
	if err.code == gbcode.CodeNil {
		return Code(err.Unwrap())
	}
	return err.code
}

// SetCode updates the internal code with given code.
func (err *Error) SetCode(code gbcode.Code) {
	if err == nil {
		return
	}
	err.code = code
}
