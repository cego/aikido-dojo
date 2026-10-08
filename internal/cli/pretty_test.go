package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestPrettyCopyMatchesIndent(t *testing.T) {
	for _, in := range []string{
		`{"b":1,"a":[1,2,{"c":null}],"e":{},"f":[]}`,
		`[]`,
		`{}`,
		`"x"`,
		`12345678901234567890`,
		`[true,false,null]`,
		`{"u":"æøå ☺","q":"a\"b\\c","t":"\t<&>"}`,
		`[[[]],[{}]]`,
	} {
		var want bytes.Buffer
		if err := json.Indent(&want, []byte(in), "", "  "); err != nil {
			t.Fatal(err)
		}
		want.WriteByte('\n')
		var got bytes.Buffer
		if err := prettyCopy(&got, strings.NewReader(in)); err != nil {
			t.Fatalf("%s: %v", in, err)
		}
		if got.String() != want.String() {
			t.Errorf("prettyCopy(%s) =\n%s\nwant\n%s", in, got.String(), want.String())
		}
	}
}

// Aikido escapes every / as \/, which a person reading a URL shouldn't see.
func TestPrettyCopyDropsEscapedSlashes(t *testing.T) {
	var got bytes.Buffer
	if err := prettyCopy(&got, strings.NewReader(`{"url":"https:\/\/x","u":"☺"}`)); err != nil {
		t.Fatal(err)
	}
	if want := "{\n  \"url\": \"https://x\",\n  \"u\": \"☺\"\n}\n"; got.String() != want {
		t.Errorf("got %q, want %q", got.String(), want)
	}
}

func TestPrettyCopyReportsBadJSON(t *testing.T) {
	if err := prettyCopy(&bytes.Buffer{}, strings.NewReader(`{"a":`)); err == nil {
		t.Error("want an error for truncated JSON")
	}
}

// Past the 32 KiB buffer the output is flushed in pieces; the result is the same.
func TestPrettyCopyFlushesALargeExport(t *testing.T) {
	in := "[" + strings.TrimSuffix(strings.Repeat(`{"id":12345,"name":"a repository name","tags":["x","y"]},`, 5000), ",") + "]"
	var want bytes.Buffer
	if err := json.Indent(&want, []byte(in), "", "  "); err != nil {
		t.Fatal(err)
	}
	want.WriteByte('\n')
	var got bytes.Buffer
	if err := prettyCopy(&got, strings.NewReader(in)); err != nil {
		t.Fatal(err)
	}
	if got.Len() < 64<<10 || got.String() != want.String() {
		t.Errorf("got %d bytes, want the %d bytes json.Indent gives", got.Len(), want.Len())
	}
}
