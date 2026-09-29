package specfetch

import (
	"testing"
	"time"
)

func TestExtractSpec(t *testing.T) {
	tests := []struct {
		name    string
		page    string
		want    string
		wantErr string
	}{
		{
			name: "takes the first json fence after the heading",
			page: "---\nupdatedAt: 2026-05-27T10:04:31.000Z\n---\n\n# List\n\nExample:\n```json\n{\"example\":true}\n```\n\n" +
				"# OpenAPI definition\n\n```json\n{\"openapi\":\"3.1.0\"}\n```\n\n```json\n{\"later\":true}\n```\n",
			want: `{"openapi":"3.1.0"}`,
		},
		{name: "no heading", page: "# List\n```json\n{}\n```\n", wantErr: `no "# OpenAPI definition" heading`},
		{name: "no fence", page: "x\n# OpenAPI definition\n\nnothing\n", wantErr: "no json block"},
		{name: "unterminated fence", page: "x\n# OpenAPI definition\n\n```json\n{\"a\":1}\n", wantErr: "unterminated json block"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractSpec([]byte(tt.page))
			if tt.wantErr != "" {
				wantErr(t, err, tt.wantErr)
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestUpdatedAt(t *testing.T) {
	tests := []struct {
		name    string
		page    string
		want    string
		wantErr string
	}{
		{
			name: "parses a fractional UTC timestamp",
			page: "---\nupdatedAt: 2026-05-27T10:04:31.000Z\nagentTools:\n  projectIndex: https://docs.test/llms.txt\n---\n# x\n",
			want: "2026-05-27T10:04:31Z",
		},
		{name: "no front matter", page: "# x\n", wantErr: "no front matter"},
		{name: "unterminated front matter", page: "---\nupdatedAt: 2026-05-27T10:04:31Z\n", wantErr: "unterminated front matter"},
		{name: "missing field", page: "---\nfoo: bar\n---\n", wantErr: "front matter has no updatedAt"},
		{name: "bad timestamp", page: "---\nupdatedAt: yesterday\n---\n", wantErr: "parse updatedAt"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpdatedAt([]byte(tt.page))
			if tt.wantErr != "" {
				wantErr(t, err, tt.wantErr)
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if s := got.UTC().Format(time.RFC3339); s != tt.want {
				t.Errorf("got %s, want %s", s, tt.want)
			}
		})
	}
}
