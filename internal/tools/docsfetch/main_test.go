package main

import (
	"errors"
	"io"
	"io/fs"
	"testing"
)

func TestRunOutsideTheRepositoryRootKeepsTheCause(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := run(nil, io.Discard); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("err = %v, want it to wrap fs.ErrNotExist", err)
	}
}
