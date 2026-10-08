package cli

import (
	"cmp"
	"context"
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cego/aikido-dojo/internal/auth"
	"github.com/cego/aikido-dojo/internal/catalog"
	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/config"
)

type loginResult struct {
	Profile  string `json:"profile"`
	ClientID string `json:"client_id"`
	Region   string `json:"region"`
}

// authStatus is who calls run as. Source is "keychain" for a stored profile
// and "environment" for the AIKIDO_DOJO_CLIENT_ID pair, which has no profile.
type authStatus struct {
	Profile  string      `json:"profile,omitempty"`
	Source   string      `json:"source"`
	Method   string      `json:"method"`
	ClientID string      `json:"client_id"`
	Region   string      `json:"region"`
	Scopes   scopeReport `json:"scopes"`
}

type scopeReport struct {
	Granted []string `json:"granted"`
	Denied  []string `json:"denied"`
	Unknown []string `json:"unknown"`
}

type logoutResult struct {
	Profile string `json:"profile"`
	Removed bool   `json:"removed"`
}

func (a *app) authCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "auth",
		Short:                      "Store, check and remove an API client's credentials",
		Args:                       unknownCommand,
		RunE:                       showHelp,
		SuggestionsMinimumDistance: 2,
	}
	cmd.AddCommand(a.loginCmd(), a.statusCmd(), a.logoutCmd())
	return cmd
}

func (a *app) loginCmd() *cobra.Command {
	var clientID, region string
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Store an API client's credentials as a profile, once Aikido accepts them",
		Long: "login stores an API client's credentials as a profile: the client ID and region in the config file, " +
			"the client secret in the OS keychain. It checks them with one token request first, and stores nothing " +
			"if Aikido refuses them.\n\n" +
			"The client ID comes from --client-id, then AIKIDO_DOJO_CLIENT_ID, then the profile's entry; the region " +
			"from --region, then AIKIDO_DOJO_REGION, then the profile's entry, then eu. The secret comes from " +
			"AIKIDO_DOJO_CLIENT_SECRET, else a hidden prompt. No flag takes it: argv and shell history would keep it.\n\n" +
			"The profile is --profile, then AIKIDO_DOJO_PROFILE, then the default profile, then default; the first " +
			"profile becomes the default. Create the API client in Aikido's workspace settings with only the scopes you need.",
		Args: exactArgs(nil),
		RunE: func(cmd *cobra.Command, _ []string) error { return a.login(cmd.Context(), clientID, region) },
	}
	cmd.Flags().StringVar(&clientID, "client-id", "", "the API client's ID, from Aikido's workspace settings")
	cmd.Flags().StringVar(&region, "region", "", "the workspace's region: eu, us, au or me")
	return cmd
}

func (a *app) login(ctx context.Context, flagID, flagRegion string) error {
	if err := a.refuseWrite(http.MethodPost, "auth login"); err != nil {
		return err
	}
	path, f, err := a.configFile()
	if err != nil {
		return err
	}
	name := cmp.Or(a.profile, a.env.Getenv(config.EnvProfile), f.DefaultProfile, "default")
	old, exists := f.Profiles[name]
	region := cmp.Or(flagRegion, a.env.Getenv(config.EnvRegion), old.Region, config.DefaultRegion)
	p := config.Profile{ClientID: cmp.Or(flagID, a.env.Getenv(config.EnvClientID), old.ClientID), Region: region}
	// An entry that omits its region keeps omitting it while the region is the default.
	if exists && old.Region == "" && region == config.DefaultRegion {
		p.Region = ""
	}
	if p.ClientID == "" {
		return &clierr.Error{Code: "no_client_id", Message: "no client ID for profile " + strconv.Quote(name),
			Hint: "pass --client-id; Aikido shows it with the API client in its workspace settings", Exit: clierr.ExitUsage}
	}
	host, err := config.RegionHost(region)
	if err != nil {
		return err
	}
	secret := a.env.Getenv(config.EnvClientSecret)
	if secret == "" {
		if secret, err = a.env.ReadSecret(ctx, fmt.Sprintf("Client secret for profile %q: ", name)); err != nil {
			return err
		}
	}
	if strings.TrimSpace(secret) == "" {
		return &clierr.Error{Code: "empty_secret", Message: "the client secret is empty",
			Hint: "enter the API client's secret from Aikido's workspace settings", Exit: clierr.ExitUsage}
	}
	r := config.Resolved{Profile: name, ClientID: p.ClientID, Secret: secret, Region: region, Host: host}
	record := func() error {
		if p == old {
			return nil
		}
		if len(f.Profiles) == 0 {
			f.Profiles = map[string]config.Profile{}
			f.DefaultProfile = cmp.Or(f.DefaultProfile, name)
		}
		f.Profiles[name] = p
		if err := config.Save(path, f); err != nil {
			return &clierr.Error{Code: "config_not_saved", Message: "save profile " + strconv.Quote(name), Err: err,
				Hint: "nothing was stored; pass --config with a file you can write to keep the profile there", Exit: clierr.ExitUnexpected}
		}
		return nil
	}
	if err := auth.Login(ctx, newHTTPClient(a.env, host, a.debug), r, record); err != nil {
		return err
	}
	return a.printJSON(loginResult{Profile: name, ClientID: p.ClientID, Region: region})
}

func (a *app) statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show who calls run as and which scopes their token holds",
		Long: "status shows the profile calls run as, or that the AIKIDO_DOJO_CLIENT_ID pair is used; where its secret " +
			"is kept; its client ID and region; and the scopes its access token holds. It requests a token only when " +
			"none is cached, which also checks the credentials, and calls nothing else.\n\n" +
			"scopes.granted is what the token holds and scopes.denied what some command needs that it lacks. " +
			"Scopes are granted on the API client in Aikido's workspace settings.",
		Args: exactArgs(nil),
		RunE: func(cmd *cobra.Command, _ []string) error { return a.status(cmd.Context()) },
	}
}

func (a *app) status(ctx context.Context) error {
	r, err := a.resolve()
	if err != nil {
		return err
	}
	src, err := auth.NewSource(newHTTPClient(a.env, r.Host, a.debug), r)
	if err != nil {
		return err
	}
	tok, err := src.Token(ctx)
	if err != nil {
		return err
	}
	claim, ok := auth.Scopes(tok)
	if !ok {
		clierr.Warn(a.env.Stderr, "the access token names no scopes, so they are listed as unknown; "+
			"a command whose scope is missing fails with exit 4 and names it")
	}
	out := authStatus{Profile: r.Profile, Source: "keychain", Method: "client_credentials", ClientID: r.ClientID, Region: r.Region, Scopes: scopesOf(claim, ok)}
	if r.Profile == "" {
		out.Source = "environment"
	}
	return a.printJSON(out)
}

// scopesOf sorts every scope some command needs into granted and denied by
// the token's claim, or into unknown when it has none. granted also keeps
// claimed scopes no command needs.
func scopesOf(claim []string, ok bool) scopeReport {
	needed := map[string]bool{}
	for _, op := range catalog.All {
		if op.Scope != "" {
			needed[op.Scope] = true
		}
	}
	all := slices.Sorted(maps.Keys(needed))
	s := scopeReport{Granted: []string{}, Denied: []string{}, Unknown: []string{}}
	if !ok {
		s.Unknown = all
		return s
	}
	s.Granted = append(s.Granted, slices.Compact(slices.Sorted(slices.Values(claim)))...)
	for _, scope := range all {
		if !slices.Contains(claim, scope) {
			s.Denied = append(s.Denied, scope)
		}
	}
	return s
}

func (a *app) logoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Delete a profile's secret and cached token from the keychain",
		Long: "logout deletes the client secret and the cached access token that auth login stored for the profile: " +
			"--profile, then AIKIDO_DOJO_PROFILE, then the default profile. The profile's entry in the config file " +
			"stays, so auth login restores it with the secret alone.\n\n" +
			"A token Aikido already issued stays valid until it expires, within an hour. To revoke the API client " +
			"itself, delete or rotate it in Aikido's workspace settings.",
		Args: exactArgs(nil),
		RunE: func(*cobra.Command, []string) error { return a.logout() },
	}
}

func (a *app) logout() error {
	if err := a.refuseWrite(http.MethodPost, "auth logout"); err != nil {
		return err
	}
	_, f, err := a.configFile()
	if err != nil {
		return err
	}
	name := cmp.Or(a.profile, a.env.Getenv(config.EnvProfile), f.DefaultProfile)
	if name == "" {
		return &clierr.Error{Code: "usage", Message: "no profile to log out of",
			Hint: "pass --profile <name>; the AIKIDO_DOJO_CLIENT_ID pair is never stored, so it has nothing to delete", Exit: clierr.ExitUsage}
	}
	removed, err := auth.Forget(name)
	if err != nil {
		return err
	}
	return a.printJSON(logoutResult{Profile: name, Removed: removed})
}
