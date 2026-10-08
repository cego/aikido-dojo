package cli

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
	"github.com/zalando/go-keyring"
)

// TestMain lets the scripts in testdata/script run aikido-dojo as a command.
func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){"aikido-dojo": scriptMain})
}

// scriptMain is aikido-dojo as a script runs it: the real command line, with
// the keychain in memory and every request sent to the fake Aikido that the
// script's setup started. Only the test binary has this seam.
func scriptMain() {
	keyring.MockInit()
	rt, err := scriptTransport()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(Run(context.Background(), os.Args[1:], Env{
		Stdin:      os.Stdin,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
		Getenv:     os.Getenv,
		CacheDir:   func() (string, error) { return filepath.Join(os.Getenv("WORK"), ".cache"), nil },
		Transport:  rt,
		ReadSecret: TerminalSecret(os.Stdin, os.Stderr),
		Confirm:    TerminalConfirm(os.Stdin, os.Stderr),
	}))
}

// scriptTransport dials the fake for every host and trusts only its certificate.
func scriptTransport() (http.RoundTripper, error) {
	ca, err := os.ReadFile(os.Getenv("AIKIDO_DOJO_TEST_CA")) //nolint:gosec // the path the test's own setup wrote
	if err != nil {
		return nil, fmt.Errorf("read the fake's certificate: %w", err)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(ca) {
		return nil, fmt.Errorf("no certificate in %s", os.Getenv("AIKIDO_DOJO_TEST_CA"))
	}
	addr := os.Getenv("AIKIDO_DOJO_TEST_ADDR")
	return &http.Transport{
		DialContext: func(ctx context.Context, network, _ string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, network, addr)
		},
		// httptest's certificate names example.com, whatever host the CLI meant.
		TLSClientConfig: &tls.Config{RootCAs: pool, ServerName: "example.com", MinVersion: tls.VersionTLS12},
	}, nil
}

func TestScripts(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/script",
		Setup: func(e *testscript.Env) error {
			srv := httptest.NewTLSServer(fakeAikido())
			e.Defer(srv.Close)
			ca := filepath.Join(e.WorkDir, ".fake-ca.pem")
			if err := os.WriteFile(ca, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: srv.Certificate().Raw}), 0o600); err != nil {
				return fmt.Errorf("write the fake's certificate: %w", err)
			}
			e.Setenv("AIKIDO_DOJO_TEST_ADDR", srv.Listener.Addr().String())
			e.Setenv("AIKIDO_DOJO_TEST_CA", ca)
			e.Setenv("GIT_CONFIG_NOSYSTEM", "1")
			return nil
		},
	})
}

// fakeAikido answers the calls the docs' examples make, with bodies shaped
// like the spec's examples.
func fakeAikido() http.Handler {
	enc := base64.RawURLEncoding
	token := enc.EncodeToString([]byte(`{"alg":"RS256"}`)) + "." +
		enc.EncodeToString([]byte(`{"scope":"basics:read issues:read issues:write repositories:read teams:write"}`)) + ".signature"
	repos := `[{"id":1,"name":"api","url":"https://gitlab.example.com/acme/api.git","provider":"gitlab_server"}]`
	mux := http.NewServeMux()
	reply := func(body string) http.HandlerFunc {
		return func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, body)
		}
	}
	mux.HandleFunc("POST /api/oauth/token", reply(`{"access_token":"`+token+`","expires_in":3600}`))
	mux.HandleFunc("GET /api/public/v1/workspace", reply(`{"id":1,"name":"Pied Piper","linked_provider":"gitlab"}`))
	mux.HandleFunc("GET /api/public/v1/repositories/code", func(w http.ResponseWriter, r *http.Request) {
		body := repos
		if r.URL.Query().Get("page") != "0" || (r.URL.Query().Has("filter_name") && r.URL.Query().Get("filter_name") != "api") {
			body = `[]`
		}
		reply(body)(w, r)
	})
	mux.HandleFunc("GET /api/public/v1/open-issue-groups", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Has-Next-Page", "false")
		reply(`[{"id":12,"title":"Prototype pollution in lodash","severity":"high","group_status":"new"}]`)(w, r)
	})
	mux.HandleFunc("POST /api/public/v1/issues/groups/12/notes", reply(`{"note_id":1}`))
	return mux
}

// Every command a doc shows must run in a script, so the docs can't drift.
func TestDocsCommandsAreScripted(t *testing.T) {
	scripted := map[string]bool{}
	scripts, err := filepath.Glob("testdata/script/*.txtar")
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range scripts {
		data, err := os.ReadFile(s) //nolint:gosec // the package's own scripts
		if err != nil {
			t.Fatal(err)
		}
		for _, line := range strings.Split(string(data), "\n") {
			if cmd, ok := strings.CutPrefix(strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "!")), "exec "); ok {
				scripted[cmd] = true
			}
		}
	}
	for doc, least := range map[string]int{"../../README.md": 6, "../../AGENTS.md": 3} {
		data, err := os.ReadFile(doc) //nolint:gosec // the repository's own docs
		if err != nil {
			t.Fatal(err)
		}
		cmds := docCommands(string(data))
		if len(cmds) < least {
			t.Errorf("%s shows %d commands, want at least %d", doc, len(cmds), least)
		}
		for _, cmd := range cmds {
			if !scripted[cmd] {
				t.Errorf("%s: %q runs in no script in testdata/script", doc, cmd)
			}
		}
	}
}

// docCommands lists the aikido-dojo command lines in a document's sh blocks.
func docCommands(md string) []string {
	var cmds []string
	inSh, open := false, false
	for _, line := range strings.Split(md, "\n") {
		if strings.HasPrefix(line, "```") {
			open = !open
			inSh = open && strings.TrimPrefix(line, "```") == "sh"
			continue
		}
		if inSh && strings.HasPrefix(line, "aikido-dojo ") {
			cmds = append(cmds, line)
		}
	}
	return cmds
}
