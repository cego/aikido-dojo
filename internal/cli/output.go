package cli

import (
	"encoding/json"
	"fmt"
)

// printJSON writes v to stdout as one line of JSON, the form all of
// aikido-dojo's own output takes.
func (a *app) printJSON(v any) error {
	enc := json.NewEncoder(a.env.Stdout)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("write the output: %w", err)
	}
	return nil
}
