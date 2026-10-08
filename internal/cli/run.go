package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/cego/aikido-dojo/internal/api"
	"github.com/cego/aikido-dojo/internal/auth"
	"github.com/cego/aikido-dojo/internal/catalog"
	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/config"
	"github.com/cego/aikido-dojo/internal/ops"
	"github.com/cego/aikido-dojo/internal/ratelimit"
	"github.com/cego/aikido-dojo/internal/validate"
)

func buildRoot(env Env) (*cobra.Command, error) {
	root, a := newRoot(env)
	if err := addGenerated(root, catalog.All, a.runOp); err != nil {
		return nil, err
	}
	root.AddCommand(a.versionCmd(), a.searchCmd(), a.schemaCmd(), a.apiCmd(), a.authCmd())
	if err := a.addRepoCurrent(root); err != nil {
		return nil, err
	}
	return root, nil
}

// runOp runs one generated command. Every input is checked before any
// credential is read, so a usage mistake costs no API call.
func (a *app) runOp(cmd *cobra.Command, c command, args []string) error {
	sc, err := catalog.Schemas(c.op.ID)
	if err != nil {
		return err
	}
	path, argWarnings, err := buildPath(c.op, args, sc.Args)
	if err != nil {
		return err
	}
	query, flagWarnings, err := buildQuery(cmd.Flags(), c.op, sc.Flags)
	if err != nil {
		return err
	}
	body, bodyWarnings, err := buildBody(cmd.Flags(), c, sc.Body, a.env.Stdin)
	if err != nil {
		return err
	}
	limit, err := limitOf(cmd.Flags(), c.op)
	if err != nil {
		return err
	}
	// Only writes have the flag. A dry run prints JSON of its own, and sends
	// nothing, so read-only mode allows it.
	dry, _ := cmd.Flags().GetBool("dry-run")
	if !dry {
		if err := a.out.check(c.op, sc); err != nil {
			return err
		}
	}
	for _, w := range slices.Concat(argWarnings, flagWarnings, bodyWarnings) {
		clierr.Warn(a.env.Stderr, w)
	}
	req := api.Request{Method: c.op.Method, Path: path, Query: query, Body: body, Scope: c.op.Scope}
	if dry {
		return a.dryRun(cmd.Context(), req)
	}
	if err := a.refuseWrite(c.op.Method, c.op.Command); err != nil {
		return err
	}
	if c.op.Destructive {
		yes, _ := cmd.Flags().GetBool("yes")
		if err := a.confirm(cmd.Context(), yes, strings.Join(append([]string{c.op.Command}, args...), " ")); err != nil {
			return err
		}
	}
	client, err := a.client()
	if err != nil {
		return err
	}
	if c.op.Paging != nil {
		return withHint(a.list(cmd.Context(), client, req, *c.op.Paging, limit), c.op)
	}
	resp, err := client.Do(cmd.Context(), req)
	if err != nil {
		return withHint(err, c.op)
	}
	defer resp.Body.Close()
	if c.op.Method != http.MethodGet {
		return a.writeResult(cmd.Context(), c.op.Command, resp)
	}
	if err := a.out.response(cmd.Context(), resp.Body, resp.Header.Get("Content-Type")); err != nil {
		return fmt.Errorf("%s: %w", c.op.Command, err)
	}
	return nil
}

// writeResult prints a write's response. It reads it whole first, as a
// write's answer is small, so that if printing fails the bytes still reach
// stdout and the error says the write took effect.
func (a *app) writeResult(ctx context.Context, what string, resp *http.Response) error {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return tookEffect(what, fmt.Errorf("read the response: %w", err), "read the result with a GET")
	}
	err = a.out.response(ctx, bytes.NewReader(data), resp.Header.Get("Content-Type"))
	if err == nil {
		return nil
	}
	next := "read the result with a GET"
	if len(data) > 0 {
		if _, werr := a.env.Stdout.Write(append(data, '\n')); werr == nil {
			next = "its response is on stdout as Aikido sent it"
		}
	}
	return tookEffect(what, err, next)
}

func (a *app) flags() config.Flags { return config.Flags{Config: a.config, Profile: a.profile} }

// configFile returns the config path and its contents; a missing file is empty.
func (a *app) configFile() (string, config.File, error) {
	path, err := config.Path(a.flags(), a.env.Getenv)
	if err != nil {
		return "", config.File{}, err
	}
	f, err := config.Load(path)
	return path, f, err
}

func (a *app) resolve() (config.Resolved, error) {
	_, f, err := a.configFile()
	if err != nil {
		return config.Resolved{}, err
	}
	return config.Resolve(f, a.flags(), a.env.Getenv)
}

func (a *app) client() (*api.Client, error) {
	r, err := a.resolve()
	if err != nil {
		return nil, err
	}
	return a.clientFor(r)
}

func (a *app) clientFor(r config.Resolved) (*api.Client, error) {
	hc := newHTTPClient(a.env, r.Host, a.debug)
	src, err := auth.NewSource(hc, r)
	if err != nil {
		return nil, err
	}
	cache, err := a.env.CacheDir()
	if err != nil {
		return nil, fmt.Errorf("find the cache directory for the rate limiter: %w", err)
	}
	return api.New(hc, r.Host, src, ratelimit.New(filepath.Join(cache, "aikido-dojo"), r.ClientID)), nil
}

// list reads every page, or up to limit items, and prints them as one JSON
// array. It prints only after the last page, so a failure never leaves a
// truncated array on stdout. With --ndjson it prints each item as it arrives
// instead, and a failure leaves the lines already printed.
func (a *app) list(ctx context.Context, c *api.Client, req api.Request, p api.Paging, limit int64) error {
	if a.out.ndjson {
		var n int64
		for item, err := range c.Items(ctx, req, p) {
			if err != nil {
				return err
			}
			if err := a.out.item(ctx, item); err != nil {
				return err
			}
			if n++; n == limit {
				break
			}
		}
		return nil
	}
	items := []json.RawMessage{}
	for item, err := range c.Items(ctx, req, p) {
		if err != nil {
			return err
		}
		items = append(items, item)
		if int64(len(items)) == limit {
			break
		}
	}
	return a.out.value(ctx, items)
}

// withHint swaps in the operation's own hint for a 400, where the overlay
// knows the likely cause.
func withHint(err error, op ops.Op) error {
	var e *clierr.Error
	if op.BadRequestHint != "" && errors.As(err, &e) && e.HTTPStatus == http.StatusBadRequest {
		e.Hint = op.BadRequestHint
	}
	return err
}

func buildPath(op ops.Op, args []string, schema map[string]any) (string, []string, error) {
	values := map[string]any{}
	texts := make([]string, len(args))
	for i, p := range op.Args {
		if args[i] == "." || args[i] == ".." {
			return "", nil, invalid(fmt.Sprintf("<%s>: %q is not an ID", p.Name, args[i]), op)
		}
		v, text, err := scalar(p.Kind, args[i])
		if err != nil {
			return "", nil, invalid("<"+p.Name+">: "+err.Error(), op)
		}
		values[p.Name], texts[i] = v, text
	}
	warnings, err := validated(schema, values, op, argName)
	if err != nil {
		return "", nil, err
	}
	path := op.Path
	for i, p := range op.Args {
		// Escaped, so a value can't reach another endpoint through / or ?.
		path = strings.Replace(path, "{"+p.Name+"}", url.PathEscape(texts[i]), 1)
	}
	return path, warnings, nil
}

func buildQuery(fs *pflag.FlagSet, op ops.Op, schema map[string]any) (url.Values, []string, error) {
	query := url.Values{}
	values := map[string]any{}
	for _, p := range op.Flags {
		if !fs.Changed(ops.FlagName(p.Name)) {
			continue
		}
		v, text, err := flagValue(fs, p)
		if err != nil {
			return nil, nil, invalid("--"+ops.FlagName(p.Name)+": "+err.Error(), op)
		}
		values[p.Name] = v
		query.Set(p.Name, text)
	}
	warnings, err := validated(schema, values, op, flagName)
	if err != nil {
		return nil, nil, err
	}
	return query, warnings, nil
}

// buildBody merges --body or --body-file with the field flags, which win, and
// validates the result. An optional body nobody set is left out.
func buildBody(fs *pflag.FlagSet, c command, schema map[string]any, stdin io.Reader) ([]byte, []string, error) {
	b := c.op.Body
	if b == nil {
		return nil, nil, nil
	}
	body, given, err := bodyInput(fs, stdin, c.op)
	if err != nil {
		return nil, nil, err
	}
	// An inline --body is argv, which ps and shell history keep.
	if obj, ok := body.(map[string]any); ok && fs.Changed("body") {
		for _, name := range b.Secret {
			if _, set := obj[name]; set {
				return nil, nil, invalid("--body: "+name+" is a credential; pass the body with --body-file <path> or --body-file - (stdin)", c.op)
			}
		}
	}
	fields := map[string]any{}
	for _, p := range c.bodyFlags {
		if !fs.Changed(ops.FlagName(p.Name)) {
			continue
		}
		v, _, err := flagValue(fs, p)
		if err != nil {
			return nil, nil, invalid("--"+ops.FlagName(p.Name)+": "+err.Error(), c.op)
		}
		fields[p.Name] = v
	}
	if len(fields) > 0 {
		if !given {
			body, given = map[string]any{}, true
		}
		obj, ok := body.(map[string]any)
		if !ok {
			return nil, nil, invalid("--body must be a JSON object when field flags are set", c.op)
		}
		maps.Copy(obj, fields)
	}
	if !given {
		if !b.Required {
			return nil, nil, nil
		}
		if !b.Object {
			return nil, nil, invalid(c.op.Command+" needs a body: pass --body or --body-file", c.op)
		}
		body = map[string]any{}
	}
	warnings, err := validated(schema, body, c.op, bodyName)
	if err != nil {
		return nil, nil, err
	}
	var out bytes.Buffer
	enc := json.NewEncoder(&out)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(body); err != nil {
		return nil, nil, fmt.Errorf("encode the body: %w", err)
	}
	return bytes.TrimSuffix(out.Bytes(), []byte("\n")), warnings, nil
}

// bodyInput decodes --body, or --body-file (a path, or - for stdin). given is
// false when neither is set.
func bodyInput(fs *pflag.FlagSet, stdin io.Reader, op ops.Op) (v any, given bool, err error) {
	inline, file := fs.Changed("body"), fs.Changed("body-file")
	var raw []byte
	switch {
	case inline && file:
		return nil, false, invalid("use --body or --body-file, not both", op)
	case inline:
		s, err := fs.GetString("body")
		if err != nil {
			return nil, false, fmt.Errorf("read --body: %w", err)
		}
		raw = []byte(s)
	case file:
		path, err := fs.GetString("body-file")
		if err != nil {
			return nil, false, fmt.Errorf("read --body-file: %w", err)
		}
		if raw, err = readBodyFile(path, stdin); err != nil {
			return nil, false, invalid("--body-file: "+err.Error(), op)
		}
	default:
		return nil, false, nil
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&v); err != nil {
		return nil, false, invalid("body: invalid JSON: "+err.Error(), op)
	}
	if _, err := dec.Token(); !errors.Is(err, io.EOF) {
		return nil, false, invalid("body: unexpected data after the JSON value", op)
	}
	return v, true, nil
}

func readBodyFile(path string, stdin io.Reader) ([]byte, error) {
	if path == "-" {
		data, err := io.ReadAll(stdin)
		if err != nil {
			return nil, fmt.Errorf("read stdin: %w", err)
		}
		return data, nil
	}
	data, err := os.ReadFile(path) //nolint:gosec // the user names the file to send
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	return data, nil
}

func limitOf(fs *pflag.FlagSet, op ops.Op) (int64, error) {
	if op.Paging == nil {
		return 0, nil
	}
	n, err := fs.GetInt64("limit")
	if err != nil {
		return 0, fmt.Errorf("read --limit: %w", err)
	}
	if n < 0 {
		return 0, invalid("--limit: want 0 (every item) or more", op)
	}
	return n, nil
}

// flagValue reads a set flag as the JSON value the schema checks and the text
// the query sends. Lists go comma-separated, as the spec's descriptions of
// every list parameter ask.
func flagValue(fs *pflag.FlagSet, p ops.Param) (any, string, error) {
	name := ops.FlagName(p.Name)
	switch p.Kind {
	case ops.Boolean:
		v, err := fs.GetBool(name)
		if err != nil {
			return nil, "", fmt.Errorf("read the flag: %w", err)
		}
		return v, strconv.FormatBool(v), nil
	case ops.Integer:
		v, err := fs.GetInt64(name)
		if err != nil {
			return nil, "", fmt.Errorf("read the flag: %w", err)
		}
		text := strconv.FormatInt(v, 10)
		return json.Number(text), text, nil
	case ops.StringList, ops.IntegerList:
		items, err := fs.GetStringSlice(name)
		if err != nil {
			return nil, "", fmt.Errorf("read the flag: %w", err)
		}
		kind := ops.String
		if p.Kind == ops.IntegerList {
			kind = ops.Integer
		}
		values := make([]any, len(items))
		for i, item := range items {
			v, text, err := scalar(kind, item)
			if err != nil {
				return nil, "", err
			}
			values[i], items[i] = v, text
		}
		return values, strings.Join(items, ","), nil
	default:
		v, err := fs.GetString(name)
		if err != nil {
			return nil, "", fmt.Errorf("read the flag: %w", err)
		}
		return v, v, nil
	}
}

// scalar turns command-line text into the JSON value the schema checks and
// the text to send.
func scalar(kind ops.Kind, text string) (any, string, error) {
	if kind != ops.Integer {
		return text, text, nil
	}
	n, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return nil, "", fmt.Errorf("%q is not an integer", text)
	}
	s := strconv.FormatInt(n, 10)
	return json.Number(s), s, nil
}

// validated checks v against schema and names each issue the way the user
// typed it: <arg>, --flag or body.field.
func validated(schema map[string]any, v any, op ops.Op, name func([]string) string) ([]string, error) {
	res := validate.Check(schema, v)
	warnings := make([]string, 0, len(res.Warnings))
	for _, w := range res.Warnings {
		warnings = append(warnings, name(w.Path)+": "+w.Msg)
	}
	if len(res.Errors) == 0 {
		return warnings, nil
	}
	msgs := make([]string, 0, len(res.Errors))
	for _, e := range res.Errors {
		msgs = append(msgs, name(e.Path)+": "+e.Msg)
	}
	return nil, invalid(strings.Join(msgs, "; "), op)
}

func argName(path []string) string {
	if len(path) == 0 {
		return "arguments"
	}
	return validate.Join("<"+path[0]+">", path[1:])
}

func flagName(path []string) string {
	if len(path) == 0 {
		return "flags"
	}
	return validate.Join("--"+ops.FlagName(path[0]), path[1:])
}

func bodyName(path []string) string { return validate.Join("body", path) }

func invalid(msg string, op ops.Op) error {
	return &clierr.Error{Code: "invalid_input", Message: msg, Hint: "see aikido-dojo " + op.Command + " --help", Exit: clierr.ExitUsage}
}
