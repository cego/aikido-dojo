package auth

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

// Scopes reads the scope claim of an access token. Aikido's tokens are JWTs
// whose claim lists the API client's scopes, though the token response names
// none (live, 2026-09-29). The signature isn't checked: the claim only
// informs auth status, and Aikido enforces the scopes itself.
func Scopes(token string) ([]string, bool) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, false
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, false
	}
	var claims struct {
		Scope *string `json:"scope"`
	}
	if json.Unmarshal(payload, &claims) != nil || claims.Scope == nil {
		return nil, false
	}
	return strings.Fields(*claims.Scope), true
}
