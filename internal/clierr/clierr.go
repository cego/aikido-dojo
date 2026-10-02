// Package clierr is the CLI's error model: a failure is reported as one JSON
// object on stderr and exits with a documented code.
package clierr

import (
	"cmp"
	"encoding/json"
	"errors"
	"io"
)

// Exit codes are part of the CLI's contract and are documented in the README.
const (
	ExitOK          = 0
	ExitUnexpected  = 1
	ExitUsage       = 2
	ExitAuth        = 3
	ExitForbidden   = 4
	ExitNotFound    = 5
	ExitRateLimited = 6
	ExitRefused     = 7
)

// Error is a failure the CLI can explain. Hint tells the reader what to do next.
type Error struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"http_status,omitempty"`
	Hint       string `json:"hint,omitempty"`
	Exit       int    `json:"-"`
	Err        error  `json:"-"`
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *Error) Unwrap() error { return e.Err }

// Report writes err to w as {"error":{...}} on one line and returns the exit
// code. The message is the whole error chain, so context added by callers
// survives. Errors that aren't *Error are reported as unexpected.
func Report(w io.Writer, err error) int {
	out := Error{Code: "unexpected", Message: err.Error(), Exit: ExitUnexpected}
	var e *Error
	if errors.As(err, &e) {
		out.Code = cmp.Or(e.Code, out.Code)
		out.HTTPStatus, out.Hint = e.HTTPStatus, e.Hint
		if e.Exit != ExitOK {
			out.Exit = e.Exit
		}
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	// stderr is the only place to report to, so a failed write there has no remedy.
	_ = enc.Encode(struct {
		Error Error `json:"error"`
	}{out})
	return out.Exit
}

// Warn writes a warning to w as {"warning":{"message":...}} on one line, so an
// agent reading stderr can tell it from the error object.
func Warn(w io.Writer, msg string) {
	type warning struct {
		Message string `json:"message"`
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	// stderr is the only place to report to, so a failed write there has no remedy.
	_ = enc.Encode(struct {
		Warning warning `json:"warning"`
	}{warning{msg}})
}
