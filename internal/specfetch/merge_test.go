package specfetch

import (
	"bytes"
	"encoding/json"
	"testing"
)

const shared = `"openapi":"3.1.0","info":{"title":"T","version":"1"},"servers":[{"url":"https://api.test"}]`

// jsonOf re-encodes v with sorted keys so two documents compare as strings.
func jsonOf(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func decoded(t *testing.T, s string) any {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader([]byte(s)))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		t.Fatal(err)
	}
	return v
}

func TestMergeCombinesPagesAndKeepsNumbersExact(t *testing.T) {
	a := `{` + shared + `,"paths":{"/a":{"get":{"operationId":"getA","parameters":[{"$ref":"#/components/parameters/Page"}]}}},` +
		`"components":{"parameters":{"Page":{"name":"page","in":"query"}}}}`
	b := `{` + shared + `,"paths":{"/a":{"post":{"operationId":"postA"}},` +
		`"/b":{"get":{"operationId":"getB","responses":{"200":{"description":"<b> & c","content":{"application/json":{"schema":{"maximum":9007199254740993}}}}}}}},` +
		`"components":{"parameters":{"Page":{"name":"page","in":"query"}}}}`
	want := `{` + shared + `,"paths":{` +
		`"/a":{"get":{"operationId":"getA","parameters":[{"$ref":"#/components/parameters/Page"}]},"post":{"operationId":"postA"}},` +
		`"/b":{"get":{"operationId":"getB","responses":{"200":{"description":"<b> & c","content":{"application/json":{"schema":{"maximum":9007199254740993}}}}}}}},` +
		`"components":{"parameters":{"Page":{"name":"page","in":"query"}}}}`

	got, err := Merge(map[string][]byte{"a.md": []byte(a), "b.md": []byte(b)})
	if err != nil {
		t.Fatal(err)
	}
	if g, w := jsonOf(t, got), jsonOf(t, decoded(t, want)); g != w {
		t.Errorf("got  %s\nwant %s", g, w)
	}
}

func TestMergeRejectsDisagreement(t *testing.T) {
	op := func(path, method string) string {
		return `{` + shared + `,"paths":{"` + path + `":{"` + method + `":{"operationId":"x"}}}}`
	}
	tests := []struct {
		name    string
		b       string
		wantErr string
	}{
		{name: "duplicate operation", b: op("/a", "get"), wantErr: "b.md: duplicate operation GET /a"},
		{
			name:    "conflicting component",
			b:       `{` + shared + `,"paths":{"/b":{"get":{}}},"components":{"parameters":{"Page":{"name":"p"}}}}`,
			wantErr: "b.md: components.parameters.Page differs from earlier pages",
		},
		{
			name:    "different info",
			b:       `{"openapi":"3.1.0","info":{"title":"Other","version":"1"},"paths":{"/b":{"get":{}}}}`,
			wantErr: "b.md: info differs from earlier pages",
		},
		{name: "no operation", b: `{` + shared + `,"paths":{}}`, wantErr: "b.md: page defines no operation"},
		{name: "not an object", b: `null`, wantErr: "b.md: not a JSON object"},
		{name: "invalid json", b: `{`, wantErr: "b.md: decode"},
	}
	a := `{` + shared + `,"paths":{"/a":{"get":{"operationId":"getA"}}},"components":{"parameters":{"Page":{"name":"page"}}}}`
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Merge(map[string][]byte{"a.md": []byte(a), "b.md": []byte(tt.b)})
			wantErr(t, err, tt.wantErr)
		})
	}
}
