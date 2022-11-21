package typex

import "net/http"

// ErrorName is a restricted type alias for errors. It's restricted by design to
// prevent abuse of errors where by the errors become localised and then they're
// about conveying information to users when in fact they're designed to be
// about the state of the application ("machine").
type ErrorName [75]rune

// ErrorCode is a encapsulation of a error code. It exposes a way to identify
// what type of error it is, this is useful for understanding the context.
type ErrorCode struct {
	httpStatus int
	name       ErrorName
}

var (
	BadRequest          = makeErrorCode(http.StatusBadRequest)
	InternalServerError = makeErrorCode(http.StatusInternalServerError)
	NotFound            = makeErrorCode(http.StatusNotFound)
	Unauthorized        = makeErrorCode(http.StatusUnauthorized)
)

func makeErrorCode(code int) ErrorCode {
	return ErrorCode{
		httpStatus: code,
		name:       runes(http.StatusText(code)),
	}
}

// With provides a way to add additional context to the error message. For
// example a `BadRequest` can be given the context of `InvalidAuthorization`.
// Note: because of the limited nature of ErrorName it's possible that contexts
// can be dropped, this is by design so that abuse of error names are not
// overloaded for other purposes.
func (e ErrorCode) With(name string) ErrorCode {
	return ErrorCode{
		httpStatus: e.httpStatus,
		name:       suppliment(e.name, name),
	}
}

// Emoji allows you to add emoji code to the error codes for fun!
func (e ErrorCode) Emoji(emoji Emoji) ErrorCode {
	return e.With(emojiCodeMap[emoji])
}

// HTTPStatusCode returns the associated underlying raw HTTP status code.
func (e ErrorCode) HTTPStatusCode() int {
	return e.httpStatus
}

// Name returns the associated underlying error name.
func (e ErrorCode) Name() ErrorName {
	return e.name
}

func (e ErrorCode) String() string {
	var res []rune
	for _, v := range e.name {
		if v == 0 {
			continue
		}
		res = append(res, v)
	}
	return string(res)
}

// Is checks if the error name is associated with the error code.
func (e ErrorCode) Is(code ErrorName) bool {
	var (
		a = e.name
		b = code[:]
	)

	// Make sure we've got something useful
	if len(b) < 1 {
		return false
	}

	for k, v := range a {
		if v == b[0] && match(a[k:], b) {
			return true
		}
	}
	return false
}

func match(a, b []rune) bool {
	for k, v := range a {
		if v == 0 {
			continue
		}
		if v != b[k] {
			return false
		}
	}
	return true
}

func runes(s string) ErrorName {
	var res ErrorName
	for k, v := range s {
		if k >= len(res) {
			break
		}
		res[k] = v
	}
	return res
}

func suppliment(a ErrorName, b string) ErrorName {
	var (
		res    ErrorName
		name   = runes(" " + b)
		offset = 0
	)

	for k, v := range a {
		res[k] = v

		if v == 0 {
			offset = k
			break
		}
	}

	num := len(a)
	for k, v := range name {
		if offset+k < num {
			res[offset+k] = v
		}
	}
	return res
}

type Emoji string

const (
	ThumbsUp      Emoji = ":+1:"
	ThumbsDown    Emoji = ":-1:"
	Poop          Emoji = ":poop:"
	RotatingLight Emoji = ":rotating_light:"
	Warning       Emoji = ":warning:"
)

var emojiCodeMap = map[Emoji]string{
	ThumbsUp:      "\U0001f44d",
	ThumbsDown:    "\U0001f44e",
	Poop:          "\U0001f4a9",
	RotatingLight: "\U0001f6a8",
	Warning:       "\u26a0\ufe0f",
}
