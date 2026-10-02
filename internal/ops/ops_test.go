package ops

import "testing"

func TestFlagName(t *testing.T) {
	for in, want := range map[string]string{"filter_code_repo_id": "filter-code-repo-id", "search": "search", "per_page": "per-page"} {
		if got := FlagName(in); got != want {
			t.Errorf("FlagName(%q) = %q, want %q", in, got, want)
		}
	}
}
