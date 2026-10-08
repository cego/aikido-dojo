package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cego/aikido-dojo/internal/api"
	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/ops"
)

const (
	dryRunUsage = "print the request this would send, with secrets redacted, and send nothing"
	yesUsage    = "make the destructive call without asking"
)

// confirm guards a destructive call: on a terminal a person answers the
// question; off one, as for an agent, --yes must be given.
func (a *app) confirm(ctx context.Context, yes bool, what string) error {
	if yes {
		return nil
	}
	ok, err := a.env.Confirm(ctx, what+" is destructive. Continue? [y/N] ")
	var e *clierr.Error
	if errors.As(err, &e) && e.Code == "no_terminal" {
		return &clierr.Error{Code: "confirmation_required", Message: what + " is destructive and needs confirmation",
			Hint: "pass --yes to confirm it", Exit: clierr.ExitRefused}
	}
	if err != nil {
		return err
	}
	if !ok {
		return &clierr.Error{Code: "not_confirmed", Message: what + " was not confirmed", Exit: clierr.ExitRefused}
	}
	return nil
}

type dryRunRequest struct {
	DryRun  bool              `json:"dry_run"`
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    json.RawMessage   `json:"body,omitempty"`
}

// dryRun prints req as a call would send it, with the token and the body's
// credential fields redacted. It resolves the profile to name the host, but
// reads no secret and sends nothing.
func (a *app) dryRun(ctx context.Context, req api.Request, secret []string) error {
	r, err := a.resolve()
	if err != nil {
		return err
	}
	out := dryRunRequest{DryRun: true, Method: req.Method, URL: api.URL(r.Host, req),
		Headers: map[string]string{"Authorization": "Bearer [redacted]"}}
	if req.Body != nil {
		out.Headers["Content-Type"] = "application/json"
		if out.Body, err = redacted(req.Body, secret); err != nil {
			return err
		}
	}
	return a.out.value(ctx, out)
}

// redacted replaces the credential fields of a JSON object body. A body that
// isn't an object has no named fields to redact.
func redacted(body []byte, secret []string) (json.RawMessage, error) {
	if len(secret) == 0 || !bytes.HasPrefix(bytes.TrimSpace(body), []byte("{")) {
		return body, nil
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(body, &obj); err != nil {
		return nil, fmt.Errorf("decode the body to redact it: %w", err)
	}
	for _, k := range secret {
		if _, ok := obj[k]; ok {
			obj[k] = json.RawMessage(`"[redacted]"`)
		}
	}
	return compact(obj)
}

func secretFields(op ops.Op) []string {
	if op.Body == nil {
		return nil
	}
	return op.Body.Secret
}

// refuseWrite enforces read-only mode: every call that isn't a GET is
// refused before any credential is read. The hint names no way out, since
// the mode is how a person hands aikido-dojo to an agent.
func (a *app) refuseWrite(method, what string) error {
	if !a.readOnly || method == http.MethodGet {
		return nil
	}
	return &clierr.Error{Code: "read_only", Message: "read-only mode refuses " + what,
		Hint: "read-only mode (--read-only or AIKIDO_DOJO_READ_ONLY) allows only GET calls; ask a person to run this write", Exit: clierr.ExitRefused}
}
