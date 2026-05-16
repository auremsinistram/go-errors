package errors

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	stderrors "errors"
)

type Error interface {
	error

	Code() int

	Location() string
	Description() string

	Unwrap() error
}

type customError struct {
	code int

	location    string
	description string

	wrapped error
}

func New(code int, format string, args ...any) error {
	return &customError{
		code:        code,
		location:    getLocation(),
		description: fmt.Sprintf(format, args...),
	}
}

func Errorf(format string, args ...any) error {
	return &customError{
		location:    getLocation(),
		description: fmt.Sprintf(format, args...),
	}
}

func Mark(err error, code int) error {
	return &customError{
		code:     code,
		location: getLocation(),
		wrapped:  err,
	}
}

func Markf(err error, code int, format string, args ...any) error {
	return &customError{
		code:        code,
		location:    getLocation(),
		description: fmt.Sprintf(format, args...),
		wrapped:     err,
	}
}

func Wrap(err error, description string) error {
	return &customError{
		location:    getLocation(),
		description: description,
		wrapped:     err,
	}
}

func Wrapf(err error, format string, args ...any) error {
	return &customError{
		location:    getLocation(),
		description: fmt.Sprintf(format, args...),
		wrapped:     err,
	}
}

func GetCode(err error) (int, bool) {
	var e Error

	for err != nil {
		if As(err, &e) {
			code := e.Code()

			if code != 0 {
				return code, true
			}
		}

		err = Unwrap(err)
	}

	return 0, false
}

func GetRoot(err error) error {
	for err != nil {
		if e := Unwrap(err); e != nil {
			err = e

			continue
		}

		return err
	}

	return nil
}

func Is(err error, target error) bool {
	return stderrors.Is(err, target)
}

func As(err error, target any) bool {
	return stderrors.As(err, target)
}

func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}

func Join(errs ...error) error {
	return stderrors.Join(errs...)
}

func getLocation() string {
	_, file, line, _ := runtime.Caller(2)

	return fmt.Sprintf(
		"%s:%d",
		file,
		line,
	)
}

func (e *customError) Error() string {
	var builder strings.Builder

	first := true

	writeSeparator := func() {
		if !first {
			builder.WriteString(", ")
		}

		first = false
	}

	if e.code != 0 {
		writeSeparator()
		builder.WriteString("code: ")
		builder.WriteString(strconv.Itoa(e.code))
	}

	if e.location != "" {
		writeSeparator()
		builder.WriteString("location: ")
		builder.WriteString(e.location)
	}

	if e.description != "" {
		writeSeparator()
		builder.WriteString("description: ")
		builder.WriteString(e.description)
	}

	if e.wrapped != nil {
		writeSeparator()
		builder.WriteString("wrapped: ")
		builder.WriteString(e.wrapped.Error())
	}

	return builder.String()
}

func (e *customError) Code() int {
	return e.code
}

func (e *customError) Location() string {
	return e.location
}

func (e *customError) Description() string {
	return e.description
}

func (e *customError) Unwrap() error {
	return e.wrapped
}
