package typex

import (
	"fmt"
	"runtime"
	"strings"
)

type ErrorSource string

const (
	Unknown ErrorSource = "Unknown Source"
)

// Error represents an error that can occur during the runtime.
type Error struct {
	source      ErrorSource
	code        ErrorCode
	description string
	errorStack  []error
	stack       string
}

// New creates an error
func New(source ErrorSource, code ErrorCode, description string) *Error {
	return &Error{
		source:      source,
		code:        code,
		description: description,
	}
}

func (e *Error) Source() ErrorSource {
	return e.source
}

// Code returns the code from the error.
func (e *Error) Code() ErrorCode {
	return e.code
}

// Error only returns the header of the error, not the error stack.
func (e *Error) Error() string {
	return fmt.Sprintf("error_source=%q code=\"(%d:%s)\" fmt=%q", e.source,
		e.code.HTTPStatusCode(), e.code.String(), e.description)
}

// With allows the concatination of errors to produce a longer error stack. This
// is useful when you want to understand the origin of errors when you're taking
// ownership of the errors.
func (e *Error) With(errs ...error) *Error {
	e.errorStack = append(e.errorStack, errs...)
	return e
}

// ErrorStack returns a slice of all the errors associated with the error.
func (e *Error) ErrorStack() []error {
	return e.errorStack
}

// IncludeStack allows you to include the stack which is useful when bubbling
// up and want to know the original issue.
func (e *Error) IncludeStack() *Error {
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, false)
	e.stack = string(buf)

	return e
}

// Is returns if the error is of the type ErrorName
func (e *Error) Is(name ErrorName) bool {
	return e.code.Is(name)
}

// Inspect goes through the whole errors and tries to print them out in a
// reliable manor!
func (e *Error) Inspect() string {
	return fmt.Sprintf("> id=0 %s", e.inspect(1))
}

func (e *Error) inspect(depth int) string {
	stack := []string{
		e.Error(),
	}
	for k, v := range e.errorStack {
		// Implosion alert!
		if v == e {
			continue
		}

		var val string
		if a, ok := v.(*Error); ok {
			val = a.inspect(depth + 1)
		} else if a, ok := v.(error); ok {
			val = fmt.Sprintf("error_source=%q native=\"true\"", a.Error())
		} else {
			continue
		}

		stack = append(stack, fmt.Sprintf("%s> id=%d %s", strings.Repeat("\t", depth), k, val))
	}

	if len(e.stack) > 0 {
		stack = append(stack, formatStack(e.stack, depth))
	}

	return strings.Join(stack, "\n")
}

func formatStack(stack string, depth int) string {
	var (
		lines = strings.Split(stack, "\n")
		res   = make([]string, len(lines))

		header  = "> Runtime Stack:"
		largest = len(header)
		offset  = strings.Repeat("\t", depth)
	)
	for k, v := range lines {
		val := fmt.Sprintf("%s| %s", offset, v)
		res[k] = strings.TrimRightFunc(val, func(r rune) bool {
			if r < ' ' {
				return true
			}
			return false

		})
		if num := len(res[k]); num > largest {
			largest = num
		}
	}
	return fmt.Sprintf(
		"%s%s\n%s%s\n%s\n%s%s",
		offset,
		header,
		offset,
		strings.Repeat("-", largest),
		strings.Join(res, "\n"),
		offset,
		strings.Repeat("-", largest),
	)
}

// Errorf formats according to a format specifier and returns the string as a
// value that satisfies error.
func Errorf(source ErrorSource, code ErrorCode, format string, a ...interface{}) *Error {
	return New(source, code, fmt.Sprintf(format, a...))
}

// ErrCode returns the error code from the potential error. If it's an Error it
// will return the code, otherwise it'll return -1.
func ErrCode(err error) int {
	if e, ok := err.(*Error); ok {
		return e.Code().HTTPStatusCode()
	}
	return -1
}

// ErrName returns the error name from the potential error. If it's an Error it
// will return the name, otherwise it'll return "".
func ErrName(err error) string {
	if e, ok := err.(*Error); ok {
		return e.Code().String()
	}
	return ""
}

// Inspect returns the full inspection of the error if it's an Error otherwise
// it'll return the result from `.Error()`
func Inspect(err error) string {
	if e, ok := err.(*Error); ok {
		return e.Inspect()
	}
	return err.Error()
}

// Is returns true if the error name matches
func Is(err error, name ErrorName) bool {
	if e, ok := err.(*Error); ok {
		return e.Is(name)
	}
	return false
}

// As converts a string to a error name
func As(name string) ErrorName {
	return runes(name)
}

// Lift a normal error into the whole Error pipeline.
func Lift(err error) *Error {
	if e, ok := err.(*Error); ok {
		return e
	}

	return Errorf(Unknown, InternalServerError, err.Error())
}
