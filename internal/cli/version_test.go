package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"runtime/debug"
	"testing"
	"time"
)

func TestVersionOf(t *testing.T) {
	bi := &debug.BuildInfo{GoVersion: "go1.27.1", Main: debug.Module{Version: "v0.1.0"},
		Settings: []debug.BuildSetting{{Key: "vcs.revision", Value: "7cae439"}, {Key: "vcs.modified", Value: "false"}}}
	if got := versionOf(bi, true); got.Version != "v0.1.0" || got.Commit != "7cae439" || got.Go != "go1.27.1" {
		t.Errorf("versionOf = %+v", got)
	}
	if got := versionOf(nil, false); got.Version != "(unknown)" || got.Commit != "" {
		t.Errorf("without build info = %+v", got)
	}
}

func TestVersionNamesTheVendoredSpec(t *testing.T) {
	spec, err := os.ReadFile("../../spec/openapi.json")
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile("../../spec/snapshot.json")
	if err != nil {
		t.Fatal(err)
	}
	var snapshot struct {
		UpdatedAt time.Time `json:"updated_at"`
	}
	if err := json.Unmarshal(raw, &snapshot); err != nil {
		t.Fatal(err)
	}
	env, _ := testEnv(t)
	var got versionInfo
	if err := json.Unmarshal([]byte(mustRun(t, env, "version")), &got); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(spec)
	if got.Spec.SHA256 != hex.EncodeToString(sum[:]) || got.Spec.UpdatedAt != snapshot.UpdatedAt.UTC().Format(time.RFC3339) {
		t.Errorf("spec = %+v, want the vendored spec's hash and snapshot date", got.Spec)
	}
}
