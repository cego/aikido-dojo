// Package overlay adds what the spec lacks: each operation's command name,
// how it pages, and corrections to the spec. Every correction carries a
// comment with its evidence; a test enforces that.
package overlay

import "github.com/cego/aikido-dojo/internal/api"

type Op struct {
	Cmd  string // "<resource> <verb>"
	Page Page   // zero for an operation that doesn't page
	// Destructive marks an operation whose verb doesn't already say so; every
	// delete, deactivate and rotate is destructive.
	Destructive bool

	// Corrections and additions to the spec. Each needs an evidence comment.
	Scope          string   // the scope the operation really needs
	ResponseArray  bool     // the 2xx body is an array of what the spec types as one object
	PageSize       int      // the page size to request when the spec gives no maximum
	NotPaged       bool     // takes page parameters but returns everything without them
	Secret         []string // body fields holding credentials, kept off the command line
	Help           string   // added to the command's help
	BadRequestHint string   // the hint for a 400 answer
	Unbounded      []string // parameters whose minimum and maximum the spec got wrong
}

// Page says where a list's items are and how its last page is marked.
type Page struct {
	End   api.End
	Items string // the field holding the items when the body is an object
	More  string // api.EndField: the boolean field saying more pages exist
}

var (
	UntilEmpty  = Page{End: api.EndEmpty}
	UntilHeader = Page{End: api.EndHeader}
)

var Ops = map[string]Op{
	// Spec gives per_page no maximum; its schema description says max 50.
	"listActivityLog": {Cmd: "activity-log list", Page: UntilEmpty, PageSize: 50},
	// The spec says to page by X-Has-Next-Page, but live on 2026-10-02 no page carried that header and a filtered page came back full, so an empty page ends the list.
	"listAikidoImages":                {Cmd: "aikido-image list", Page: UntilEmpty},
	"getContainerAutofixSettings":     {Cmd: "autofix-container-config get"},
	"updateContainerAutofixSettings":  {Cmd: "autofix-container-config update"},
	"getDependencyAutofixSettings":    {Cmd: "autofix-dependency-config get"},
	"updateDependencyAutofixSettings": {Cmd: "autofix-dependency-config update"},
	"getPentestAutofixSettings":       {Cmd: "autofix-pentest-config get"},
	"updatePentestAutofixSettings":    {Cmd: "autofix-pentest-config update"},
	"getSastAutofixSettings":          {Cmd: "autofix-sast-config get"},
	"updateSastAutofixSettings":       {Cmd: "autofix-sast-config update"},
	"listAutofixHistory":              {Cmd: "autofix-task list", Page: UntilHeader},
	"listAutofixTaskIssues":           {Cmd: "autofix-task-issue list"},
	"listBrokerApps":                  {Cmd: "broker-app list"},
	"addBrokerResourceUrl":            {Cmd: "broker-resource-url create"},
	"deleteBrokerResourceUrl":         {Cmd: "broker-resource-url delete"},
	"addBugBountyReport":              {Cmd: "bug-bounty-report create"},
	"getChangelogSummary":             {Cmd: "changelog-summary get"},
	"removeCloud":                     {Cmd: "cloud delete"},
	"listClouds":                      {Cmd: "cloud list", Page: UntilEmpty},
	// Spec gives limit no maximum; live on 2026-09-30, limit=100 answered 200 and limit=100000 answered 400.
	"getCloudAssets":  {Cmd: "cloud-asset list", Page: Page{End: api.EndField, Items: "assets", More: "hasMore"}, PageSize: 100},
	"connectAwsCloud": {Cmd: "cloud-aws create"},
	// Spec: key_value is the generated secret of the registered application.
	"connectAzureCloud": {Cmd: "cloud-azure create", Secret: []string{"key_value"}},
	// Spec: client_secret is the generated secret of the registered application.
	"updateAzureCloudCredentials": {Cmd: "cloud-azure-credentials update", Secret: []string{"client_secret"}},
	// Spec: access_key is the service account's access key.
	"connectGcpCloud":              {Cmd: "cloud-gcp create", Secret: []string{"access_key"}},
	"createKubernetesCloud":        {Cmd: "cloud-kubernetes create"},
	"listCloudRules":               {Cmd: "cloud-rule list"},
	"getCodeCoverageForRepository": {Cmd: "code-coverage get"},
	// Spec: {code_coverage_repositories} with page/per_page; not checked live, the READ clients lack code_coverage:read (2026-10-08).
	"listAllCodeCoverageRepositories":       {Cmd: "code-coverage list", Page: Page{End: api.EndEmpty, Items: "code_coverage_repositories"}},
	"listCodeQualityFullRepoFindings":       {Cmd: "code-quality-finding list", Page: Page{End: api.EndHeader, Items: "findings"}},
	"listCodeQualityFindingsForPullRequest": {Cmd: "code-quality-pr-finding list"},
	// Spec: access_token is the access token used for code scanning.
	"setCodeScanningAccessToken":  {Cmd: "code-scanning-token update", Secret: []string{"access_token"}},
	"getCisComplianceOverview":    {Cmd: "compliance-cis get"},
	"getCisAwsComplianceOverview": {Cmd: "compliance-cis-aws get"},
	"getDoraComplianceOverview":   {Cmd: "compliance-dora get"},
	"getFedrampDastReport":        {Cmd: "compliance-fedramp-dast get"},
	"getGdprComplianceOverview":   {Cmd: "compliance-gdpr get"},
	"getIsoComplianceOverview":    {Cmd: "compliance-iso get"},
	"getNis2ComplianceOverview":   {Cmd: "compliance-nis2 get"},
	"getPciComplianceOverview":    {Cmd: "compliance-pci get"},
	"getSoc2ComplianceOverview":   {Cmd: "compliance-soc2 get"},
	"activateContainer":           {Cmd: "container activate"},
	"cloneContainer":              {Cmd: "container clone"},
	"addPublicContainer":          {Cmd: "container create"},
	"deactivateContainer":         {Cmd: "container deactivate"},
	"deleteContainer":             {Cmd: "container delete"},
	"getContainerRepo":            {Cmd: "container get"},
	// Spec types the 200 as one object; a live call on 2026-09-30 returned a JSON array.
	"listContainerRepos": {Cmd: "container list", Page: UntilEmpty, ResponseArray: true},
	"scanContainer":      {Cmd: "container scan"},
	// Spec says container:write, which it never declares; the token's scope claim in the auth spike (2026-09-29) has containers:write.
	"updateContainerInternetConnection": {Cmd: "container-connectivity update", Scope: "containers:write"},
	"addContainerLabel":                 {Cmd: "container-label create"},
	"removeContainerLabel":              {Cmd: "container-label delete"},
	"updateContainerLabel":              {Cmd: "container-label update"},
	"exportContainerRepoLicenses":       {Cmd: "container-license export"},
	"getContainerRegistry":              {Cmd: "container-registry get"},
	// Spec: accesstoken is the registry's access token.
	"addAzureContainerRegistry": {Cmd: "container-registry-acr create", Secret: []string{"accesstoken"}},
	// Spec: service_account_key is the service account's access key.
	"addGcpArtifactRegistry": {Cmd: "container-registry-gcp create", Secret: []string{"service_account_key"}},
	// Spec: accesstoken is a GitLab access token.
	"addGitlabSelfManagedContainerRegistry": {Cmd: "container-registry-gitlab create", Secret: []string{"accesstoken"}},
	"linkCodeRepoToContainer":               {Cmd: "container-repo-link create"},
	"unlinkCodeRepoToContainer":             {Cmd: "container-repo-link delete"},
	"listContainerRepoRunners":              {Cmd: "container-runner list"},
	"exportRawContainerSbom":                {Cmd: "container-sbom export"},
	"uploadContainerSBOM":                   {Cmd: "container-sbom import"},
	"generateContainerSBOM":                 {Cmd: "container-sbom-bulk export"},
	// Spec says container:write, which it never declares; the token's scope claim in the auth spike (2026-09-29) has containers:write.
	"updateContainerSensitivity": {Cmd: "container-sensitivity update", Scope: "containers:write"},
	"updateTagFilter":            {Cmd: "container-tag-filter update"},
	"getCveDetails":              {Cmd: "cve get"},
	"listCweDescriptions":        {Cmd: "cwe list"},
	"createDomain":               {Cmd: "domain create"},
	"removeDomain":               {Cmd: "domain delete"},
	"listDomains":                {Cmd: "domain list", Page: UntilEmpty},
	"startDomainScan":            {Cmd: "domain scan"},
	// Spec: http_headers is "The authentication headers", name/value pairs a scan sends.
	"updateDomainAuthenticationHeaders": {Cmd: "domain-auth-header update", Secret: []string{"http_headers"}},
	"updateDomainOpenAPISpec":           {Cmd: "domain-openapi-spec update"},
	// Spec: custom_scan_headers is "The authentication headers", name/value pairs a scan sends.
	"updateDomainCustomScanHeaders": {Cmd: "domain-scan-header update", Secret: []string{"custom_scan_headers"}},
	// Spec gives per_page no maximum; its description says max 50.
	"listEndpointProtectionActivityLogs":      {Cmd: "endpoint-activity-log list", Page: UntilEmpty, PageSize: 50},
	"removeEndpointProtectionDevice":          {Cmd: "endpoint-device delete"},
	"listEndpointProtectionDevices":           {Cmd: "endpoint-device list"},
	"addEndpointProtectionException":          {Cmd: "endpoint-exception create"},
	"removeEndpointProtectionException":       {Cmd: "endpoint-exception delete"},
	"listEndpointProtectionExceptions":        {Cmd: "endpoint-exception list"},
	"listEndpointProtectionInstalledPackages": {Cmd: "endpoint-package list", Page: UntilEmpty},
	"listEndpointProtectionPermissionGroups":  {Cmd: "endpoint-permission-group list"},
	// Live 2026-10-08: an array, empty past the last page; per_page 100 is honoured.
	"listEolRuntimes":                 {Cmd: "eol-runtime list", Page: UntilEmpty},
	"createApp":                       {Cmd: "firewall-app create"},
	"deleteApp":                       {Cmd: "firewall-app delete"},
	"getApp":                          {Cmd: "firewall-app get"},
	"listApps":                        {Cmd: "firewall-app list"},
	"updateApp":                       {Cmd: "firewall-app update"},
	"rotateAppToken":                  {Cmd: "firewall-app-token rotate"},
	"updateBlocking":                  {Cmd: "firewall-blocking update"},
	"getBotLists":                     {Cmd: "firewall-bot-list get"},
	"updateBotLists":                  {Cmd: "firewall-bot-list update"},
	"getCountries":                    {Cmd: "firewall-country-list get"},
	"updateCountries":                 {Cmd: "firewall-country-list update"},
	"getEvent":                        {Cmd: "firewall-event get"},
	"updateIpBlocklist":               {Cmd: "firewall-ip-blocklist update"},
	"getIpLists":                      {Cmd: "firewall-threat-list get"},
	"updateIpLists":                   {Cmd: "firewall-threat-list update"},
	"listZenUsers":                    {Cmd: "firewall-user list", Page: Page{End: api.EndEmpty, Items: "users"}},
	"updateUser":                      {Cmd: "firewall-user update"},
	"listIacRules":                    {Cmd: "iac-rule list"},
	"listIdeAdoption":                 {Cmd: "ide-adoption list"},
	"requestContainerImage":           {Cmd: "image-request create"},
	"listImageRequests":               {Cmd: "image-request list"},
	"markLicenseAsInternalForPackage": {Cmd: "internal-package create"},
	// Spec: leaving page empty exports every issue, which is what an export is for.
	// Help: the handover's live use (2026-09) found both filter traps.
	"exportIssues": {Cmd: "issue export", NotPaged: true,
		Help: "Aikido ignores --filter-issue-group-id, and --filter-status all without --filter-code-repo-id fails with 400 \"Request too big\"."},
	"getIssueDetail": {Cmd: "issue get"},
	// Handover, from live use (2026-09): Aikido accepts reason and discards it.
	"ignoreIssue": {Cmd: "issue ignore", Help: "Aikido discards --reason: the issue records \"Ignored via API\". To keep a rationale, add one to its group with issue-group-note create."},
	// Live on 2026-09-30: without the feature Aikido answers 400 "This feature is not enabled on your workspace".
	"getIssueDetailsBulk":           {Cmd: "issue list", BadRequestHint: "Aikido support must enable bulk issue details for the workspace; issue get <issue_id> works without it"},
	"snoozeIssue":                   {Cmd: "issue snooze"},
	"solveIssue":                    {Cmd: "issue solve"},
	"UnignoreIssue":                 {Cmd: "issue unignore"},
	"UnsnoozeIssue":                 {Cmd: "issue unsnooze"},
	"getIssueCounts":                {Cmd: "issue-count get"},
	"getIssueCveExploitability":     {Cmd: "issue-cve-exploitability get"},
	"getIssueCveExploitabilityBulk": {Cmd: "issue-cve-exploitability list"},
	"getIssueGroupDetails":          {Cmd: "issue-group get"},
	// Handover, from live use (2026-09): Aikido accepts reason and discards it.
	"ignoreIssueGroup": {Cmd: "issue-group ignore", Help: "Aikido discards --reason: the group records \"Ignored via API\". To keep a rationale, add one with issue-group-note create."},
	// Handover, from live use (2026-09): a fully ignored group drops out of this list.
	"listOpenIssueGroups": {Cmd: "issue-group list", Page: UntilHeader, Help: "Lists open groups only. A group whose issues are all ignored is not listed, so a repo can look clean while ignored issues exist."},
	"snoozeIssueGroup":    {Cmd: "issue-group snooze"},
	"unignoreIssueGroup":  {Cmd: "issue-group unignore"},
	"unsnoozeIssueGroup":  {Cmd: "issue-group unsnooze"},
	"addNoteToIssueGroup": {Cmd: "issue-group-note create"},
	// Spec types the 200 as one object; a live call on 2026-09-30 returned a JSON array.
	"listNotesForIssueGroup":     {Cmd: "issue-group-note list", ResponseArray: true},
	"adjustGroupSeverity":        {Cmd: "issue-group-severity update"},
	"linkTaskToIssueGroup":       {Cmd: "issue-group-task create"},
	"getIssueGroupTasks":         {Cmd: "issue-group-task list"},
	"getIssueReachability":       {Cmd: "issue-reachability get"},
	"adjustSeverity":             {Cmd: "issue-severity update"},
	"listLicenses":               {Cmd: "license list", Page: UntilEmpty},
	"getLatestLocalScanInfo":     {Cmd: "local-scanner-version get"},
	"getMalwarePackages":         {Cmd: "malware-package list", Page: UntilEmpty},
	"listMobileRules":            {Cmd: "mobile-rule list"},
	"getOpenApiSpec":             {Cmd: "openapi-spec get"},
	"getPackageHealthBulk":       {Cmd: "package-health list"},
	"overwriteLicenseForPackage": {Cmd: "package-license update"},
	// Spec: custom_headers maps header names to the values a pentest sends, which hold its credentials.
	"createPentestDraftAssessment": {Cmd: "pentest-assessment create", Secret: []string{"custom_headers"}},
	"getPentestAssessmentDetail":   {Cmd: "pentest-assessment get"},
	"getAttackAnalysis":            {Cmd: "pentest-attack-analysis get"},
	"getCiScanDetail":              {Cmd: "pr-check get"},
	"listCiScans":                  {Cmd: "pr-check list", Page: UntilEmpty},
	"listCiScanIssueActions":       {Cmd: "pr-check-action list", Page: UntilEmpty},
	// Spec types the 200 as one object; a live call on 2026-09-30 returned a JSON array.
	"listPRChecksConfigurations": {Cmd: "pr-check-config list", Page: UntilEmpty, ResponseArray: true},
	"savePRChecksConfiguration":  {Cmd: "pr-check-config update"},
	// Replaces the PR-check configuration of every repository in one call.
	"savePRChecksConfigurationForAllRepos": {Cmd: "pr-check-config-all update", Destructive: true},
	"getDefaultPRChecksConfiguration":      {Cmd: "pr-check-config-default get"},
	"saveDefaultPRChecksConfiguration":     {Cmd: "pr-check-config-default update"},
	// Spec: payload is the registry connection or .npmrc content, which carries the registry's credentials.
	"addPrivateRegistry":        {Cmd: "private-registry update", Secret: []string{"payload"}},
	"activateCodeRepo":          {Cmd: "repo activate"},
	"cloneCodeRepo":             {Cmd: "repo clone"},
	"deactivateCodeRepo":        {Cmd: "repo deactivate"},
	"deleteCodeRepo":            {Cmd: "repo delete"},
	"detailCodeRepo":            {Cmd: "repo get"},
	"importNewCodeRepositories": {Cmd: "repo import"},
	// Spec types the 200 as one object; a live call on 2026-09-29 returned a JSON array.
	"listCodeRepos":                 {Cmd: "repo list", Page: UntilEmpty, ResponseArray: true},
	"scanCodeRepo":                  {Cmd: "repo scan"},
	"updateCodeRepoConnectivity":    {Cmd: "repo-connectivity update"},
	"updateCodeRepoDevDepScan":      {Cmd: "repo-dev-dep-scan update"},
	"addCodeRepoExcludePath":        {Cmd: "repo-exclude-path create"},
	"removeCodeRepoExcludePath":     {Cmd: "repo-exclude-path delete"},
	"addCodeRepoLabel":              {Cmd: "repo-label create"},
	"removeCodeRepoLabel":           {Cmd: "repo-label delete"},
	"updateCodeRepoLabel":           {Cmd: "repo-label update"},
	"exportCodeRepoLicenses":        {Cmd: "repo-license export"},
	"updateCodeRepoSensitivity":     {Cmd: "repo-sensitivity update"},
	"exportReportPdf":               {Cmd: "report export"},
	"createCustomRule":              {Cmd: "sast-custom-rule create"},
	"removeCustomRule":              {Cmd: "sast-custom-rule delete"},
	"getCustomRule":                 {Cmd: "sast-custom-rule get"},
	"listCustomRules":               {Cmd: "sast-custom-rule list"},
	"editCustomRule":                {Cmd: "sast-custom-rule update"},
	"listSastRules":                 {Cmd: "sast-rule list"},
	"getSlaSettings":                {Cmd: "sla-config get"},
	"addSubdomain":                  {Cmd: "subdomain create"},
	"listSubdomains":                {Cmd: "subdomain list", Page: UntilEmpty},
	"getTasks":                      {Cmd: "task list"},
	"listTaskTrackingIntegrations":  {Cmd: "task-integration list"},
	"getProjectMapping":             {Cmd: "task-mapping get"},
	"getProjects":                   {Cmd: "task-project list"},
	"mapReposToProject":             {Cmd: "task-repo-mapping update"},
	"mapTeamsToProject":             {Cmd: "task-team-mapping update"},
	"createTeam":                    {Cmd: "team create"},
	"deleteTeam":                    {Cmd: "team delete"},
	"listTeams":                     {Cmd: "team list", Page: UntilEmpty},
	"updateTeam":                    {Cmd: "team update"},
	"exportCodeRepoLicensesForTeam": {Cmd: "team-license export"},
	"addUserToTeam":                 {Cmd: "team-member create"},
	"removeUserFromTeam":            {Cmd: "team-member delete"},
	"updateTeamRepoPathLimitation":  {Cmd: "team-repo-path-limitation update"},
	"linkResourceToTeam":            {Cmd: "team-resource create"},
	"unlinkResourceFromTeam":        {Cmd: "team-resource delete"},
	// Spec bounds user_id to 0..1, copied from listUsers.include_inactive; live on 2026-10-02, GET /users/{id} with an id above 1 answered 200.
	"getUser":              {Cmd: "user get", Unbounded: []string{"user_id"}},
	"listUsers":            {Cmd: "user list"},
	"listPendingInvites":   {Cmd: "user-invite list"},
	"listUserLoginHistory": {Cmd: "user-login list", Page: UntilEmpty},
	// Spec: the description says the list pages by X-Has-Next-Page, but it defines no page parameters; not checked live (the read clients lack request_inbox:read, 2026-10-02).
	"listUserRequests": {Cmd: "user-request list",
		Help: "Aikido documents this list as paged by an X-Has-Next-Page header but defines no page parameters, so the command returns what one call returns."},
	"reviewUserRequest": {Cmd: "user-request review"},
	// Spec bounds user_id to 0..1, copied from listUsers.include_inactive; live on 2026-10-02, GET /users/{id} with an id above 1 answered 200.
	"updateUserRights":                {Cmd: "user-role update", Unbounded: []string{"user_id"}},
	"listVirtualMachines":             {Cmd: "virtual-machine list"},
	"exportVirtualMachineSBOM":        {Cmd: "virtual-machine-sbom export"},
	"getWalletBalance":                {Cmd: "wallet-balance get"},
	"addWebhook":                      {Cmd: "webhook create"},
	"removeWebhook":                   {Cmd: "webhook delete"},
	"listWebhooks":                    {Cmd: "webhook list"},
	"getWorkspaceInfo":                {Cmd: "workspace get"},
	"getWorkspaceConfigurationErrors": {Cmd: "workspace-config-error list"},
}
