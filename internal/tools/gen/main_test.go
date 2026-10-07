package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunWritesTheCatalog(t *testing.T) {
	dir := t.TempDir()
	if err := run([]string{"-spec", "../../../spec/openapi.json", "-snapshot", "../../../spec/snapshot.json", "-dir", dir}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"zz_ops.go", "schemas.json", "search.json"} {
		if info, err := os.Stat(filepath.Join(dir, name)); err != nil || info.Size() == 0 {
			t.Errorf("%s: %v, want a non-empty file", name, err)
		}
	}
}

func TestRunFailures(t *testing.T) {
	spec, snapshot := "../../../spec/openapi.json", "../../../spec/snapshot.json"
	tests := []struct {
		args []string
		want string
	}{
		{[]string{"-spec", "x.json"}, "usage: gen"},
		{[]string{"-spec", spec, "-dir", t.TempDir()}, "usage: gen"},
		{[]string{"-spec", filepath.Join(t.TempDir(), "missing.json"), "-snapshot", snapshot, "-dir", t.TempDir()}, "read the spec"},
		{[]string{"-spec", spec, "-snapshot", filepath.Join(t.TempDir(), "missing.json"), "-dir", t.TempDir()}, "read the snapshot"},
		{[]string{"-spec", spec, "-snapshot", snapshot, "-dir", filepath.Join(t.TempDir(), "absent")}, "write zz_ops.go"},
	}
	for _, tt := range tests {
		if err := run(tt.args); err == nil || !strings.Contains(err.Error(), tt.want) {
			t.Errorf("run(%q) = %v, want an error containing %q", tt.args, err, tt.want)
		}
	}
}
