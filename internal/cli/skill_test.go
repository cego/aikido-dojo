package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/config"
)

func installed(t *testing.T, out string) []installedSkill {
	t.Helper()
	var got []installedSkill
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("decode %q: %v", out, err)
	}
	return got
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path) //nolint:gosec // a path the test chose
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func mkdir(t *testing.T, path ...string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(path...), 0o700); err != nil {
		t.Fatal(err)
	}
}

// Agents that share ~/.agents/skills get one copy; agents not installed get none.
func TestSkillInstallFindsTheInstalledAgents(t *testing.T) {
	env, _ := testEnv(t)
	home, _ := env.Home()
	mkdir(t, home, ".claude")
	mkdir(t, home, ".codex")
	mkdir(t, home, ".copilot")
	mkdir(t, home, ".config", "opencode")
	got := installed(t, mustRun(t, env, "skill", "install"))
	want := []installedSkill{
		{Path: filepath.Join(home, ".claude", "skills", "aikido-dojo", "SKILL.md"), Agents: []string{"claude"}},
		{Path: filepath.Join(home, ".agents", "skills", "aikido-dojo", "SKILL.md"), Agents: []string{"codex", "copilot"}},
		{Path: filepath.Join(home, ".config", "opencode", "skills", "aikido-dojo", "SKILL.md"), Agents: []string{"opencode"}},
	}
	if !slices.EqualFunc(got, want, func(a, b installedSkill) bool { return a.Path == b.Path && slices.Equal(a.Agents, b.Agents) }) {
		t.Fatalf("installed %+v, want %+v", got, want)
	}
	for _, s := range got {
		if body := readFile(t, s.Path); body != string(env.Primer) {
			t.Errorf("%s = %q, want the primer", s.Path, body)
		}
		info, err := os.Stat(s.Path)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm() != 0o644 {
			t.Errorf("%s mode = %v, want 0644", s.Path, info.Mode().Perm())
		}
	}
	if _, err := os.Stat(filepath.Join(home, ".cursor")); err == nil {
		t.Error("installing created ~/.cursor for an agent that isn't installed")
	}
}

func TestSkillInstallNeedsAnAgent(t *testing.T) {
	env, _ := testEnv(t)
	_, stderr, code := run(t, env, "skill", "install")
	if e := errorOf(t, stderr); code != clierr.ExitUsage || e.Code != "no_agent" {
		t.Errorf("exit %d, error %+v", code, e)
	}
}

// --agent installs even when the agent's config directory doesn't exist yet.
func TestSkillInstallForNamedAgents(t *testing.T) {
	env, _ := testEnv(t)
	home, _ := env.Home()
	mkdir(t, home, ".claude")
	got := installed(t, mustRun(t, env, "skill", "install", "--agent", "codex", "--agent", "Gemini"))
	if len(got) != 1 || got[0].Path != filepath.Join(home, ".agents", "skills", "aikido-dojo", "SKILL.md") ||
		!slices.Equal(got[0].Agents, []string{"codex", "gemini"}) {
		t.Errorf("installed %+v, want one copy in ~/.agents/skills", got)
	}
	_, stderr, code := run(t, env, "skill", "install", "--agent", "vim")
	if e := errorOf(t, stderr); code != clierr.ExitUsage || e.Code != "invalid_input" {
		t.Errorf("unknown agent: exit %d, error %+v", code, e)
	}
}

func TestSkillInstallFollowsTheAgentsHomeVariables(t *testing.T) {
	env, vars := testEnv(t)
	home, _ := env.Home()
	claude, codex := t.TempDir(), t.TempDir()
	vars["CLAUDE_CONFIG_DIR"], vars["CODEX_HOME"] = claude, codex
	got := installed(t, mustRun(t, env, "skill", "install"))
	want := []string{filepath.Join(claude, "skills", "aikido-dojo", "SKILL.md"), filepath.Join(home, ".agents", "skills", "aikido-dojo", "SKILL.md")}
	if len(got) != 2 || got[0].Path != want[0] || got[1].Path != want[1] {
		t.Errorf("installed %+v, want %q", got, want)
	}
}

// A re-run after an update replaces the old copy and leaves no temp file behind.
func TestSkillInstallReplacesAnOlderCopy(t *testing.T) {
	env, _ := testEnv(t)
	dir := t.TempDir()
	old := filepath.Join(dir, "aikido-dojo", "SKILL.md")
	mkdir(t, filepath.Dir(old))
	if err := os.WriteFile(old, []byte("an older primer"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := installed(t, mustRun(t, env, "skill", "install", "--dir", dir)); len(got) != 1 || got[0].Path != old || got[0].Agents != nil {
		t.Errorf("installed %+v, want %q", got, old)
	}
	if got := readFile(t, old); got != string(env.Primer) {
		t.Errorf("SKILL.md = %q, want the primer", got)
	}
	entries, err := os.ReadDir(filepath.Dir(old))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Errorf("skill directory holds %d entries, want only SKILL.md", len(entries))
	}
}

// Read-only mode guards Aikido, not local files, as with docs sync.
func TestSkillInstallRunsInReadOnlyMode(t *testing.T) {
	env, vars := testEnv(t)
	vars[config.EnvReadOnly] = "1"
	mustRun(t, env, "--read-only", "skill", "install", "--dir", t.TempDir())
}

func TestSkillInstallFailsOnAnUnwritableDir(t *testing.T) {
	env, _ := testEnv(t)
	file := filepath.Join(t.TempDir(), "not-a-dir")
	if err := os.WriteFile(file, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	_, stderr, code := run(t, env, "skill", "install", "--dir", file)
	if e := errorOf(t, stderr); code != clierr.ExitUnexpected || e.Code != "skill_install_failed" || e.Hint == "" {
		t.Errorf("exit %d, error %+v", code, e)
	}
}

func TestSkillInstallUsageErrors(t *testing.T) {
	env, _ := testEnv(t)
	for _, args := range [][]string{
		{"skill", "install", "extra"},
		{"skill", "instal"},
		{"skill", "install", "--dir", t.TempDir(), "--agent", "claude"},
	} {
		if _, _, code := run(t, env, args...); code != clierr.ExitUsage {
			t.Errorf("%q: exit %d, want usage", args, code)
		}
	}
}
