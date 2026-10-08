package cli

import (
	"net/http"

	"github.com/cego/aikido-dojo/internal/clierr"
)

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
