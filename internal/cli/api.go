package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cego/aikido-dojo/internal/api"
	"github.com/cego/aikido-dojo/internal/catalog"
	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/ops"
)

// apiMethods are the methods the spec uses; anything else is a typo.
var apiMethods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}

// apiFlags are api's own flags.
type apiFlags struct {
	fields   []string
	input    string
	hasInput bool
	dryRun   bool
	yes      bool
}

func (a *app) apiCmd() *cobra.Command {
	var f apiFlags
	cmd := &cobra.Command{
		Use:   "api <method> <path>",
		Short: "Call any public API path, with the same credentials, pacing, retries and errors",
		Long: "api calls a path under /api/public/v1 and prints the response as it arrives. It covers what the " +
			"generated commands get wrong, and shares their credentials, rate limiting, retries and errors.\n\n" +
			"-f key=value adds a query parameter to a GET, or to any call with --input; otherwise it adds a string " +
			"field to a JSON body. --input sends a JSON file, or - for stdin, as the body. When the method and path " +
			"match a generated command, a 403 names the scope that command needs.",
		Example: "  aikido-dojo api GET /repositories/code -f per_page=5\n" +
			`  aikido-dojo api POST /issues/groups/12/notes -f note="false positive: a test fixture"`,
		Args: exactArgs([]ops.Param{{Name: "method"}, {Name: "path"}}),
		RunE: func(cmd *cobra.Command, args []string) error {
			f.hasInput = cmd.Flags().Changed("input")
			return a.callAPI(cmd.Context(), args[0], args[1], f)
		},
	}
	// StringArray, not StringSlice: a value may hold commas.
	cmd.Flags().StringArrayVarP(&f.fields, "field", "f", nil, "key=value: a query parameter for GET or with --input, else a string field of the JSON body")
	cmd.Flags().StringVar(&f.input, "input", "", "a file holding the JSON request body, or - for stdin")
	cmd.Flags().BoolVar(&f.dryRun, "dry-run", false, dryRunUsage)
	cmd.Flags().BoolVar(&f.yes, "yes", false, yesUsage+": a DELETE, or a path a destructive command calls")
	return cmd
}

func (a *app) callAPI(ctx context.Context, method, target string, f apiFlags) error {
	method = strings.ToUpper(method)
	if !slices.Contains(apiMethods, method) {
		return apiUsage(fmt.Sprintf("method %q: use GET, POST, PUT or DELETE", method))
	}
	path, query, err := apiTarget(target)
	if err != nil {
		return err
	}
	var body []byte
	if f.hasInput {
		if body, err = jsonInput(ctx, f.input, a.env.Stdin); err != nil {
			return err
		}
	}
	op, matched := catalog.Match(method, path)
	toQuery := method == http.MethodGet || f.hasInput
	obj := map[string]string{}
	for _, field := range f.fields {
		k, v, ok := strings.Cut(field, "=")
		if !ok || k == "" {
			return apiUsage(fmt.Sprintf("-f %q: want key=value", field))
		}
		// argv and shell history keep -f values, and a query also reaches logs, so
		// credentials come in through --input, as with the generated command.
		if slices.Contains(credentialFields(), k) {
			return apiUsage(fmt.Sprintf("-f %s: a credential; pass the body with --input <file> or --input - (stdin)", k))
		}
		if toQuery {
			query.Add(k, v)
			continue
		}
		if _, dup := obj[k]; dup {
			return apiUsage(fmt.Sprintf("-f %s is given twice", k))
		}
		obj[k] = v
	}
	if len(obj) > 0 {
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(obj); err != nil {
			return fmt.Errorf("encode the body: %w", err)
		}
		body = bytes.TrimSuffix(buf.Bytes(), []byte("\n"))
	}
	req := api.Request{Method: method, Path: path, Query: query, Body: body, Scope: op.Scope}
	if f.dryRun {
		return a.dryRun(ctx, req)
	}
	if matched {
		sc, err := catalog.Schemas(op.ID)
		if err != nil {
			return err
		}
		if err := a.out.check(op, sc); err != nil {
			return err
		}
		// The raw page of an envelope list is an object; the generated command unwraps it.
		if a.out.ndjson && op.Paging != nil && op.Paging.Items != "" {
			return &clierr.Error{Code: "invalid_input", Message: "--ndjson needs a list, but " + path + " returns its items inside an object",
				Hint: "use aikido-dojo " + op.Command + " --ndjson, which pages and unwraps them", Exit: clierr.ExitUsage}
		}
	}
	if err := a.refuseWrite(method, "api "+method+" "+path); err != nil {
		return err
	}
	if method == http.MethodDelete || op.Destructive {
		if err := a.confirm(ctx, f.yes, apiCommandLine(method, target, f)); err != nil {
			return err
		}
	}
	client, err := a.client()
	if err != nil {
		return err
	}
	resp, err := client.Do(ctx, req)
	if err != nil {
		return withHint(err, op)
	}
	defer resp.Body.Close()
	if method != http.MethodGet {
		return a.writeResult(ctx, "api "+method+" "+path, resp)
	}
	if err := a.out.response(ctx, resp.Body, resp.Header.Get("Content-Type")); err != nil {
		return fmt.Errorf("api %s %s: %w", method, path, err)
	}
	return nil
}

// apiTarget splits /repositories/code?per_page=5 into the path under
// /api/public/v1 and its query. It refuses what would leave that prefix: a
// URL of its own, or a . or .. segment, escaped or not, including one that an
// escaped / or a \ splits off inside a segment, as a server may do.
func apiTarget(s string) (string, url.Values, error) {
	path, rawQuery, _ := strings.Cut(s, "?")
	if strings.Contains(path, "://") {
		return "", nil, apiUsage(fmt.Sprintf("%q: give the path under /api/public/v1, such as /repositories/code", s))
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	// Checked here, so a path no URL can hold, such as one with a control
	// character, costs no token request.
	if _, err := url.Parse("https://host" + path); err != nil {
		var ue *url.Error
		if errors.As(err, &ue) {
			err = ue.Err
		}
		return "", nil, apiUsage(fmt.Sprintf("path %q: %v", path, err))
	}
	if strings.Contains(path, "#") {
		return "", nil, apiUsage(fmt.Sprintf("path %q: a # starts a fragment, which is never sent; write it as %%23", path))
	}
	for _, seg := range strings.Split(path, "/") {
		dec, err := url.PathUnescape(seg)
		if err != nil {
			return "", nil, apiUsage(fmt.Sprintf("path %q: %v", path, err))
		}
		for _, part := range strings.FieldsFunc(dec, func(r rune) bool { return r == '/' || r == '\\' }) {
			if part == "." || part == ".." {
				return "", nil, apiUsage(fmt.Sprintf("path %q: . and .. segments could leave /api/public/v1", path))
			}
		}
	}
	query, err := url.ParseQuery(rawQuery)
	if err != nil {
		return "", nil, apiUsage(fmt.Sprintf("query %q: %v", rawQuery, err))
	}
	return path, query, nil
}

// jsonInput reads --input and checks it holds one JSON value, so a typo
// costs no call. The bytes are sent as they are.
func jsonInput(ctx context.Context, path string, stdin io.Reader) ([]byte, error) {
	data, err := readBodyFile(ctx, path, stdin)
	if err != nil {
		return nil, apiUsage("--input: " + err.Error())
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, apiUsage("--input: invalid JSON: " + err.Error())
	}
	if _, err := dec.Token(); !errors.Is(err, io.EOF) {
		return nil, apiUsage("--input: unexpected data after the JSON value")
	}
	return data, nil
}

// apiCommandLine is the call as typed, for the destructive question: -f and
// --input may name its target.
func apiCommandLine(method, target string, f apiFlags) string {
	parts := []string{"api", method, target}
	for _, field := range f.fields {
		parts = append(parts, "-f", quoted(field))
	}
	if f.hasInput {
		parts = append(parts, "--input", quoted(f.input))
	}
	return strings.Join(parts, " ")
}

func apiUsage(msg string) error {
	return &clierr.Error{Code: "invalid_input", Message: msg, Hint: "see aikido-dojo api --help", Exit: clierr.ExitUsage}
}
