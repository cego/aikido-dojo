package auth

import (
	"encoding/base64"
	"slices"
	"testing"
)

func jwt(claims string) string {
	enc := base64.RawURLEncoding
	return enc.EncodeToString([]byte(`{"alg":"RS256","typ":"JWT"}`)) + "." + enc.EncodeToString([]byte(claims)) + ".signature"
}

func TestScopes(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  []string
		ok    bool
	}{
		{"Aikido's space-separated claim", jwt(`{"client_id":"x","scope":"basics:read issues:read"}`), []string{"basics:read", "issues:read"}, true},
		{"an empty claim grants nothing", jwt(`{"scope":""}`), []string{}, true},
		{"no claim", jwt(`{"client_id":"x"}`), nil, false},
		{"a claim that isn't a string", jwt(`{"scope":["a"]}`), nil, false},
		{"an opaque token", "tok", nil, false},
		{"a payload that isn't base64", "a.!!!.c", nil, false},
	}
	for _, tt := range tests {
		got, ok := Scopes(tt.token)
		if ok != tt.ok || !slices.Equal(got, tt.want) {
			t.Errorf("%s: Scopes = %q, %v; want %q, %v", tt.name, got, ok, tt.want, tt.ok)
		}
	}
}
