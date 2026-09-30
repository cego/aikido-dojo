package clierr

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

func TestReport(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantJSON string
		wantExit int
	}{
		{
			name:     "a plain error is unexpected",
			err:      errors.New("boom"),
			wantJSON: `{"error":{"code":"unexpected","message":"boom"}}`,
			wantExit: ExitUnexpected,
		},
		{
			name: "keeps the caller's context and the hint",
			err: fmt.Errorf("list repos: %w", &Error{
				Code: "missing_scope", Message: "GET /repositories/code: forbidden", HTTPStatus: 403,
				Hint: "grant repositories:read", Exit: ExitForbidden,
			}),
			wantJSON: `{"error":{"code":"missing_scope","message":"list repos: GET /repositories/code: forbidden","http_status":403,"hint":"grant repositories:read"}}`,
			wantExit: ExitForbidden,
		},
		{
			name:     "includes the wrapped cause",
			err:      &Error{Code: "bad_config", Message: "parse config", Err: errors.New("unexpected EOF"), Exit: ExitUsage},
			wantJSON: `{"error":{"code":"bad_config","message":"parse config: unexpected EOF"}}`,
			wantExit: ExitUsage,
		},
		{
			name:     "an unset exit code is unexpected, never success",
			err:      &Error{Code: "x", Message: "m"},
			wantJSON: `{"error":{"code":"x","message":"m"}}`,
			wantExit: ExitUnexpected,
		},
		{
			name:     "does not escape html characters",
			err:      errors.New("a <b> & c"),
			wantJSON: `{"error":{"code":"unexpected","message":"a <b> & c"}}`,
			wantExit: ExitUnexpected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if got := Report(&buf, tt.err); got != tt.wantExit {
				t.Errorf("exit = %d, want %d", got, tt.wantExit)
			}
			if got := buf.String(); got != tt.wantJSON+"\n" {
				t.Errorf("stderr = %s, want %s", got, tt.wantJSON)
			}
		})
	}
}
