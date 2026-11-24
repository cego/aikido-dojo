package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cego/aikido-dojo/internal/clierr"
)

// skillName is the directory the skill goes in, and the name an agent sees it by.
const skillName = "aikido-dojo"

// skillAgent is an agent that loads skills from a directory in the user's home.
type skillAgent struct {
	name string
	home []string // where the agent keeps its config, relative to the home directory; it is installed when one exists
	env  string   // a variable that moves home, when the agent has one
	dir  []string // the skills directory, relative to the home directory; nil means home/skills
}

// The shared ~/.agents/skills serves every agent that reads it, so each gets one copy.
var skillAgents = []skillAgent{
	{name: "claude", home: []string{".claude"}, env: "CLAUDE_CONFIG_DIR"},
	{name: "codex", home: []string{".codex"}, env: "CODEX_HOME", dir: []string{".agents", "skills"}},
	{name: "copilot", home: []string{".copilot"}, dir: []string{".agents", "skills"}},
	{name: "gemini", home: []string{".gemini"}, dir: []string{".agents", "skills"}},
	{name: "cursor", home: []string{".cursor"}},
	{name: "opencode", home: []string{".config", "opencode"}},
	{name: "windsurf", home: []string{".codeium", "windsurf"}},
}

func skillAgentNames() []string {
	names := make([]string, len(skillAgents))
	for i, s := range skillAgents {
		names[i] = s.name
	}
	return names
}

// installed is the agent's config directory, and whether it exists.
func (s skillAgent) installed(home string, getenv func(string) string) (string, bool) {
	root := filepath.Join(append([]string{home}, s.home...)...)
	if s.env != "" && getenv(s.env) != "" {
		root = getenv(s.env)
	}
	info, err := os.Stat(root)
	return root, err == nil && info.IsDir()
}

func (s skillAgent) skillsDir(home string, getenv func(string) string) string {
	if s.dir != nil {
		return filepath.Join(append([]string{home}, s.dir...)...)
	}
	root, _ := s.installed(home, getenv)
	return filepath.Join(root, "skills")
}

type installedSkill struct {
	Path   string   `json:"path"`
	Agents []string `json:"agents,omitempty"`
}

func (a *app) skillCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "skill",
		Short:                      "Install AGENTS.md as an agent skill",
		Args:                       unknownCommand,
		RunE:                       showHelp,
		SuggestionsMinimumDistance: 2,
	}
	cmd.AddCommand(a.skillInstallCmd())
	return cmd
}

func (a *app) skillInstallCmd() *cobra.Command {
	var dir string
	var agents []string
	names := strings.Join(skillAgentNames(), ", ")
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Write AGENTS.md to AI agents' skills directories as aikido-dojo/SKILL.md",
		Long: "skill install writes the AGENTS.md built into aikido-dojo to <skills dir>/aikido-dojo/SKILL.md, so an AI " +
			"agent loads it as a skill. It replaces the file a previous install wrote; run it again after updating aikido-dojo.\n\n" +
			"Without --agent it installs for every agent whose config directory exists in your home directory. The agents " +
			"are " + names + ". Codex, Copilot and Gemini share ~/.agents/skills, which also serves Cline, Pi, Warp and Zed. " +
			"--dir names any other skills directory.\n\n" +
			"It prints a JSON array of {path, agents}: each file it wrote, and the agents that read it.",
		Example: "  aikido-dojo skill install\n  aikido-dojo skill install --agent claude --agent codex\n" +
			"  aikido-dojo skill install --dir ~/.config/agents/skills",
		Args: exactArgs(nil),
		RunE: func(cmd *cobra.Command, _ []string) error { return a.skillInstall(cmd.Context(), dir, agents) },
	}
	cmd.Flags().StringVar(&dir, "dir", "", "a skills directory to install into instead of the agents'")
	cmd.Flags().StringSliceVar(&agents, "agent", nil, "an agent to install for, repeatable: "+names+" (default: every agent installed)")
	return cmd
}

func (a *app) skillInstall(ctx context.Context, dir string, names []string) error {
	if len(a.env.Primer) == 0 {
		return errors.New("this build of aikido-dojo carries no AGENTS.md")
	}
	targets, err := a.skillTargets(dir, names)
	if err != nil {
		return err
	}
	for i := range targets {
		path, err := writeSkill(filepath.Join(targets[i].Path, skillName), a.env.Primer)
		if err != nil {
			return &clierr.Error{Code: "skill_install_failed", Message: "skill install failed", Err: err,
				Hint: "check that the skills directory is writable, or name another with --dir", Exit: clierr.ExitUnexpected}
		}
		targets[i].Path = path
	}
	return a.printJSON(ctx, targets)
}

// skillTargets lists the skills directories to write, once each, with the agents that read each.
func (a *app) skillTargets(dir string, names []string) ([]installedSkill, error) {
	if dir != "" && len(names) > 0 {
		return nil, &clierr.Error{Code: "usage", Message: "--dir and --agent name the skills directory two ways",
			Hint: "pass one of them", Exit: clierr.ExitUsage}
	}
	if dir != "" {
		return []installedSkill{{Path: dir}}, nil
	}
	home, err := a.env.Home()
	if err != nil {
		return nil, &clierr.Error{Code: "no_home", Message: "find the home directory for the skills directories", Err: err,
			Hint: "name the skills directory with --dir", Exit: clierr.ExitUsage}
	}
	var chosen []skillAgent
	for _, n := range names {
		i := slices.IndexFunc(skillAgents, func(s skillAgent) bool { return s.name == strings.ToLower(strings.TrimSpace(n)) })
		if i < 0 {
			return nil, &clierr.Error{Code: "invalid_input", Message: fmt.Sprintf("--agent %q is not an agent skill install knows", n),
				Hint: "use one of " + strings.Join(skillAgentNames(), ", ") + ", or name its skills directory with --dir", Exit: clierr.ExitUsage}
		}
		chosen = append(chosen, skillAgents[i])
	}
	if len(names) == 0 {
		for _, s := range skillAgents {
			if _, ok := s.installed(home, a.env.Getenv); ok {
				chosen = append(chosen, s)
			}
		}
		if len(chosen) == 0 {
			return nil, &clierr.Error{Code: "no_agent", Message: "found no AI agent's config directory in " + home,
				Hint: "name one with --agent (" + strings.Join(skillAgentNames(), ", ") + "), or a skills directory with --dir", Exit: clierr.ExitUsage}
		}
	}
	var targets []installedSkill
	for _, s := range chosen {
		d := s.skillsDir(home, a.env.Getenv)
		if i := slices.IndexFunc(targets, func(t installedSkill) bool { return t.Path == d }); i >= 0 {
			if !slices.Contains(targets[i].Agents, s.name) {
				targets[i].Agents = append(targets[i].Agents, s.name)
			}
			continue
		}
		targets = append(targets, installedSkill{Path: d, Agents: []string{s.name}})
	}
	return targets, nil
}

// writeSkill replaces dir/SKILL.md in one rename, so an agent never reads half a file.
func writeSkill(dir string, body []byte) (string, error) {
	path := filepath.Join(dir, "SKILL.md")
	if err := os.MkdirAll(dir, 0o755); err != nil { //nolint:gosec // agents read skills as the user; none is secret
		return "", fmt.Errorf("create %s: %w", dir, err)
	}
	tmp, err := os.CreateTemp(dir, ".SKILL-*.md")
	if err != nil {
		return "", fmt.Errorf("write %s: %w", path, err)
	}
	defer os.Remove(tmp.Name()) // fails harmlessly once the rename has moved it
	if _, err := tmp.Write(body); err != nil {
		return "", fmt.Errorf("write %s: %w", path, errors.Join(err, tmp.Close()))
	}
	if err := tmp.Chmod(0o644); err != nil {
		return "", fmt.Errorf("write %s: %w", path, errors.Join(err, tmp.Close()))
	}
	if err := tmp.Close(); err != nil {
		return "", fmt.Errorf("write %s: %w", path, err)
	}
	if err := os.Rename(tmp.Name(), path); err != nil {
		return "", fmt.Errorf("write %s: %w", path, err)
	}
	return path, nil
}
